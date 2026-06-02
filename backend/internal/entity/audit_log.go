package entity

import (
	"time"
)

// AuditLog représente une entrée dans le journal d'audit.
type AuditLog struct {
	ID        int64      `json:"id"`
	ActorID   int64      `json:"actor_id"`
	Action    string     `json:"action"`
	Entity    string     `json:"entity"`
	EntityID  NullInt64  `json:"entity_id"`
	Metadata  NullString `json:"metadata"` // JSON
	CreatedAt time.Time  `json:"created_at"`
}
