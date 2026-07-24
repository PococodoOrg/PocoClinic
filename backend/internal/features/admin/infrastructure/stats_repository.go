package infrastructure

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	admindomain "github.com/PococodoOrg/PocoClinic/internal/features/admin/domain"
	authdomain "github.com/PococodoOrg/PocoClinic/internal/features/auth/domain"
	formdomain "github.com/PococodoOrg/PocoClinic/internal/features/forms/domain"
	patientdomain "github.com/PococodoOrg/PocoClinic/internal/features/patients/domain"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/backup"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/pathsafe"
)

// StatsRepository loads admin metrics from SQL or in-memory stores.
type StatsRepository struct {
	pool          *database.DB
	patientRepo   patientdomain.PatientRepository
	userRepo      authdomain.UserRepository
	formRepo      formdomain.Repository
	backupDir     string
	documentsDir  string
	startedAt     time.Time
	environment   string
	appVersion    string
}

func NewStatsRepository(
	pool *database.DB,
	patientRepo patientdomain.PatientRepository,
	userRepo authdomain.UserRepository,
	formRepo formdomain.Repository,
	backupDir, documentsDir string,
	startedAt time.Time,
	environment, appVersion string,
) *StatsRepository {
	return &StatsRepository{
		pool:         pool,
		patientRepo:  patientRepo,
		userRepo:     userRepo,
		formRepo:     formRepo,
		backupDir:    backupDir,
		documentsDir: documentsDir,
		startedAt:    startedAt,
		environment:  environment,
		appVersion:   appVersion,
	}
}

func (r *StatsRepository) GetSystemStatus(ctx context.Context) (*admindomain.SystemStatus, error) {
	status := &admindomain.SystemStatus{
		StorageMode:       "memory",
		DatabaseConnected: false,
		PendingMigrations: []string{},
		UptimeSeconds:     int64(time.Since(r.startedAt).Seconds()),
		Environment:       r.environment,
		AppVersion:        r.appVersion,
		Backup:            backupStatusFromDir(r.backupDir),
		Security: admindomain.SecurityOverview{
			DocumentsStorageReady: true,
		},
	}

	if r.pool != nil {
		status.StorageMode = "database"
		if err := database.Ping(ctx, r.pool); err == nil {
			status.DatabaseConnected = true
		}

		pending, err := database.PendingMigrations(ctx, r.pool)
		if err != nil {
			return nil, fmt.Errorf("list pending migrations: %w", err)
		}
		status.PendingMigrations = pending

		applied, err := database.AppliedMigrations(ctx, r.pool)
		if err != nil {
			return nil, fmt.Errorf("list applied migrations: %w", err)
		}
		if len(applied) > 0 {
			status.LatestMigration = applied[len(applied)-1]
		}

		counts, err := r.countsFromDatabase(ctx)
		if err != nil {
			return nil, err
		}
		status.Counts = counts
		return status, nil
	}

	counts, err := r.countsFromMemory(ctx)
	if err != nil {
		return nil, err
	}
	status.Counts = counts
	return status, nil
}

func (r *StatsRepository) GetPatientCensus(ctx context.Context) (*admindomain.PatientCensus, error) {
	if r.pool != nil {
		return r.patientCensusFromDatabase(ctx)
	}
	return r.patientCensusFromMemory(ctx)
}

