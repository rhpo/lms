package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) FindByID(id int64) (*entity.AuditLog, error) {
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

func (r *AuditLogRepository) Insert(a *entity.AuditLog) error {
	result, err := r.db.Exec(`INSERT INTO audit_logs (actor_id, action, entity, entity_id, metadata) VALUES (?, ?, ?, ?, ?)`,
		a.ActorID, a.Action, a.Entity, a.EntityID, a.Metadata)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (r *AuditLogRepository) FindByEntityType(entityType string, entityID int64) ([]*entity.AuditLog, error) {
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

func (r *AuditLogRepository) FindByActor(actorID int64) ([]*entity.AuditLog, error) {
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
