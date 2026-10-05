//go:build examples

package authors

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()

	// DuckDB runs in the process: an empty data source is a database in
	// memory that lives as long as the connection.
	sdb, err := sql.Open("duckdb", "")
	if err != nil {
		t.Fatal(err)
	}
	defer sdb.Close()

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sdb.ExecContext(ctx, string(schema)); err != nil {
		t.Fatal(err)
	}

	db := New(sdb)

	// list all authors
	authors, err := db.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authors)

	// create an author
	insertedAuthor, err := db.CreateAuthor(ctx, CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := db.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedAuthor.Name != "Brian Kernighan" || !fetchedAuthor.Bio.Valid {
		t.Fatalf("unexpected author: %+v", fetchedAuthor)
	}
	if len(fetchedAuthor.DataJSON) != 0 || len(fetchedAuthor.PreferencesJSON) != 0 {
		t.Fatalf("unexpected JSON defaults: %+v", fetchedAuthor)
	}
	t.Log(fetchedAuthor)

	authors, err = db.FindAuthorsByName(ctx, "Brian Kernighan")
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 || authors[0].ID != insertedAuthor.ID {
		t.Fatalf("unexpected authors: %+v", authors)
	}

	// batch create authors
	batchAuthors := []CreateAuthorsParams{
		{
			Name: "Dennis Ritchie",
			Bio:  sql.NullString{String: "Creator of C", Valid: true},
		},
		{
			Name: "Ken Thompson",
			Bio:  sql.NullString{String: "Co-creator of Unix", Valid: true},
		},
	}
	var callbackIndexes []int
	db.CreateAuthors(ctx, batchAuthors).Exec(func(i int, err error) {
		if err != nil {
			t.Fatalf("batch create author %d: %v", i, err)
		}
		callbackIndexes = append(callbackIndexes, i)
	})
	if len(callbackIndexes) != len(batchAuthors) {
		t.Fatalf("expected %d batch callbacks, got %d", len(batchAuthors), len(callbackIndexes))
	}
	for i, callbackIndex := range callbackIndexes {
		if callbackIndex != i {
			t.Fatalf("callback %d had index %d", i, callbackIndex)
		}
	}

	authors, err = db.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != len(batchAuthors)+1 {
		t.Fatalf("expected %d authors, got %+v", len(batchAuthors)+1, authors)
	}
	for _, batchAuthor := range batchAuthors {
		matched := false
		for _, author := range authors {
			if author.Name == batchAuthor.Name && author.Bio == batchAuthor.Bio {
				matched = true
				break
			}
		}
		if !matched {
			t.Fatalf("missing batch-created author: %+v", batchAuthor)
		}
	}

	tx, err := sdb.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	db.WithTx(tx).CreateAuthors(ctx, []CreateAuthorsParams{{Name: "Rolled Back"}}).Exec(func(i int, err error) {
		if err != nil {
			t.Fatalf("transaction batch create author %d: %v", i, err)
		}
	})
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	authors, err = db.FindAuthorsByName(ctx, "Rolled Back")
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 0 {
		t.Fatalf("batch insert survived rollback: %+v", authors)
	}

	closedBatch := db.CreateAuthors(ctx, []CreateAuthorsParams{{Name: "Closed Batch"}})
	if err := closedBatch.Close(); err != nil {
		t.Fatal(err)
	}
	closedBatch.Exec(func(i int, err error) {
		if !errors.Is(err, ErrBatchAlreadyClosed) {
			t.Fatalf("closed batch entry %d returned %v", i, err)
		}
	})
	authors, err = db.FindAuthorsByName(ctx, "Closed Batch")
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 0 {
		t.Fatalf("closed batch inserted authors: %+v", authors)
	}

	// delete the author
	if err := db.DeleteAuthor(ctx, insertedAuthor.ID); err != nil {
		t.Fatal(err)
	}
}
