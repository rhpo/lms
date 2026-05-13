package entity

import (
	"database/sql"
	"time"
)

// Notification représente une notification in-app.
type Notification struct {
	ID          string       `json:"id"`
	RecipientID string       `json:"recipient_id"`
	Type        string       `json:"type"`
	Payload     string       `json:"payload"` // JSON
	ReadAt      sql.NullTime `json:"read_at"`
	CreatedAt   time.Time    `json:"created_at"`
}
