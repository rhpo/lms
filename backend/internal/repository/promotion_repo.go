package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type PromotionRepository struct {
	db *sql.DB
}

func NewPromotionRepository(db *sql.DB) *PromotionRepository {
	return &PromotionRepository{db: db}
}

func (r *PromotionRepository) FindByID(id int64) (*entity.Promotion, error) {
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

func (r *PromotionRepository) Insert(p *entity.Promotion) error {
	result, err := r.db.Exec(`INSERT INTO promotions (label, academic_year_id) VALUES (?, ?)`, p.Label, p.AcademicYearID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *PromotionRepository) Update(p *entity.Promotion) error {
	_, err := r.db.Exec(`UPDATE promotions SET label = ?, academic_year_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, p.Label, p.AcademicYearID, p.ID)
	return err
}

func (r *PromotionRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM promotions WHERE id = ?`, id)
	return err
}
