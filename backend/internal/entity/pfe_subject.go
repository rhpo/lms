package entity

import (
	"database/sql"
	"time"
)

// PfeSubject représente un sujet de PFE proposé par un enseignant ou une entreprise.
type PfeSubject struct {
	ID                    string         `json:"id"`
	Title                 string         `json:"title"`
	Description           string         `json:"description"`
	GroupType             string         `json:"group_type"` // monome/binome/trinome
	ProposerID            string         `json:"proposer_id"`
	ProposerRole          string         `json:"proposer_role"` // teacher/company
	CompanyID             sql.NullString `json:"company_id"`
	AcademicYearID        string         `json:"academic_year_id"`
	Validator1ID          sql.NullString `json:"validator1_id"`
	Validator2ID          sql.NullString `json:"validator2_id"`
	Validator1Decision    sql.NullString `json:"validator1_decision"` // valide/accepte_sous_reserve/refuse
	Validator2Decision    sql.NullString `json:"validator2_decision"`
	Validator1Comment     sql.NullString `json:"validator1_comment"`
	Validator2Comment     sql.NullString `json:"validator2_comment"`
	Status                string         `json:"status"` // en_attente/valide/accepte_sous_reserve/refuse/expire
	CoSupervisorID        sql.NullString `json:"co_supervisor_id"`
	PreAssignedStudentIDs sql.NullString `json:"pre_assigned_student_ids"` // JSON array d'IDs
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`

	// Relations
	Proposer     *Profile  `json:"proposer,omitempty"`
	Company      *Company  `json:"company,omitempty"`
	Validator1   *Teacher  `json:"validator1,omitempty"`
	Validator2   *Teacher  `json:"validator2,omitempty"`
	CoSupervisor *Teacher  `json:"co_supervisor,omitempty"`
	Domains      []*Domain `json:"domains,omitempty"`
}
