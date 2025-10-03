package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_SortSingleFieldDefault(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.SortBy(Sort("name")).Apply()

	expected := "select * from accounts ORDER BY name;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_SortSingleFieldDesc(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.SortBy(Sort("name", SortDesc)).Apply()

	expected := "select * from accounts ORDER BY name DESC;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_SortMultipleFields(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.SortBy(Sort("id"), Sort("name", SortDesc)).Apply()

	expected := "select * from accounts ORDER BY id, name DESC;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_SortWithWhere(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByIntColumn("id", 1)).SortBy(Sort("name")).Apply()

	expected := "select * from accounts WHERE id = %v ORDER BY name;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_SortWithLimitOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.SortBy(Sort("name", SortDesc)).Limit(10).Offset(5).Apply()

	expected := "select * from accounts ORDER BY name DESC LIMIT 10 OFFSET 5;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_CompleteQuery(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.
		Where(ByIntColumn("id", 1, 2, 3)).
		SortBy(Sort("name"), Sort("created_at", SortDesc)).
		Limit(10).
		Offset(5).
		Apply()

	expected := "select * from accounts WHERE id IN %v ORDER BY name, created_at DESC LIMIT 10 OFFSET 5;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}
