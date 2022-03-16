package views

import "time"

type VCQRVerifyView struct {
	VCs       []string   `json:"vcs"`
	SenderDID string     `json:"sender_did"`
	CreatedAt *time.Time `json:"created_at"`
}
