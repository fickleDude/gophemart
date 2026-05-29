package model

import "time"

const (
	New        = "NEW"
	Processing = "PROCESSING"
	Invalid    = "INVALID"
	Processed  = "PROCESSED"
)

type Order struct {
	Login      string `json:"-"`
	Number     string
	Status     string
	Accrual    float64   `json:",omitempty"`
	UploadedAt time.Time `json:",omitempty"`
}
