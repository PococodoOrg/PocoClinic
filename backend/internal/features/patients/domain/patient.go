package domain

import (
	"time"

	"github.com/google/uuid"
)

// Gender represents the patient's gender identity
type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderUnknown Gender = "unknown"
)

// Patient represents the core patient domain model
type Patient struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	MiddleName    string    `json:"middleName,omitempty"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	Gender        Gender    `json:"gender"`
	Email         string    `json:"email,omitempty"`
	PhoneNumber   string    `json:"phoneNumber,omitempty"`
	Address       Address   `json:"address,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	MedicalNumber string    `json:"medicalNumber"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// NewPatient creates a new patient with a generated ID and timestamps
func NewPatient(firstName, lastName string, dateOfBirth time.Time, gender Gender) *Patient {
	now := time.Now()
	return &Patient{
		ID:          uuid.New(),
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dateOfBirth,
		Gender:      gender,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the patient's updatedAt timestamp
func (p *Patient) Update() {
	p.UpdatedAt = time.Now()
}

// FullName returns the patient's full name
func (p *Patient) FullName() string {
	if p.MiddleName != "" {
		return p.FirstName + " " + p.MiddleName + " " + p.LastName
	}
	return p.FirstName + " " + p.LastName
}

// Age calculates the patient's current age
func (p *Patient) Age() int {
	now := time.Now()
	age := now.Year() - p.DateOfBirth.Year()
	
	// Adjust age if birthday hasn't occurred this year
	if now.Month() < p.DateOfBirth.Month() || 
		(now.Month() == p.DateOfBirth.Month() && now.Day() < p.DateOfBirth.Day()) {
		age--
	}
	
	return age
} 