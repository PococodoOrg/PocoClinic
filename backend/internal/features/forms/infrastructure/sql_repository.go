package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dksch/pococlinic/internal/features/forms/domain"
	"github.com/dksch/pococlinic/internal/pkg/database"
)

type SQLRepository struct {
	pool *database.DB
}

func NewSQLRepository(pool *database.DB) *SQLRepository {
	return &SQLRepository{pool: pool}
}

func (r *SQLRepository) EnsureDefaultGroup(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO form_groups (id, name, sort_order)
		SELECT $1, 'General', 0
		WHERE NOT EXISTS (SELECT 1 FROM form_groups WHERE id = $1)
	`, domain.DefaultGroupID)
	return err
}

func (r *SQLRepository) CreateGroup(ctx context.Context, group *domain.FormGroup) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO form_groups (id, name, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, group.ID, group.Name, group.SortOrder, group.CreatedAt, group.UpdatedAt)
	return err
}

func (r *SQLRepository) UpdateGroup(ctx context.Context, group *domain.FormGroup) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE form_groups SET name = $2, sort_order = $3, updated_at = $4 WHERE id = $1
	`, group.ID, group.Name, group.SortOrder, group.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("form group not found")
	}
	return nil
}

func (r *SQLRepository) DeleteGroup(ctx context.Context, id string) error {
	if id == domain.DefaultGroupID.String() {
		return fmt.Errorf("cannot delete the default form group")
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE form_templates SET group_id = $1 WHERE group_id = $2
	`, domain.DefaultGroupID, id)
	if err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `DELETE FROM form_groups WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("form group not found")
	}
	return nil
}

func (r *SQLRepository) GetGroupByID(ctx context.Context, id string) (*domain.FormGroup, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, sort_order, created_at, updated_at FROM form_groups WHERE id = $1
	`, id)
	return scanGroup(row)
}

func (r *SQLRepository) ListGroups(ctx context.Context) ([]*domain.FormGroup, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, sort_order, created_at, updated_at FROM form_groups ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*domain.FormGroup, 0)
	for rows.Next() {
		group, err := scanGroup(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

func (r *SQLRepository) CreateTemplate(ctx context.Context, template *domain.FormTemplate) error {
	fields, err := json.Marshal(template.Fields)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO form_templates (id, group_id, name, form_type, fields, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, template.ID, template.GroupID, template.Name, template.FormType, fields, template.CreatedAt, template.UpdatedAt)
	return err
}

func (r *SQLRepository) UpdateTemplate(ctx context.Context, template *domain.FormTemplate) error {
	fields, err := json.Marshal(template.Fields)
	if err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE form_templates
		SET group_id = $2, name = $3, form_type = $4, fields = $5, updated_at = $6
		WHERE id = $1
	`, template.ID, template.GroupID, template.Name, template.FormType, fields, template.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("form template not found")
	}
	return nil
}

func (r *SQLRepository) DeleteTemplate(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM form_templates WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("form template not found")
	}
	return nil
}

func (r *SQLRepository) GetTemplateByID(ctx context.Context, id string) (*domain.FormTemplate, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, group_id, name, form_type, fields, created_at, updated_at
		FROM form_templates WHERE id = $1
	`, id)
	return scanTemplate(row)
}

func (r *SQLRepository) ListTemplates(ctx context.Context) ([]*domain.FormTemplate, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, group_id, name, form_type, fields, created_at, updated_at
		FROM form_templates ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]*domain.FormTemplate, 0)
	for rows.Next() {
		template, err := scanTemplate(rows)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, rows.Err()
}

func (r *SQLRepository) CreateSubmission(ctx context.Context, submission *domain.FormSubmission) error {
	answers, err := json.Marshal(submission.Answers)
	if err != nil {
		return err
	}
	fieldSnapshot, err := json.Marshal(submission.FieldSnapshot)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO form_submissions (
			id, entry_id, version, is_current, template_id, patient_id,
			submitted_by, updated_by, answers, field_snapshot, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, submission.ID, submission.EntryID, submission.Version, submission.IsCurrent,
		submission.TemplateID, submission.PatientID, submission.SubmittedBy, submission.UpdatedBy,
		answers, fieldSnapshot, submission.CreatedAt, submission.UpdatedAt)
	return err
}

func (r *SQLRepository) MarkEntryNotCurrent(ctx context.Context, entryID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE form_submissions SET is_current = 0 WHERE entry_id = $1 AND is_current = 1
	`, entryID)
	return err
}

