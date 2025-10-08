package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_Limit(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Limit(10).Commit()

	expected := "select * from accounts LIMIT 10;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_Offset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Offset(20).Commit()

	// Offset without explicit limit applies default limit of 10
	expected := "select * from accounts LIMIT 10 OFFSET 20;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_LimitAndOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Limit(10).Offset(20).Commit()

	expected := "select * from accounts LIMIT 10 OFFSET 20;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_WhereWithLimitAndOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Where(ByIntColumn("id", []int{1})).Limit(10).Offset(5).Commit()

	expected := "select * from accounts WHERE id = $1 LIMIT 10 OFFSET 5;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_ComplexQueryWithLimitOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(
		Or(
			ByIntColumn("id", []int{1}),
			ByStringColumn("name", []string{"carlos"}),
		),
	).Limit(5).Offset(10).Commit()

	expected := "select * from accounts WHERE (id = $1 OR name = $2) LIMIT 5 OFFSET 10;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
}
