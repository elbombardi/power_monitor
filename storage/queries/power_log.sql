-- name: InsertPowerLog :one
INSERT INTO power_log (
    timestamp,
    power_state
)
VALUES ($1, $2)
RETURNING *;