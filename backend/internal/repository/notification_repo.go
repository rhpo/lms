package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// NotificationRepository gère les opérations base de données pour les notifications.
type NotificationRepository struct {
	db *sql.DB
}

// NewNotificationRepository crée un nouveau NotificationRepository.
func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// FindByID cherche une notification par son ID.
func (r *NotificationRepository) FindByID(id string) (*entity.Notification, error) {
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

// FindByRecipient retourne les notifications d'un destinataire.
func (r *NotificationRepository) FindByRecipient(recipientID string) ([]*entity.Notification, error) {
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

// CountUnread retourne le nombre de notifications non lues pour un destinataire.
func (r *NotificationRepository) CountUnread(recipientID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE recipient_id = ? AND read_at IS NULL`, recipientID).Scan(&count)
	return count, err
}

// Insert crée une nouvelle notification.
func (r *NotificationRepository) Insert(n *entity.Notification) error {
	_, err := r.db.Exec(`INSERT INTO notifications (id, recipient_id, type, payload) VALUES (?, ?, ?, ?)`,
		n.ID, n.RecipientID, n.Type, n.Payload)
	return err
}

// MarkAsRead marque une notification comme lue.
func (r *NotificationRepository) MarkAsRead(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE id = ? AND read_at IS NULL`, id)
	return err
}

// MarkAllAsRead marque toutes les notifications d'un destinataire comme lues.
func (r *NotificationRepository) MarkAllAsRead(recipientID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = CURRENT_TIMESTAMP WHERE recipient_id = ? AND read_at IS NULL`, recipientID)
	return err
}
