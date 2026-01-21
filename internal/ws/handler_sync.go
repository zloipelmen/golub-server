package ws

import (
	"context"
	"encoding/json"
	"time"
)

func (h *Handler) handleSync(ctx context.Context, c *Conn, userID string, in Envelope) {
	var p SyncPayload
	if err := json.Unmarshal(in.Payload, &p); err != nil {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "BAD_REQUEST", Message: "invalid sync payload",
		})})
		return
	}
	if p.ChatID == "" {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "BAD_REQUEST", Message: "chat_id required",
		})})
		return
	}

	isMember, err := h.store.IsChatMember(ctx, p.ChatID, userID)
	if err != nil {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "INTERNAL", Message: "db error",
		})})
		return
	}
	if !isMember {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "FORBIDDEN", Message: "not a member of chat",
		})})
		return
	}

	msgs, err := h.store.ListMessages(ctx, p.ChatID, p.Limit)
	if err != nil {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "INTERNAL", Message: "db error",
		})})
		return
	}

	out := make([]MessagePayload, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, MessagePayload{
			MessageID:    m.ID,
			ChatID:       m.ChatID,
			SenderUserID: m.SenderUserID,
			Text:         m.Text,
			CreatedAt:    m.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}

	c.Send(Envelope{Type: "sync_ok", ReqID: in.ReqID, Payload: mustJSON(SyncOKPayload{
		ChatID:   p.ChatID,
		Messages: out,
	})})
}
