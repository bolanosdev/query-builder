package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_Limit(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Limit(10).Apply()

	expected := "select * from accounts LIMIT 10;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_Offset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Offset(20).Apply()

	expected := "select * from accounts OFFSET 20;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_LimitAndOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Limit(10).Offset(20).Apply()

	expected := "select * from accounts LIMIT 10 OFFSET 20;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_WhereWithLimitAndOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByIntColumn("id", 1)).Limit(10).Offset(5).Apply()

	expected := "select * from accounts WHERE id = %v LIMIT 10 OFFSET 5;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_ComplexQueryWithLimitOffset(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(
		Or(
			ByIntColumn("id", 1),
			ByStringColumn("name", "carlos"),
		),
	).Limit(5).Offset(10).Apply()

	expected := "select * from accounts WHERE (id = %v OR name = %s) LIMIT 5 OFFSET 10;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}
