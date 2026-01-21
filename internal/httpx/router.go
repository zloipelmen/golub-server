package httpx

import (
	"net/http"

	"messenger/internal/auth"
	"messenger/internal/storage"
	"messenger/internal/ws"
)

func NewRouter(hub *ws.Hub, authSvc *auth.Service, store *storage.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", Healthz)

	wsHandler := ws.NewHandler(hub, authSvc, store)
	mux.HandleFunc("/ws", wsHandler.ServeWS)

	return mux
}
