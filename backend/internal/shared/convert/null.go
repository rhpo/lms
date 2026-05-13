package convert

import (
	"database/sql"
	"time"
)

// StringPtr converts sql.NullString to *string.
func StringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// TimePtr converts sql.NullTime to *time.Time.
func TimePtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// NullString converts *string to sql.NullString.
func NullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

// NullTime converts *time.Time to sql.NullTime.
func NullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}
