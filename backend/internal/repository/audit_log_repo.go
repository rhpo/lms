package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// AuditLogRepository gère les opérations base de données pour les logs d'audit.
type AuditLogRepository struct {
	db *sql.DB
}

// NewAuditLogRepository crée un nouveau AuditLogRepository.
func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// FindByID cherche un log d'audit par son ID.
func (r *AuditLogRepository) FindByID(id string) (*entity.AuditLog, error) {
	row := r.db.QueryRow(`SELECT id, actor_id, action, entity, entity_id, metadata, created_at FROM audit_logs WHERE id = ?`, id)
	a := &entity.AuditLog{}
	err := row.Scan(&a.ID, &a.ActorID, &a.Action, &a.Entity, &a.EntityID, &a.Metadata, &a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

// FindAll retourne tous les logs d'audit.
func (r *AuditLogRepository) FindAll() ([]*entity.AuditLog, error) {
	rows, err := r.db.Query(`SELECT id, actor_id, action, entity, entity_id, metadata, created_at FROM audit_logs ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		a := &entity.AuditLog{}
		if err := rows.Scan(&a.ID, &a.ActorID, &a.Action, &a.Entity, &a.EntityID, &a.Metadata, &a.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, a)
	}
	return logs, rows.Err()
}

// Insert crée un nouveau log d'audit.
func (r *AuditLogRepository) Insert(a *entity.AuditLog) error {
	_, err := r.db.Exec(`INSERT INTO audit_logs (id, actor_id, action, entity, entity_id, metadata) VALUES (?, ?, ?, ?, ?, ?)`,
		a.ID, a.ActorID, a.Action, a.Entity, a.EntityID, a.Metadata)
	return err
}

// FindByEntityType retourne les logs pour une entité spécifique.
func (r *AuditLogRepository) FindByEntityType(entityType, entityID string) ([]*entity.AuditLog, error) {
	rows, err := r.db.Query(`SELECT id, actor_id, action, entity, entity_id, metadata, created_at
		FROM audit_logs WHERE entity = ? AND entity_id = ? ORDER BY created_at DESC`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		a := &entity.AuditLog{}
		if err := rows.Scan(&a.ID, &a.ActorID, &a.Action, &a.Entity, &a.EntityID, &a.Metadata, &a.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, a)
	}
	return logs, rows.Err()
}

// FindByActor retourne les logs pour un acteur spécifique.
func (r *AuditLogRepository) FindByActor(actorID string) ([]*entity.AuditLog, error) {
	rows, err := r.db.Query(`SELECT id, actor_id, action, entity, entity_id, metadata, created_at
		FROM audit_logs WHERE actor_id = ? ORDER BY created_at DESC`, actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		a := &entity.AuditLog{}
		if err := rows.Scan(&a.ID, &a.ActorID, &a.Action, &a.Entity, &a.EntityID, &a.Metadata, &a.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, a)
	}
	return logs, rows.Err()
}
