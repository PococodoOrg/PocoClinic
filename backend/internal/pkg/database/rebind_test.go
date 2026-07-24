package database

import (
	"context"
	"path/filepath"
	"testing"
)

func TestRebindDoesNotMangleDoubleDigitPlaceholders(t *testing.T) {
	in := `INSERT INTO users (a,b,c,d,e,f,g,h,i,j,k) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	got := rebind(in)
	want := `INSERT INTO users (a,b,c,d,e,f,g,h,i,j,k) VALUES (?1,?2,?3,?4,?5,?6,?7,?8,?9,?10,?11)`
	if got != want {
		t.Fatalf("rebind mangled placeholders:\n got: %s\nwant: %s", got, want)
	}
}

func TestRebindPreservesReusedPlaceholders(t *testing.T) {
	in := `(lower(first_name) LIKE $1 OR lower(last_name) LIKE $1 OR lower(coalesce(email, '')) LIKE $1)`
	got := rebind(in)
	want := `(lower(first_name) LIKE ?1 OR lower(last_name) LIKE ?1 OR lower(coalesce(email, '')) LIKE ?1)`
	if got != want {
		t.Fatalf("rebind lost reuse:\n got: %s\nwant: %s", got, want)
	}
}

func TestRebindLeavesPlainSQL(t *testing.T) {
	in := `SELECT COUNT(*) FROM users`
	if got := rebind(in); got != in {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestRebindReuseExecutesWithSingleArg(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "reuse.db")
	ctx := context.Background()
	db, err := Connect(ctx, path)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(ctx, `
		INSERT INTO patients (
			id, first_name, last_name, date_of_birth, gender, email, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, "p1", "Ada", "Lovelace", "1815-12-10", "female", "ada@example.com", "2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	var n int
	err = db.QueryRow(ctx, `
		SELECT COUNT(*) FROM patients
		WHERE lower(first_name) LIKE $1 OR lower(last_name) LIKE $1 OR lower(coalesce(email, '')) LIKE $1
	`, "%ada%").Scan(&n)
	if err != nil {
		t.Fatalf("reuse query: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 row, got %d", n)
	}
}

func TestHighPlaceholderInsert(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "high.db")
	ctx := context.Background()
	db, err := Connect(ctx, path)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(ctx, `
		INSERT INTO users (
			id, email, name, role, key_hash, key_salt, pin_hash, pin_salt, key_lookup,
			failed_attempts, locked_until, last_login, must_change_pin, is_active, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`,
		"u1", "a@b.c", "Admin", "admin",
		[]byte{1}, []byte{2}, []byte{3}, []byte{4}, nil,
		0, nil, nil, 1, 1, "2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("16-arg insert (would fail with $10->$1 mangling): %v", err)
	}
}
