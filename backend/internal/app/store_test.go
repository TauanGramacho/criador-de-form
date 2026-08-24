package app

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestDeleteFormRequiresOwnerAndRemovesResponses(t *testing.T) {
	db, err := OpenDB("file:" + filepath.ToSlash(filepath.Join(t.TempDir(), "test.db")))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}

	store := NewStore(db)
	owner, err := store.CreateUserWithPassword("Ana", "ana@example.com", "hash")
	if err != nil {
		t.Fatal(err)
	}
	other, err := store.CreateUserWithPassword("Bia", "bia@example.com", "hash")
	if err != nil {
		t.Fatal(err)
	}

	form, err := store.CreateForm(owner.ID, "Cafe", "", []FormField{{ID: "name", Label: "Nome", Type: FieldTypeText}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateResponse(form.ID, map[string]any{"name": "Lia"}, "127.0.0.1"); err != nil {
		t.Fatal(err)
	}

	if err := store.DeleteForm(other.ID, form.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found for another owner, got %v", err)
	}
	if _, err := store.GetForm(owner.ID, form.ID); err != nil {
		t.Fatalf("expected owner form to remain, got %v", err)
	}

	if err := store.DeleteForm(owner.ID, form.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetForm(owner.ID, form.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted form to be missing, got %v", err)
	}
	responses, err := store.ListResponses(owner.ID, form.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(responses) != 0 {
		t.Fatalf("expected no responses after deleting form, got %d", len(responses))
	}
}
