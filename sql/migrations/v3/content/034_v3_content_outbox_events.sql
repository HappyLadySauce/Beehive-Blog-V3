-- v3 content: outbox events
-- Content domain events are written to outbox first, then reliably dispatched to RabbitMQ.
-- 内容领域事件先写入 outbox，再由 dispatcher 可靠投递到 RabbitMQ

CREATE TABLE content.outbox_events (
  id BIGSERIAL PRIMARY KEY,
  event_id VARCHAR(64) NOT NULL,
  event_type VARCHAR(128) NOT NULL,
  resource_type VARCHAR(64) NOT NULL,
  resource_id BIGINT NOT NULL,
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  attempts INT NOT NULL DEFAULT 0,
  next_retry_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_error TEXT NOT NULL DEFAULT '',
  published_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_content_outbox_events_status CHECK (status IN ('pending', 'processing', 'done', 'failed')),
  CONSTRAINT chk_content_outbox_events_attempts CHECK (attempts >= 0)
);

CREATE UNIQUE INDEX ux_content_outbox_events_event_id ON content.outbox_events (event_id);
CREATE INDEX idx_content_outbox_events_poll ON content.outbox_events (status, next_retry_at, id);
CREATE INDEX idx_content_outbox_events_resource ON content.outbox_events (resource_type, resource_id, id DESC);
