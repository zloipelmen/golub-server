package ws

import (
	"context"
	"encoding/json"
	"time"

	"messenger/internal/storage"
)

func (h *Handler) handleSendMessage(ctx context.Context, c *Conn, userID, deviceID string, in Envelope) {
	var p SendMessagePayload
	if err := json.Unmarshal(in.Payload, &p); err != nil {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "BAD_REQUEST", Message: "invalid send_message payload",
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

	msg, err := h.store.CreateMessage(ctx, storage.CreateMessageParams{
		ChatID:       p.ChatID,
		SenderUserID: userID,
		ClientMsgID:  p.ClientMsgID,
		Text:         p.Text,
	})
	if err != nil {
		c.Send(Envelope{Type: "error", ReqID: in.ReqID, Payload: mustJSON(ErrorPayload{
			Code: "INTERNAL", Message: err.Error(),
		})})
		return
	}

	c.Send(Envelope{Type: "send_message_ok", ReqID: in.ReqID, Payload: mustJSON(SendMessageOKPayload{
		MessageID: msg.ID,
		CreatedAt: msg.CreatedAt.UTC().Format(time.RFC3339Nano),
	})})

	memberIDs, err := h.store.ListChatMemberUserIDs(ctx, p.ChatID)
	if err != nil {
		return
	}

	h.hub.BroadcastToUsers(memberIDs, Envelope{Type: "message", Payload: mustJSON(MessagePayload{
		MessageID:    msg.ID,
		ChatID:       msg.ChatID,
		SenderUserID: msg.SenderUserID,
		Text:         msg.Text,
		CreatedAt:    msg.CreatedAt.UTC().Format(time.RFC3339Nano),
	})})
}
