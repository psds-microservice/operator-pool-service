CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS operator_status (
  user_id UUID PRIMARY KEY,
  available BOOLEAN NOT NULL DEFAULT false,
  active_sessions INT NOT NULL DEFAULT 0,
  max_sessions INT NOT NULL DEFAULT 5,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_operator_status_available ON operator_status(available) WHERE available = true;