func (r *StatsRepository) countsFromDatabase(ctx context.Context) (admindomain.ResourceCounts, error) {
	counts := admindomain.ResourceCounts{}
	queries := []struct {
		sql   string
		field *int64
	}{
		{`SELECT COUNT(*) FROM patients`, &counts.Patients},
		{`SELECT COUNT(*) FROM users`, &counts.Users},
		{`SELECT COUNT(*) FROM form_templates`, &counts.FormTemplates},
		{`SELECT COUNT(*) FROM form_submissions WHERE is_current = 1`, &counts.FormSubmissions},
		{`SELECT COUNT(*) FROM sessions WHERE expires_at > $1`, &counts.ActiveSessions},
		{`SELECT COUNT(*) FROM users WHERE locked_until IS NOT NULL AND locked_until > $1`, &counts.LockedAccounts},
		{`SELECT COUNT(*) FROM users WHERE must_change_pin = 1`, &counts.DefaultPinAccounts},
		{`SELECT COUNT(*) FROM patient_documents`, &counts.PatientDocuments},
	}

	now := time.Now().UTC()
	for _, query := range queries {
		var err error
		if strings.Contains(query.sql, "$1") {
			err = r.pool.QueryRow(ctx, query.sql, now).Scan(query.field)
		} else {
			err = r.pool.QueryRow(ctx, query.sql).Scan(query.field)
		}
		if err != nil {
			return counts, err
		}
	}
	return counts, nil
}

func (r *StatsRepository) countsFromMemory(ctx context.Context) (admindomain.ResourceCounts, error) {
	counts := admindomain.ResourceCounts{}

	_, totalPatients, err := r.patientRepo.ListPaginated(ctx, 1, 1, patientdomain.PatientListFilter{})
	if err != nil {
		return counts, err
	}
	counts.Patients = totalPatients

	_, totalUsers, err := r.userRepo.ListPaginated(ctx, 1, 1, "")
	if err != nil {
		return counts, err
	}
	counts.Users = totalUsers

	templates, err := r.formRepo.ListTemplates(ctx)
	if err != nil {
		return counts, err
	}
	counts.FormTemplates = int64(len(templates))

	return counts, nil
}

func (r *StatsRepository) patientCensusFromDatabase(ctx context.Context) (*admindomain.PatientCensus, error) {
	census := &admindomain.PatientCensus{ByGender: []admindomain.GenderCount{}}

	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM patients`).Scan(&census.TotalPatients); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT gender, COUNT(*) FROM patients GROUP BY gender ORDER BY gender
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item admindomain.GenderCount
		if err := rows.Scan(&item.Gender, &item.Count); err != nil {
			return nil, err
		}
		census.ByGender = append(census.ByGender, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM patients WHERE created_at >= $1
	`, time.Now().UTC().AddDate(0, 0, -30)).Scan(&census.AddedLast30Days); err != nil {
		return nil, err
	}

	return census, nil
}

func (r *StatsRepository) patientCensusFromMemory(ctx context.Context) (*admindomain.PatientCensus, error) {
	patients, err := r.patientRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	census := &admindomain.PatientCensus{
		TotalPatients: int64(len(patients)),
		ByGender:      []admindomain.GenderCount{},
	}

	genderCounts := map[string]int64{}
	cutoff := time.Now().AddDate(0, 0, -30)
	for _, patient := range patients {
		genderCounts[string(patient.Gender)]++
		if patient.CreatedAt.After(cutoff) {
			census.AddedLast30Days++
		}
	}

	for gender, count := range genderCounts {
		census.ByGender = append(census.ByGender, admindomain.GenderCount{
			Gender: gender,
			Count:  count,
		})
	}

	return census, nil
}

func backupStatusFromDir(dir string) admindomain.BackupStatus {
	latest, err := backup.FindLatest(dir)
	if err != nil || latest == nil {
		return admindomain.BackupStatus{Status: "unknown"}
	}

	status := admindomain.BackupStatus{
		LastBackupAt:   &latest.CreatedAt,
		LastBackupFile: latest.Filename,
	}

	age := time.Since(latest.CreatedAt).Hours()
	status.BackupAgeHours = &age
	switch {
	case age <= 24:
		status.Status = "ok"
	case age <= 48:
		status.Status = "warning"
	default:
		status.Status = "critical"
	}
	return status
}

