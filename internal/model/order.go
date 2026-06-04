package model

const (
	New        = "NEW"
	Processing = "PROCESSING"
	Invalid    = "INVALID"
	Processed  = "PROCESSED"
)

type Order struct {
	Login      string  `json:"-"`
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at",omitempty`
}
