package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Conn struct {
	ws     *websocket.Conn
	send   chan Envelope
	mu     sync.Mutex
	closed bool
}

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{ws: ws, send: make(chan Envelope, 64)}
}

func (c *Conn) WriteLoop() {
	defer c.Close()
	for msg := range c.send {
		if err := c.ws.WriteJSON(msg); err != nil {
			log.Printf("ws write: %v", err)
			return
		}
	}
}

func (c *Conn) Send(msg Envelope) {
	select {
	case c.send <- msg:
	default:
		_ = c.ws.Close()
	}
}

func (c *Conn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	close(c.send)
	_ = c.ws.Close()
}
