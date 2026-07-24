package handlers

import (
	"time"

	"github.com/dksch/pococlinic/internal/features/auth/domain"
)

// StaffUserJSON is the admin-facing staff account shape returned by auth APIs.
type StaffUserJSON struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	Role          domain.Role `json:"role"`
	LastLogin     *time.Time `json:"lastLogin,omitempty"`
	MustChangePIN bool       `json:"mustChangePin"`
	IsActive      bool       `json:"isActive"`
	IsLocked      bool       `json:"isLocked"`
	LockedUntil   *time.Time `json:"lockedUntil,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func toStaffUserJSON(user *domain.User) StaffUserJSON {
	return StaffUserJSON{
		ID:            user.ID.String(),
		Email:         user.Email,
		Name:          user.Name,
		Role:          user.Role,
		LastLogin:     user.LastLogin,
		MustChangePIN: user.MustChangePIN,
		IsActive:      user.IsActive,
		IsLocked:      user.IsLocked(),
		LockedUntil:   user.LockedUntil,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
