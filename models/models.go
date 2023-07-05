package models

import "time"

type PowerState string

const (
	POWERSTATE_ON  PowerState = "ON"
	POWERSTATE_OFF PowerState = "OFF"
)

type PowerStateChange struct {
	Id         int        `json:"id"`
	Time       time.Time  `json:"time"`
	PowerState PowerState `json:"powerState"`
}
