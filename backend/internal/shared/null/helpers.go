package null

import (
	"database/sql"
	"time"
)

// StringToPtr convertit une sql.NullString en *string.
func StringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// PtrToString convertit un *string en sql.NullString.
func PtrToString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

// TimeToPtr convertit une sql.NullTime en *time.Time.
func TimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}

// PtrToTime convertit un *time.Time en sql.NullTime.
func PtrToTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

// Float64ToPtr convertit une sql.NullFloat64 en *float64.
func Float64ToPtr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

// PtrToFloat64 convertit un *float64 en sql.NullFloat64.
func PtrToFloat64(f *float64) sql.NullFloat64 {
	if f != nil {
		return sql.NullFloat64{Float64: *f, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

// NowPtr retourne un pointeur vers le temps actuel.
func NowPtr() *time.Time {
	t := time.Now()
	return &t
}

// StringPtr retourne un pointeur vers une chaîne.
func StringPtr(s string) *string {
	return &s
}

// Float64Ptr retourne un pointeur vers un float64.
func Float64Ptr(f float64) *float64 {
	return &f
}
