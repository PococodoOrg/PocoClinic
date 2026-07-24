package opshelper_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/PococodoOrg/PocoClinic/internal/pkg/backup"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/database"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/opshelper"
	"github.com/gin-gonic/gin"
)

func TestVerifyBackupRejectsInvalidFilename(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	opshelper.NewServer(nil, t.TempDir(), "", "test", "http://localhost:3000", "").RegisterRoutes(router)

	body := bytes.NewBufferString(`{"filename":"../etc/passwd"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/backups/verify", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestVerifyBackupValidArchive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	dataDir := t.TempDir()
	backupDir := filepath.Join(dataDir, "backups")
	docsDir := filepath.Join(dataDir, "documents")
	dbPath := filepath.Join(dataDir, "clinic.db")

	db, err := database.Connect(ctx, dbPath)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(ctx, `
		INSERT INTO patients (id, first_name, last_name, date_of_birth, gender, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, "p-verify", "Pat", "Verify", "1990-01-02", "other", "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("seed patient: %v", err)
	}

	archivePath, err := backup.Create(ctx, db, backupDir, docsDir, "test")
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}

	router := gin.New()
	opshelper.NewServer(db, backupDir, docsDir, "test", "http://localhost:3000", "").RegisterRoutes(router)

	payload, _ := json.Marshal(map[string]string{"filename": filepath.Base(archivePath)})
	req := httptest.NewRequest(http.MethodPost, "/api/backups/verify", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !result.Valid {
		t.Fatalf("expected valid backup, got message=%q", result.Message)
	}
}

func TestStatusIncludesMainAppOnline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	healthSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(healthSrv.Close)

	router := gin.New()
	opshelper.NewServer(nil, t.TempDir(), "", "test", "http://localhost:3000", healthSrv.URL).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var status struct {
		MainAppOnline bool `json:"mainAppOnline"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !status.MainAppOnline {
		t.Fatal("expected main app to be reported online")
	}
}