func (r *SQLRepository) GetCurrentSubmissionByEntryID(ctx context.Context, entryID string) (*domain.FormSubmission, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT s.id, s.entry_id, s.version, s.is_current, s.template_id, t.name, t.form_type,
			s.patient_id, s.submitted_by, s.updated_by, s.answers, s.field_snapshot, s.created_at, s.updated_at
		FROM form_submissions s
		JOIN form_templates t ON t.id = s.template_id
		WHERE s.entry_id = $1 AND s.is_current = 1
	`, entryID)
	return scanSubmission(row)
}

func (r *SQLRepository) GetCurrentSingletonSubmission(ctx context.Context, patientID, templateID string) (*domain.FormSubmission, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT s.id, s.entry_id, s.version, s.is_current, s.template_id, t.name, t.form_type,
			s.patient_id, s.submitted_by, s.updated_by, s.answers, s.field_snapshot, s.created_at, s.updated_at
		FROM form_submissions s
		JOIN form_templates t ON t.id = s.template_id
		WHERE s.patient_id = $1 AND s.template_id = $2 AND s.is_current = 1 AND t.form_type = 'singleton'
	`, patientID, templateID)

	submission, err := scanSubmission(row)
	if errors.Is(err, database.ErrNoRows) {
		return nil, nil
	}
	return submission, err
}

func (r *SQLRepository) ListCurrentSubmissionsByPatient(ctx context.Context, patientID string) ([]*domain.FormSubmission, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.entry_id, s.version, s.is_current, s.template_id, t.name, t.form_type,
			s.patient_id, s.submitted_by, s.updated_by, s.answers, s.field_snapshot, s.created_at, s.updated_at
		FROM form_submissions s
		JOIN form_templates t ON t.id = s.template_id
		WHERE s.patient_id = $1 AND s.is_current = 1
		ORDER BY s.updated_at DESC
	`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*domain.FormSubmission, 0)
	for rows.Next() {
		submission, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, rows.Err()
}

func (r *SQLRepository) ListCurrentSubmissionsByTemplate(
	ctx context.Context,
	templateID string,
	filter domain.TemplateSubmissionFilter,
) (*domain.SubmissionReportPage, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM form_submissions s
		WHERE s.template_id = $1 AND s.is_current = 1
			AND ($2 IS NULL OR s.updated_at >= $2)
			AND ($3 IS NULL OR s.updated_at <= $3)
	`, templateID, filter.From, filter.To).Scan(&total)
	if err != nil {
		return nil, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.entry_id, s.version, s.is_current, s.template_id, t.name, t.form_type,
			s.patient_id, s.submitted_by, s.updated_by, s.answers, s.field_snapshot, s.created_at, s.updated_at,
			TRIM(p.first_name || ' ' || COALESCE(p.middle_name || ' ', '') || p.last_name) AS patient_name,
			u.name AS submitter_name
		FROM form_submissions s
		JOIN form_templates t ON t.id = s.template_id
		JOIN patients p ON p.id = s.patient_id
		JOIN users u ON u.id = s.submitted_by
		WHERE s.template_id = $1 AND s.is_current = 1
			AND ($2 IS NULL OR s.updated_at >= $2)
			AND ($3 IS NULL OR s.updated_at <= $3)
		ORDER BY s.updated_at DESC
		LIMIT $4 OFFSET $5
	`, templateID, filter.From, filter.To, filter.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.FormSubmissionReport, 0)
	for rows.Next() {
		report, err := scanSubmissionReport(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, report)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &domain.SubmissionReportPage{
		Items:    items,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (r *SQLRepository) ListSubmissionHistory(ctx context.Context, entryID string) ([]*domain.FormSubmission, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.entry_id, s.version, s.is_current, s.template_id, t.name, t.form_type,
			s.patient_id, s.submitted_by, s.updated_by, s.answers, s.field_snapshot, s.created_at, s.updated_at
		FROM form_submissions s
		JOIN form_templates t ON t.id = s.template_id
		WHERE s.entry_id = $1
		ORDER BY s.version DESC
	`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]*domain.FormSubmission, 0)
	for rows.Next() {
		submission, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, rows.Err()
}

