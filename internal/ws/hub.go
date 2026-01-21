package ws

import "sync"

type Hub struct {
	mu    sync.RWMutex
	conns map[string]map[string]*Conn // userID -> deviceID -> Conn
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]map[string]*Conn)}
}

func (h *Hub) Run() {
	select {}
}

func (h *Hub) Register(userID, deviceID string, c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.conns[userID]; !ok {
		h.conns[userID] = make(map[string]*Conn)
	}
	h.conns[userID][deviceID] = c
}

func (h *Hub) Unregister(userID, deviceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if m, ok := h.conns[userID]; ok {
		delete(m, deviceID)
		if len(m) == 0 {
			delete(h.conns, userID)
		}
	}
}

func (h *Hub) BroadcastToUsers(userIDs []string, msg Envelope) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, uid := range userIDs {
		devs := h.conns[uid]
		for _, c := range devs {
			c.Send(msg)
		}
	}
}

