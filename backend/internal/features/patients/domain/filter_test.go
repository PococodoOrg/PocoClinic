package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMatchesFilter(t *testing.T) {
	created := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	patient := &Patient{
		ID:          uuid.New(),
		FirstName:   "Jane",
		LastName:    "Smith",
		Email:       "jane@clinic.local",
		DateOfBirth: Date(time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)),
		Gender:      GenderFemale,
		CreatedAt:   created,
	}

	if !MatchesFilter(patient, PatientListFilter{Search: "smith"}) {
		t.Fatal("expected name search to match")
	}
	if MatchesFilter(patient, PatientListFilter{Gender: "male"}) {
		t.Fatal("expected gender mismatch")
	}
	if !MatchesFilter(patient, PatientListFilter{RegisteredSince: "2026-03-01"}) {
		t.Fatal("expected registered since to match")
	}
	if MatchesFilter(patient, PatientListFilter{DateOfBirthTo: "1980-01-01"}) {
		t.Fatal("expected dob upper bound to exclude patient")
	}
}

func TestValidGender(t *testing.T) {
	if !ValidGender("female") || ValidGender("invalid") {
		t.Fatal("expected gender validation to accept known values only")
	}
}
