package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_Integration_AndSingleCondition(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByIntColumn("id", []int{1}))
	ids := mapUserIDs(u)

	require.Equal(t, []int{1}, ids)
}

func TestQueryBuilder_Integration_AndMultipleConditions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db,
		ByIntColumn("id", []int{1}),
		ByStringColumn("name", []string{"carlos"}),
	)
	ids := mapUserIDs(u)

	require.Equal(t, []int{1}, ids)
}

func TestQueryBuilder_Integration_ORSingleConditions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db,
		Or(ByIntColumn("id", []int{1}), ByStringColumn("name", []string{"john"})),
	)
	ids := mapUserIDs(u)

	require.Equal(t, []int{1, 2}, ids)
}

func TestQueryBuilder_Integration_ORMultipleConditions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db,
		Or(ByIntColumn("id", []int{1}), ByStringColumn("name", []string{"john", "alice"})),
	)
	ids := mapUserIDs(u)

	require.Equal(t, []int{1, 2, 4}, ids)
}

func TestQueryBuilder_Integration_ORConditionsSingleResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db,
		Or(ByIntColumn("id", []int{1}), ByStringColumn("name", []string{"carlos"})),
	)
	ids := mapUserIDs(u)

	require.Equal(t, []int{1}, ids)
}

func TestQueryBuilder_Integration_AndMultipleClause(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db, ByIntColumn("id", []int{1, 2}))
	ids := mapUserIDs(u)

	require.Equal(t, []int{1, 2}, ids)

	u = executeWhereQuery(t, db, ByStringColumn("name", []string{"carlos", "john"}))
	ids = mapUserIDs(u)

	require.Equal(t, []int{1, 2}, ids)
}

func TestQueryBuilder_Integration_NestedConditions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	u := executeWhereQuery(t, db,
		ByIntColumn("id", []int{1}),
		Or(ByIntColumn("id", []int{2}), ByStringColumn("name", []string{"carlos"})),
	)
	ids := mapUserIDs(u)

	require.Equal(t, []int{1}, ids)
}
