package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type WebhookInfo struct {
	ID          uuid.NullUUID `json:"id" gorm:"primaryKey;default=uuid_generate_v4()"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   sql.NullTime  `json:"deletedAt"`
	AccountID   string        `json:"account_id" gorm:"not null;index"`
	CallbackURL string        `json:"callback_url" gorm:"not null"`
	Secret      string        `json:"secret" gorm:"not null"`
}
