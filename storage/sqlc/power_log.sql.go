// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: power_log.sql

package storage

import (
	"context"
	"time"
)

const insertPowerLog = `-- name: InsertPowerLog :one
INSERT INTO power_log (
    timestamp,
    power_state
)
VALUES ($1, $2)
RETURNING id, timestamp, power_state
`

type InsertPowerLogParams struct {
	Timestamp  time.Time
	PowerState string
}

func (q *Queries) InsertPowerLog(ctx context.Context, arg InsertPowerLogParams) (PowerLog, error) {
	row := q.db.QueryRowContext(ctx, insertPowerLog, arg.Timestamp, arg.PowerState)
	var i PowerLog
	err := row.Scan(&i.ID, &i.Timestamp, &i.PowerState)
	return i, err
}
