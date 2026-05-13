package entity

import (
	"database/sql"
	"time"
)

// AuditLog représente une entrée dans le journal d'audit.
type AuditLog struct {
	ID        string         `json:"id"`
	ActorID   string         `json:"actor_id"`
	Action    string         `json:"action"`
	Entity    string         `json:"entity"`
	EntityID  sql.NullString `json:"entity_id"`
	Metadata  sql.NullString `json:"metadata"` // JSON
	CreatedAt time.Time      `json:"created_at"`
}
