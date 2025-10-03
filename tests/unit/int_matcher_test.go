package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_IntSingleValue(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByIntColumn("id", 1)).Apply()

	expected := "select * from accounts WHERE id = %v;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_IntInClause(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByIntColumn("id", 1, 2, 3)).Apply()

	expected := "select * from accounts WHERE id IN %v;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}
