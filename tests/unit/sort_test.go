package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_SortSingleFieldDefault(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.SortBy(Sort("name")).Commit()

	expected := "select * from accounts ORDER BY name;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 0 {
		t.Errorf("Expected 0 values, got %d", len(values))
	}
}

func TestQueryBuilder_SortSingleFieldDesc(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.SortBy(Sort("name", SortDesc)).Commit()

	expected := "select * from accounts ORDER BY name DESC;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 0 {
		t.Errorf("Expected 0 values, got %d", len(values))
	}
}

func TestQueryBuilder_SortMultipleFields(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.SortBy(Sort("id"), Sort("name", SortDesc)).Commit()

	expected := "select * from accounts ORDER BY id, name DESC;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 0 {
		t.Errorf("Expected 0 values, got %d", len(values))
	}
}

func TestQueryBuilder_SortWithWhere(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByIntColumn("id", []int{1})).SortBy(Sort("name")).Commit()

	expected := "select * from accounts WHERE id = $1 ORDER BY name;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}

func TestQueryBuilder_SortWithLimitOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.SortBy(Sort("name", SortDesc)).Limit(10).Offset(5).Commit()

	expected := "select * from accounts ORDER BY name DESC LIMIT 10 OFFSET 5;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 0 {
		t.Errorf("Expected 0 values, got %d", len(values))
	}
}

func TestQueryBuilder_CompleteQuery(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.
		Where(ByIntColumn("id", []int{1, 2, 3})).
		SortBy(Sort("name"), Sort("created_at", SortDesc)).
		Limit(10).
		Offset(5).
		Commit()

	expected := "select * from accounts WHERE id IN ($1, $2, $3) ORDER BY name, created_at DESC LIMIT 10 OFFSET 5;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}
