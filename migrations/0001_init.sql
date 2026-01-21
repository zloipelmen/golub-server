CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE invites (
  code text PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  expires_at timestamptz NULL,
  max_uses int NOT NULL DEFAULT 1,
  uses int NOT NULL DEFAULT 0,
  disabled boolean NOT NULL DEFAULT false
);

CREATE TABLE devices (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  device_key text NOT NULL UNIQUE,
  device_name text NOT NULL DEFAULT '',
  last_seen_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE chats (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  title text NOT NULL DEFAULT '',
  is_direct boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE chat_members (
  chat_id uuid NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  joined_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  chat_id uuid NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
  sender_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  client_msg_id uuid NOT NULL,
  text text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (sender_user_id, client_msg_id)
);

CREATE INDEX idx_messages_chat_created ON messages(chat_id, created_at);

CREATE TABLE message_acks (
  message_id uuid NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  device_id uuid NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
  acked_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (message_id, device_id)
);

CREATE TABLE attachments (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  message_id uuid NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  type text NOT NULL,
  url text NOT NULL DEFAULT '',
  meta jsonb NOT NULL DEFAULT '{}'::jsonb
);
