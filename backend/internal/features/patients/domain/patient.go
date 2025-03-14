package domain

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Date is a custom type that can handle both ISO 8601 and YYYY-MM-DD formats
type Date time.Time

// UnmarshalJSON implements json.Unmarshaler interface
func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")

	// Try parsing as YYYY-MM-DD first
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// If that fails, try parsing as ISO 8601
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}

	*d = Date(t)
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

// Time returns the underlying time.Time value
func (d Date) Time() time.Time {
	return time.Time(d)
}

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
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	MiddleName  string    `json:"middleName,omitempty"`
	DateOfBirth Date      `json:"dateOfBirth"`
	Gender      Gender    `json:"gender"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	Address     Address   `json:"address,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
		DateOfBirth: Date(dateOfBirth),
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
	age := now.Year() - p.DateOfBirth.Time().Year()

	// Adjust age if birthday hasn't occurred this year
	if now.Month() < p.DateOfBirth.Time().Month() ||
		(now.Month() == p.DateOfBirth.Time().Month() && now.Day() < p.DateOfBirth.Time().Day()) {
		age--
	}

	return age
}
