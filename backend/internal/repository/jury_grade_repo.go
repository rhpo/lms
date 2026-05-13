package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// JuryGradeRepository gère les opérations base de données pour les notes du jury.
type JuryGradeRepository struct {
	db *sql.DB
}

// NewJuryGradeRepository crée un nouveau JuryGradeRepository.
func NewJuryGradeRepository(db *sql.DB) *JuryGradeRepository {
	return &JuryGradeRepository{db: db}
}

// FindByID cherche une note par son ID.
func (r *JuryGradeRepository) FindByID(id string) (*entity.JuryGrade, error) {
	row := r.db.QueryRow(`SELECT id, defense_id, jury_member_id, criterion1, criterion2, criterion3, criterion4, created_at, updated_at
		FROM jury_grades WHERE id = ?`, id)
	g := &entity.JuryGrade{}
	err := row.Scan(&g.ID, &g.DefenseID, &g.JuryMemberID, &g.Criterion1, &g.Criterion2, &g.Criterion3, &g.Criterion4, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return g, nil
}

// FindByDefense retourne les notes pour une soutenance.
func (r *JuryGradeRepository) FindByDefense(defenseID string) ([]*entity.JuryGrade, error) {
	rows, err := r.db.Query(`SELECT id, defense_id, jury_member_id, criterion1, criterion2, criterion3, criterion4, created_at, updated_at
		FROM jury_grades WHERE defense_id = ?`, defenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*entity.JuryGrade
	for rows.Next() {
		g := &entity.JuryGrade{}
		if err := rows.Scan(&g.ID, &g.DefenseID, &g.JuryMemberID, &g.Criterion1, &g.Criterion2, &g.Criterion3, &g.Criterion4, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		grades = append(grades, g)
	}
	return grades, rows.Err()
}

// FindByDefenseAndMember retourne la note d'un membre du jury pour une soutenance.
func (r *JuryGradeRepository) FindByDefenseAndMember(defenseID, juryMemberID string) (*entity.JuryGrade, error) {
	row := r.db.QueryRow(`SELECT id, defense_id, jury_member_id, criterion1, criterion2, criterion3, criterion4, created_at, updated_at
		FROM jury_grades WHERE defense_id = ? AND jury_member_id = ?`, defenseID, juryMemberID)
	g := &entity.JuryGrade{}
	err := row.Scan(&g.ID, &g.DefenseID, &g.JuryMemberID, &g.Criterion1, &g.Criterion2, &g.Criterion3, &g.Criterion4, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return g, nil
}

// Insert crée une nouvelle note de jury.
func (r *JuryGradeRepository) Insert(g *entity.JuryGrade) error {
	_, err := r.db.Exec(`INSERT INTO jury_grades (id, defense_id, jury_member_id, criterion1, criterion2, criterion3, criterion4)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, g.ID, g.DefenseID, g.JuryMemberID, g.Criterion1, g.Criterion2, g.Criterion3, g.Criterion4)
	return err
}

// Update met à jour une note de jury.
func (r *JuryGradeRepository) Update(g *entity.JuryGrade) error {
	_, err := r.db.Exec(`UPDATE jury_grades SET criterion1 = ?, criterion2 = ?, criterion3 = ?, criterion4 = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		g.Criterion1, g.Criterion2, g.Criterion3, g.Criterion4, g.ID)
	return err
}

// DeleteByDefense supprime toutes les notes d'une soutenance.
func (r *JuryGradeRepository) DeleteByDefense(defenseID string) error {
	_, err := r.db.Exec(`DELETE FROM jury_grades WHERE defense_id = ?`, defenseID)
	return err
}
