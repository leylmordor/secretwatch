package db

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/leylmordor/secretwatch/internal/store"
)

func openTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestSaveAndList(t *testing.T) {
	db := openTestDB(t)

	now := time.Now().UTC().Truncate(time.Second)
	remaining := 30
	secrets := []*store.Secret{
		{
			Name:            "prod-db-password",
			FullPath:        "arn:aws:secretsmanager:us-east-1:123:secret/prod-db-password",
			StoreType:       "aws",
			Account:         "123456789",
			LastRotated:     &now,
			MaxAgeDays:      90,
			Status:          store.StatusOK,
			DaysUntilExpiry: &remaining,
			Tags:            map[string]string{"env": "prod"},
		},
	}

	if err := db.Save(secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}

	all, err := db.List("", "", "", false)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(all))
	}

	got := all[0]
	if got.Name != "prod-db-password" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Account != "123456789" {
		t.Errorf("Account = %q", got.Account)
	}
	if got.Status != store.StatusOK {
		t.Errorf("Status = %q", got.Status)
	}
}

func TestListFilterByStatus(t *testing.T) {
	db := openTestDB(t)

	old := time.Now().AddDate(0, 0, -200)
	recent := time.Now().AddDate(0, 0, -5)
	overdueRemaining := -110
	okRemaining := 85

	secrets := []*store.Secret{
		{Name: "old-key", FullPath: "old", StoreType: "aws", LastRotated: &old,
			MaxAgeDays: 90, Status: store.StatusOverdue, DaysUntilExpiry: &overdueRemaining, Tags: map[string]string{}},
		{Name: "fresh-key", FullPath: "fresh", StoreType: "aws", LastRotated: &recent,
			MaxAgeDays: 90, Status: store.StatusOK, DaysUntilExpiry: &okRemaining, Tags: map[string]string{}},
	}
	if err := db.Save(secrets); err != nil {
		t.Fatal(err)
	}

	overdue, err := db.List("", "", "overdue", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(overdue) != 1 || overdue[0].Name != "old-key" {
		t.Errorf("expected 1 overdue secret named old-key, got %v", overdue)
	}
}

func TestSetExcluded(t *testing.T) {
	db := openTestDB(t)

	remaining := 60
	secrets := []*store.Secret{
		{Name: "secret", FullPath: "fp/secret", StoreType: "aws",
			MaxAgeDays: 90, Status: store.StatusOK, DaysUntilExpiry: &remaining, Tags: map[string]string{}},
	}
	if err := db.Save(secrets); err != nil {
		t.Fatal(err)
	}

	if err := db.SetExcluded("fp/secret", true); err != nil {
		t.Fatalf("SetExcluded: %v", err)
	}

	// default List hides excluded
	visible, err := db.List("", "", "", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(visible) != 0 {
		t.Errorf("expected excluded secret to be hidden, got %d", len(visible))
	}

	// showExcluded=true returns it
	all, err := db.List("", "", "", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || !all[0].Excluded {
		t.Errorf("expected excluded secret with Excluded=true")
	}
}
