package infrastructure

import (
	"fmt"
	"strings"

	"github.com/dksch/pococlinic/internal/features/patients/domain"
)

func buildPatientListWhere(filter domain.PatientListFilter) (string, []any) {
	var clauses []string
	args := make([]any, 0, 6)
	argN := 1

	search := strings.TrimSpace(filter.Search)
	if search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		clauses = append(clauses, fmt.Sprintf(
			`(lower(first_name) LIKE $%d OR lower(last_name) LIKE $%d OR lower(coalesce(email, '')) LIKE $%d)`,
			argN, argN, argN,
		))
		args = append(args, pattern)
		argN++
	}

	if filter.Gender != "" {
		clauses = append(clauses, fmt.Sprintf(`gender = $%d`, argN))
		args = append(args, filter.Gender)
		argN++
	}

	if filter.DateOfBirthFrom != "" {
		clauses = append(clauses, fmt.Sprintf(`date_of_birth >= $%d`, argN))
		args = append(args, filter.DateOfBirthFrom)
		argN++
	}

	if filter.DateOfBirthTo != "" {
		clauses = append(clauses, fmt.Sprintf(`date_of_birth <= $%d`, argN))
		args = append(args, filter.DateOfBirthTo)
		argN++
	}

	if filter.RegisteredSince != "" {
		clauses = append(clauses, fmt.Sprintf(`created_at >= $%d`, argN))
		args = append(args, filter.RegisteredSince)
		argN++
	}

	if len(clauses) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(clauses, " AND "), args
}
