package ws

import "encoding/json"

type Envelope struct {
	Type    string          `json:"type"`
	ReqID   string          `json:"req_id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AuthPayload struct {
	InviteCode string `json:"invite_code"`
	DeviceKey  string `json:"device_key"`
	DeviceName string `json:"device_name"`
	AppVersion string `json:"app_version"`
}

type AuthOKPayload struct {
	UserID     string `json:"user_id"`
	DeviceID   string `json:"device_id"`
	ServerTime string `json:"server_time"`
}

type SendMessagePayload struct {
	ChatID      string `json:"chat_id"`
	ClientMsgID string `json:"client_msg_id"`
	Text        string `json:"text"`
}

type SendMessageOKPayload struct {
	MessageID string `json:"message_id"`
	CreatedAt string `json:"created_at"`
}

type MessagePayload struct {
	MessageID    string `json:"message_id"`
	ChatID       string `json:"chat_id"`
	SenderUserID string `json:"sender_user_id"`
	Text         string `json:"text"`
	CreatedAt    string `json:"created_at"`
}

type SyncPayload struct {
	ChatID string `json:"chat_id"`
	Limit  int    `json:"limit"`
}

type SyncOKPayload struct {
	ChatID    string          `json:"chat_id"`
	Messages  []MessagePayload `json:"messages"`
}
