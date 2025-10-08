package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_IntSingleValue(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByIntColumn("id", []int{1})).Commit()

	expected := "select * from accounts WHERE id = $1;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}

func TestQueryBuilder_IntInClause(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByIntColumn("id", []int{1, 2, 3})).Commit()

	expected := "select * from accounts WHERE id IN ($1, $2, $3);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}

func TestQueryBuilder_IntSliceInClause(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	ids := []int{1, 2, 3}
	result, values := qb.Where(ByIntColumn("id", ids)).Commit()

	expected := "select * from accounts WHERE id IN ($1, $2, $3);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}
