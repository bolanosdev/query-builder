package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_Integration_StringExact(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"carlos"}))
	ids := mapUserIDs(u)

	require.Equal(t, []int{1}, ids)
}

func TestQueryBuilder_Integration_StringExactMultiple(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"carlos", "john", "alice"}))
	ids := mapUserIDs(u)

	require.Equal(t, []int{1, 2, 4}, ids)
}

func TestQueryBuilder_Integration_StringContains(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Find all names containing "ar" (carlos, charlie, maria, yara, clara, karen, oscar)
	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"ar"}, StringContains))
	ids := mapUserIDs(u)

	require.Equal(t, []int{1, 6, 16, 28, 32, 40, 44}, ids)
}

func TestQueryBuilder_Integration_StringStartsWith(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Find names starting with "j" (john, jane, jack, james)
	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"j"}, StringStartsWith))
	ids := mapUserIDs(u)

	require.Equal(t, []int{2, 3, 13, 39}, ids)
}

func TestQueryBuilder_Integration_StringEndsWith(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Find names ending with "n" (john, nathan, quinn, brian, karen, quentin)
	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"n"}, StringEndsWith))
	ids := mapUserIDs(u)

	require.Equal(t, []int{2, 17, 20, 31, 40, 46}, ids)
}

func TestQueryBuilder_Integration_StringCaseInsensitiveExact(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"CARLOS"}, StringExact, NonSensitive))
	ids := mapUserIDs(u)

	require.Equal(t, []int{1}, ids)
}

func TestQueryBuilder_Integration_StringCaseInsensitiveStartsWith(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"J"}, StringStartsWith, NonSensitive))
	ids := mapUserIDs(u)

	require.Equal(t, []int{2, 3, 13, 39}, ids)
}

func TestQueryBuilder_Integration_StringCaseInsensitiveEndsWith(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByStringColumn("name", []string{"N"}, StringEndsWith, NonSensitive))
	ids := mapUserIDs(u)

	require.Equal(t, []int{2, 17, 20, 31, 40, 46}, ids)
}
