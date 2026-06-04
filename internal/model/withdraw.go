package model

type Withdraw struct {
	Login       string  `json:"-"`
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
