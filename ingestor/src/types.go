package main

import (
	"time"
	"uuid"
)

type Reading struct {
	Schema_version uint8     `json:"schema_version"`
	Id             uuid.UUID `json:"id" binding:"required"`
	Device_id      uuid.UUID `json:"device_id" binding:"required"`
	Metric         string    `json:"metric" binding:"required"`
	Value          float32   `json:"value"`
	Timestamp      time.Time `json:"timestamp" binding:"required"`
	Reliability    uint8     `json:"reliability"`
}
