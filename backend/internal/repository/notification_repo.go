package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) FindByID(id int64) (*entity.Notification, error) {
	row := r.db.QueryRow(`SELECT id, recipient_id, type, payload, read_at, created_at FROM notifications WHERE id = ?`, id)
	n := &entity.Notification{}
	err := row.Scan(&n.ID, &n.RecipientID, &n.Type, &n.Payload, &n.ReadAt, &n.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return n, nil
}

func (r *NotificationRepository) FindByRecipient(recipientID int64) ([]*entity.Notification, error) {
	rows, err := r.db.Query(`SELECT id, recipient_id, type, payload, read_at, created_at
		FROM notifications WHERE recipient_id = ? ORDER BY created_at DESC`, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []*entity.Notification
	for rows.Next() {
		n := &entity.Notification{}
		if err := rows.Scan(&n.ID, &n.RecipientID, &n.Type, &n.Payload, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *NotificationRepository) CountUnread(recipientID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE recipient_id = ? AND read_at IS NULL`, recipientID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) Insert(n *entity.Notification) error {
	result, err := r.db.Exec(`INSERT INTO notifications (recipient_id, type, payload) VALUES (?, ?, ?)`,
		n.RecipientID, n.Type, n.Payload)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = id
	return nil
}

func (r *NotificationRepository) MarkAsRead(id int64) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE id = ? AND read_at IS NULL`, id)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(recipientID int64) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE recipient_id = ? AND read_at IS NULL`, recipientID)
	return err
}
