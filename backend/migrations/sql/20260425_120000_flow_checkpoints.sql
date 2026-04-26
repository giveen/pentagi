-- +goose Up
-- +goose StatementBegin
CREATE TABLE flow_checkpoints (
  id          BIGINT       PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  flow_id     BIGINT       NOT NULL REFERENCES flows(id) ON DELETE CASCADE,
  kind        TEXT         NOT NULL DEFAULT 'pause_checkpoint_v1',
  payload     JSONB        NOT NULL,
  active      BOOLEAN      NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  consumed_at TIMESTAMPTZ
);

CREATE INDEX flow_checkpoints_flow_id_idx ON flow_checkpoints(flow_id);
CREATE INDEX flow_checkpoints_flow_active_idx ON flow_checkpoints(flow_id, active, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS flow_checkpoints;
-- +goose StatementEnd