func scanGroup(row database.Row) (*domain.FormGroup, error) {
	group := &domain.FormGroup{}
	err := row.Scan(&group.ID, &group.Name, &group.SortOrder, database.Time(&group.CreatedAt), database.Time(&group.UpdatedAt))
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, fmt.Errorf("form group not found")
		}
		return nil, err
	}
	return group, nil
}

func scanTemplate(row database.Row) (*domain.FormTemplate, error) {
	template := &domain.FormTemplate{}
	var fields []byte
	err := row.Scan(
		&template.ID, &template.GroupID, &template.Name, &template.FormType,
		&fields, database.Time(&template.CreatedAt), database.Time(&template.UpdatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, fmt.Errorf("form template not found")
		}
		return nil, err
	}
	if len(fields) > 0 {
		if err := json.Unmarshal(fields, &template.Fields); err != nil {
			return nil, err
		}
	}
	return template, nil
}

func scanSubmission(row database.Row) (*domain.FormSubmission, error) {
	submission := &domain.FormSubmission{}
	var answers []byte
	var fieldSnapshot []byte
	err := row.Scan(
		&submission.ID,
		&submission.EntryID,
		&submission.Version,
		&submission.IsCurrent,
		&submission.TemplateID,
		&submission.TemplateName,
		&submission.FormType,
		&submission.PatientID,
		&submission.SubmittedBy,
		&submission.UpdatedBy,
		&answers,
		&fieldSnapshot,
		database.Time(&submission.CreatedAt),
		database.Time(&submission.UpdatedAt),
	)
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return nil, fmt.Errorf("form entry not found")
		}
		return nil, err
	}
	if len(answers) > 0 {
		if err := json.Unmarshal(answers, &submission.Answers); err != nil {
			return nil, err
		}
	}
	if len(fieldSnapshot) > 0 {
		if err := json.Unmarshal(fieldSnapshot, &submission.FieldSnapshot); err != nil {
			return nil, err
		}
	}
	return submission, nil
}

func scanSubmissionReport(row database.Row) (*domain.FormSubmissionReport, error) {
	report := &domain.FormSubmissionReport{}
	var answers []byte
	var fieldSnapshot []byte
	err := row.Scan(
		&report.ID,
		&report.EntryID,
		&report.Version,
		&report.IsCurrent,
		&report.TemplateID,
		&report.TemplateName,
		&report.FormType,
		&report.PatientID,
		&report.SubmittedBy,
		&report.UpdatedBy,
		&answers,
		&fieldSnapshot,
		database.Time(&report.CreatedAt),
		database.Time(&report.UpdatedAt),
		&report.PatientName,
		&report.SubmittedByName,
	)
	if err != nil {
		return nil, err
	}
	if len(answers) > 0 {
		if err := json.Unmarshal(answers, &report.Answers); err != nil {
			return nil, err
		}
	}
	if len(fieldSnapshot) > 0 {
		if err := json.Unmarshal(fieldSnapshot, &report.FieldSnapshot); err != nil {
			return nil, err
		}
	}
	return report, nil
}

var _ domain.Repository = (*SQLRepository)(nil)
