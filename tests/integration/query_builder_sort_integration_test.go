package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_Integration_NoSort(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeSortQuery(t, db)
	ids := mapUserIDs(u)

	// With 50 users, verify first 5 IDs are in natural order
	require.Equal(t, 50, len(ids))
	require.Equal(t, []int{1, 2, 3, 4, 5}, ids[:5])
}

func TestQueryBuilder_Integration_SortByIdDesc(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeSortQuery(t, db, Sort("id", SortDesc))
	ids := mapUserIDs(u)

	// Verify first 5 IDs are in descending order
	require.Equal(t, 50, len(ids))
	require.Equal(t, []int{50, 49, 48, 47, 46}, ids[:5])
}

func TestQueryBuilder_Integration_SortByNameAsc(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeSortQuery(t, db, Sort("name", SortAsc))
	ids := mapUserIDs(u)

	// Verify first 5 names alphabetically (alice, anna, bob, brian, carlos)
	require.Equal(t, 50, len(ids))
	require.Equal(t, []int{4, 30, 5, 31, 1}, ids[:5])
}

func TestQueryBuilder_Integration_SortByNameDesc(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeSortQuery(t, db, Sort("name", SortDesc))
	ids := mapUserIDs(u)

	// Verify first 5 names in reverse alphabetical order (zack, yara, xavier, wendy, victor)
	require.Equal(t, 50, len(ids))
	require.Equal(t, []int{29, 28, 27, 26, 25}, ids[:5])
}

func TestQueryBuilder_Integration_SortMultiple(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeSortQuery(t, db, Sort("name"), Sort("id", SortDesc))
	ids := mapUserIDs(u)

	// Verify first 5 sorted by name ASC, then id DESC (alice, anna, bob, brian, carlos)
	require.Equal(t, 50, len(ids))
	require.Equal(t, []int{4, 30, 5, 31, 1}, ids[:5])
}
