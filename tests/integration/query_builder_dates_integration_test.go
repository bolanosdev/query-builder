package querybuilder_test

import (
	"testing"
	"time"

	. "github.com/bolanosdev/query-builder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_Integration_DateAfter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	date := time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC)
	u := executeDateQuery(t, db, "created_at", Dates{After: date})
	ids := mapUserIDs(u)

	// All users after ID 2 (48 users from ID 3 to 50)
	require.Equal(t, 48, len(ids))
	require.Equal(t, []int{3, 4, 5}, ids[:3])
}

func TestQueryBuilder_Integration_DateOn(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	date := time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC)
	u := executeDateQuery(t, db, "created_at", Dates{On: date})
	ids := mapUserIDs(u)

	require.Equal(t, []int{2}, ids)
}

func TestQueryBuilder_Integration_DateBefore(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Query for dates before 2024-01-03 (should return carlos and john)
	date := time.Date(2024, 1, 2, 23, 59, 59, 0, time.UTC)
	u := executeDateQuery(t, db, "created_at", Dates{Before: date})
	ids := mapUserIDs(u)

	require.Equal(t, []int{1, 2}, ids)
}

func TestQueryBuilder_Integration_DateBetween(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Query for dates after 2024-01-02 11:00 and before 2024-01-04 (should return only jane)
	afterDate := time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC)
	beforeDate := time.Date(2024, 1, 3, 23, 59, 59, 0, time.UTC)
	u := executeDateQuery(t, db, "created_at", Dates{After: afterDate, Before: beforeDate})
	ids := mapUserIDs(u)

	require.Equal(t, []int{3}, ids)
}
