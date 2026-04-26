-- name: CreateFlowCheckpoint :one
INSERT INTO flow_checkpoints (
  flow_id,
  kind,
  payload,
  active
)
VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetLatestActiveFlowCheckpoint :one
SELECT
  fc.*
FROM flow_checkpoints fc
INNER JOIN flows f ON fc.flow_id = f.id
WHERE fc.flow_id = $1 AND fc.active = TRUE AND f.deleted_at IS NULL
ORDER BY fc.created_at DESC, fc.id DESC
LIMIT 1;

-- name: DeactivateFlowCheckpoints :exec
UPDATE flow_checkpoints
SET active = FALSE,
    consumed_at = CURRENT_TIMESTAMP
WHERE flow_id = $1 AND active = TRUE;

-- name: MarkFlowCheckpointConsumed :one
UPDATE flow_checkpoints
SET active = FALSE,
    consumed_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
