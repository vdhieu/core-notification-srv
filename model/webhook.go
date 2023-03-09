package model

import "time"

type WebhookInfo struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	AccountID   string    `json:"account_id"`
	CallbackURL string    `json:"callback_url"`
	Secret      string    `json:"secret"`
}
