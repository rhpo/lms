package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// PromotionRepository gère les opérations base de données pour les promotions.
type PromotionRepository struct {
	db *sql.DB
}

// NewPromotionRepository crée un nouveau PromotionRepository.
func NewPromotionRepository(db *sql.DB) *PromotionRepository {
	return &PromotionRepository{db: db}
}

// FindByID cherche une promotion par son ID.
func (r *PromotionRepository) FindByID(id string) (*entity.Promotion, error) {
	row := r.db.QueryRow(`SELECT id, label, academic_year_id, created_at, updated_at FROM promotions WHERE id = ?`, id)
	p := &entity.Promotion{}
	err := row.Scan(&p.ID, &p.Label, &p.AcademicYearID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// FindAll retourne toutes les promotions.
func (r *PromotionRepository) FindAll() ([]*entity.Promotion, error) {
	rows, err := r.db.Query(`SELECT id, label, academic_year_id, created_at, updated_at FROM promotions ORDER BY label`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promotions []*entity.Promotion
	for rows.Next() {
		p := &entity.Promotion{}
		if err := rows.Scan(&p.ID, &p.Label, &p.AcademicYearID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		promotions = append(promotions, p)
	}
	return promotions, rows.Err()
}

// Insert crée une nouvelle promotion.
func (r *PromotionRepository) Insert(p *entity.Promotion) error {
	_, err := r.db.Exec(`INSERT INTO promotions (id, label, academic_year_id) VALUES (?, ?, ?)`, p.ID, p.Label, p.AcademicYearID)
	return err
}

// Update met à jour une promotion.
func (r *PromotionRepository) Update(p *entity.Promotion) error {
	_, err := r.db.Exec(`UPDATE promotions SET label = ?, academic_year_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, p.Label, p.AcademicYearID, p.ID)
	return err
}

// Delete supprime une promotion.
func (r *PromotionRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM promotions WHERE id = ?`, id)
	return err
}
