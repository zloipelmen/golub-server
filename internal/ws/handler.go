package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"net"
	"net/url"
	"strings"

	"messenger/internal/auth"
	"messenger/internal/storage"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = 25 * time.Second
	writeWait  = 10 * time.Second
)

type Handler struct {
	hub     *Hub
	authSvc *auth.Service
	store   *storage.Store
}

func NewHandler(hub *Hub, authSvc *auth.Service, store *storage.Store) *Handler {
	return &Handler{hub: hub, authSvc: authSvc, store: store}
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        if origin == "" {
            // Non-browser clients (wscat, etc.)
            return true
        }

        u, err := url.Parse(origin)
        if err != nil {
            log.Printf("ws origin parse failed: origin=%q err=%v", origin, err)
            return false
        }

        originHost := u.Hostname()

        reqHost := r.Host
        if h, _, err := net.SplitHostPort(reqHost); err == nil {
            reqHost = h
        }

        // Allow same-host
        if strings.EqualFold(originHost, reqHost) {
            return true
        }

        // Allow local dev
        if originHost == "localhost" || originHost == "127.0.0.1" {
            return true
        }

        log.Printf("ws origin rejected: origin=%q originHost=%q reqHost=%q hostHeader=%q",
            origin, originHost, reqHost, r.Host,
        )
        return false
    },
}


func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("ws upgrade request: host=%q origin=%q ua=%q remote=%q",
		r.Host,
		r.Header.Get("Origin"),
		r.Header.Get("User-Agent"),
		r.RemoteAddr,
	)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}
	wsConn.SetReadLimit(1 << 20) // 1MB
	_ = wsConn.SetReadDeadline(time.Now().Add(pongWait))
	wsConn.SetPongHandler(func(string) error {
		_ = wsConn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	conn := NewConn(wsConn)
	go conn.WriteLoop()
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for range ticker.C {
			_ = wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				conn.Close()
				return
			}
		}
	}()

	_ = wsConn.SetReadDeadline(time.Now().Add(20 * time.Second))
	var env Envelope
	if err := wsConn.ReadJSON(&env); err != nil || env.Type != "auth" {
		conn.Send(Envelope{Type: "error", ReqID: env.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "AUTH_REQUIRED", Message: "first message must be auth",
		})})
		conn.Close()
		return
	}

	var ap AuthPayload
	if err := json.Unmarshal(env.Payload, &ap); err != nil {
		conn.Send(Envelope{Type: "error", ReqID: env.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "BAD_REQUEST", Message: "invalid auth payload",
		})})
		conn.Close()
		return
	}

	userID, deviceID, err := h.authSvc.AuthByInvite(r.Context(), ap.InviteCode, ap.DeviceKey, ap.DeviceName)
	if err != nil {
		conn.Send(Envelope{Type: "error", ReqID: env.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "AUTH_FAILED", Message: err.Error(),
		})})
		conn.Close()
		return
	}

	h.hub.Register(userID, deviceID, conn)
	defer h.hub.Unregister(userID, deviceID)

	conn.Send(Envelope{Type: "auth_ok", ReqID: env.ReqID, Payload: mustJSON(AuthOKPayload{
		UserID: userID, DeviceID: deviceID, ServerTime: time.Now().UTC().Format(time.RFC3339Nano),
	})})

	_ = wsConn.SetReadDeadline(time.Time{})

	for {
		var in Envelope
		if err := wsConn.ReadJSON(&in); err != nil {
			log.Printf("ws read closed: %v", err)
			return
		}

		switch in.Type {
		case "sync":
			h.handleSync(r.Context(), conn, userID, in)
		case "send_message":
			h.handleSendMessage(r.Context(), conn, userID, deviceID, in)
		default:
			conn.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
				Code: "UNKNOWN_TYPE", Message: "unknown event type",
			})})
		}
	}
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
