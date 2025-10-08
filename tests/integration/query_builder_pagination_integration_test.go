package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_Integration_Offset(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.Offset(1).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	u, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	// Offset without explicit limit applies default limit of 10
	ids := mapUserIDs(u)
	require.Equal(t, 10, len(ids))
	require.Equal(t, []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, ids)
}

func TestQueryBuilder_Integration_Limit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.Limit(2).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	u, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	ids := mapUserIDs(u)
	require.Equal(t, []int{1, 2}, ids)
}

func TestQueryBuilder_Integration_LimitOffset(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.Limit(2).Offset(2).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	u, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	ids := mapUserIDs(u)
	require.Equal(t, []int{3, 4}, ids)
}

func TestQueryBuilder_Integration_CompleteQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.Where(ByIntColumn("id", []int{1, 2, 3, 4, 5})).
		SortBy(Sort("name", SortDesc)).Limit(3).Offset(1).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	u, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	require.Equal(t, len(u), 3)
	require.Equal(t, u[0].ID, 3)
	require.Equal(t, u[0].Name, "jane")
	require.Equal(t, u[1].ID, 1)
	require.Equal(t, u[1].Name, "carlos")
	require.Equal(t, u[2].ID, 5)
	require.Equal(t, u[2].Name, "bob")
}