func (r *StatsRepository) GetActivitySummary(ctx context.Context, periodDays int) (*admindomain.ActivitySummary, error) {
	if periodDays <= 0 {
		periodDays = 7
	}
	if periodDays > 90 {
		periodDays = 90
	}

	summary := &admindomain.ActivitySummary{PeriodDays: periodDays}
	if r.pool == nil {
		return summary, nil
	}

	queries := []struct {
		sql   string
		field *int64
	}{
		{
			`SELECT COUNT(*) FROM audit_logs WHERE event_type = 'auth.login.success' AND created_at >= $1`,
			&summary.SignIns,
		},
		{
			`SELECT COUNT(*) FROM audit_logs WHERE event_type = 'auth.login.failure' AND created_at >= $1`,
			&summary.FailedSignIns,
		},
		{
			`SELECT COUNT(*) FROM patients WHERE created_at >= $1`,
			&summary.PatientsCreated,
		},
		{
			`SELECT COUNT(*) FROM patient_notes WHERE created_at >= $1`,
			&summary.ChartNotesAdded,
		},
		{
			`SELECT COUNT(*) FROM audit_logs WHERE event_type = 'patient.viewed' AND created_at >= $1`,
			&summary.PatientViews,
		},
		{
			`SELECT COUNT(*) FROM audit_logs WHERE event_type = 'patient.document.uploaded' AND created_at >= $1`,
			&summary.DocumentsUploaded,
		},
		{
			`SELECT COUNT(*) FROM audit_logs WHERE event_type = 'system.backup.created' AND success = 1 AND created_at >= $1`,
			&summary.BackupsCreated,
		},
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -periodDays)
	for _, query := range queries {
		if err := r.pool.QueryRow(ctx, query.sql, cutoff).Scan(query.field); err != nil {
			return nil, err
		}
	}

	return summary, nil
}

func (r *StatsRepository) GetStaffActivity(ctx context.Context, periodDays int) (*admindomain.StaffActivityReport, error) {
	if periodDays <= 0 {
		periodDays = 30
	}
	if periodDays > 90 {
		periodDays = 90
	}

	report := &admindomain.StaffActivityReport{PeriodDays: periodDays}
	if r.pool == nil {
		return report, nil
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -periodDays)
	rows, err := r.pool.Query(ctx, `
		SELECT
			u.id,
			u.name,
			u.email,
			COALESCE(s.sign_ins, 0),
			COALESCE(f.failed_sign_ins, 0),
			COALESCE(v.patient_views, 0),
			ls.last_sign_in
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) AS sign_ins
			FROM audit_logs
			WHERE event_type = 'auth.login.success'
			  AND created_at >= $1
			GROUP BY user_id
		) s ON s.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) AS failed_sign_ins
			FROM audit_logs
			WHERE event_type = 'auth.login.failure'
			  AND created_at >= $1
			GROUP BY user_id
		) f ON f.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) AS patient_views
			FROM audit_logs
			WHERE event_type = 'patient.viewed'
			  AND created_at >= $1
			GROUP BY user_id
		) v ON v.user_id = u.id
		LEFT JOIN (
			SELECT user_id, MAX(created_at) AS last_sign_in
			FROM audit_logs
			WHERE event_type = 'auth.login.success'
			  AND created_at >= $1
			GROUP BY user_id
		) ls ON ls.user_id = u.id
		ORDER BY COALESCE(s.sign_ins, 0) DESC, u.name ASC
	`, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row := admindomain.StaffActivityRow{}
		var lastSignIn *time.Time
		if err := rows.Scan(
			&row.UserID,
			&row.Name,
			&row.Email,
			&row.SignIns,
			&row.FailedSignIns,
			&row.PatientViews,
			database.NullTime(&lastSignIn),
		); err != nil {
			return nil, err
		}
		row.LastSignIn = lastSignIn
		report.Staff = append(report.Staff, row)
	}
	return report, rows.Err()
}

func documentsStorageReady(dir string) bool {
	if strings.TrimSpace(dir) == "" {
		return false
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return false
	}
	probe, err := pathsafe.JoinRoot(dir, ".write-probe")
	if err != nil {
		return false
	}
	if err := os.WriteFile(probe, []byte("ok"), 0o640); err != nil {
		return false
	}
	_ = os.Remove(probe)
	return true
}

var _ admindomain.StatsReader = (*StatsRepository)(nil)
