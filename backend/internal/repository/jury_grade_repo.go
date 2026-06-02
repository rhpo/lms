package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type JuryGradeRepository struct {
	db *sql.DB
}

func NewJuryGradeRepository(db *sql.DB) *JuryGradeRepository {
	return &JuryGradeRepository{db: db}
}

const juryGradeSelect = `SELECT id, defense_id, jury_member_id, criterion1, criterion2, criterion3, criterion4, archive_decision, created_at, updated_at FROM jury_grades`

func scanGrade(row interface{ Scan(...any) error }) (*entity.JuryGrade, error) {
	g := &entity.JuryGrade{}
	err := row.Scan(&g.ID, &g.DefenseID, &g.JuryMemberID,
		&g.Criterion1, &g.Criterion2, &g.Criterion3, &g.Criterion4,
		&g.ArchiveDecision, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return g, nil
}

func (r *JuryGradeRepository) FindByID(id int64) (*entity.JuryGrade, error) {
	return scanGrade(r.db.QueryRow(juryGradeSelect+` WHERE id = ?`, id))
}

func (r *JuryGradeRepository) FindByDefense(defenseID int64) ([]*entity.JuryGrade, error) {
	rows, err := r.db.Query(juryGradeSelect+` WHERE defense_id = ?`, defenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var grades []*entity.JuryGrade
	for rows.Next() {
		g, err := scanGrade(rows)
		if err != nil {
			return nil, err
		}
		grades = append(grades, g)
	}
	return grades, rows.Err()
}

func (r *JuryGradeRepository) FindByDefenseAndMember(defenseID, juryMemberID int64) (*entity.JuryGrade, error) {
	return scanGrade(r.db.QueryRow(juryGradeSelect+` WHERE defense_id = ? AND jury_member_id = ?`, defenseID, juryMemberID))
}

func (r *JuryGradeRepository) Insert(g *entity.JuryGrade) error {
	result, err := r.db.Exec(
		`INSERT INTO jury_grades (defense_id, jury_member_id, criterion1, criterion2, criterion3, criterion4, archive_decision)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		g.DefenseID, g.JuryMemberID, g.Criterion1, g.Criterion2, g.Criterion3, g.Criterion4, g.ArchiveDecision)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	g.ID = id
	return nil
}

func (r *JuryGradeRepository) Update(g *entity.JuryGrade) error {
	_, err := r.db.Exec(
		`UPDATE jury_grades SET criterion1 = ?, criterion2 = ?, criterion3 = ?, criterion4 = ?, archive_decision = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		g.Criterion1, g.Criterion2, g.Criterion3, g.Criterion4, g.ArchiveDecision, g.ID)
	return err
}

func (r *JuryGradeRepository) DeleteByDefense(defenseID int64) error {
	_, err := r.db.Exec(`DELETE FROM jury_grades WHERE defense_id = ?`, defenseID)
	return err
}
