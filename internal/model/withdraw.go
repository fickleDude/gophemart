package model

import "time"

type Withdraw struct {
	Login       string `json:"-"`
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
