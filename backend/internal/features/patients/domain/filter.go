package domain

import (
	"strings"
	"time"
)

// PatientListFilter holds optional list/search criteria for patient queries.
type PatientListFilter struct {
	Search          string
	Gender          string
	DateOfBirthFrom string // YYYY-MM-DD
	DateOfBirthTo   string // YYYY-MM-DD
	RegisteredSince string // YYYY-MM-DD — patients created on or after this date
}

// ValidGender reports whether gender is empty or a known patient gender value.
func ValidGender(gender string) bool {
	if gender == "" {
		return true
	}
	switch Gender(gender) {
	case GenderMale, GenderFemale, GenderOther, GenderUnknown:
		return true
	default:
		return false
	}
}

func parseFilterDate(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

// MatchesFilter returns true when a patient satisfies all non-empty filter fields.
func MatchesFilter(patient *Patient, filter PatientListFilter) bool {
	search := strings.TrimSpace(filter.Search)
	if search != "" {
		needle := strings.ToLower(search)
		haystack := strings.ToLower(patient.FullName() + " " + patient.Email)
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	if filter.Gender != "" && string(patient.Gender) != filter.Gender {
		return false
	}

	if from, ok := parseFilterDate(filter.DateOfBirthFrom); ok {
		if patient.DateOfBirth.Time().Before(from) {
			return false
		}
	}

	if to, ok := parseFilterDate(filter.DateOfBirthTo); ok {
		if patient.DateOfBirth.Time().After(to) {
			return false
		}
	}

	if since, ok := parseFilterDate(filter.RegisteredSince); ok {
		if patient.CreatedAt.Before(since) {
			return false
		}
	}

	return true
}
