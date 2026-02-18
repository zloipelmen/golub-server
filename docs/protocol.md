# Golub WS Protocol v1 (MVP)

## Transport
- WebSocket endpoint: `wss://<host>/ws`
- All messages are JSON frames.
- Timestamps are strings in RFC3339Nano UTC (example: `2026-02-18T07:20:32.123456789Z`).

## Envelope
Every WS frame is:

```json
{
  "type": "auth",
  "req_id": "uuid-optional",
  "payload": {}
}```
```
{
  "type": "error",
  "req_id": "same-as-request-if-any",
  "payload": { "code": "BAD_REQUEST", "message": "..." }
}```
```
{
  "type": "auth",
  "req_id": "1",
  "payload": {
    "invite_code": "...",
    "device_key": "...",
    "device_name": "Pixel 7",
    "app_version": "0.1.0"
  }
}```
```
{
  "type": "auth_ok",
  "req_id": "1",
  "payload": {
    "user_id": "u_...",
    "device_id": "d_...",
    "server_time": "2026-02-18T07:20:32.123456789Z"
  }
}```
```
{
  "type": "send_message",
  "req_id": "2",
  "payload": {
    "chat_id": "c_...",
    "client_msg_id": "uuid",
    "text": "hello"
  }
}```
```
{
  "type": "send_message_ok",
  "req_id": "2",
  "payload": {
    "message_id": "m_...",
    "created_at": "2026-02-18T07:21:00.000000000Z"
  }
}```
```
{
  "type": "message",
  "payload": {
    "message_id": "m_...",
    "chat_id": "c_...",
    "sender_user_id": "u_...",
    "text": "hello",
    "created_at": "2026-02-18T07:21:00.000000000Z"
  }
}```
```
{
  "type": "sync",
  "req_id": "3",
  "payload": {
    "chat_id": "c_...",
    "limit": 50
  }
}```
```
{
  "type": "sync_ok",
  "req_id": "3",
  "payload": {
    "chat_id": "c_...",
    "messages": [
      {
        "message_id": "m_...",
        "chat_id": "c_...",
        "sender_user_id": "u_...",
        "text": "hello",
        "created_at": "2026-02-18T07:21:00.000000000Z"
      }
    ]
  }
}```
