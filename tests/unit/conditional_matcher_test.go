package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_Basic(t *testing.T) {
	query := "select * from accounts"
	conditions := []QueryCondition{
		ByIntColumn("id", 1),
		ByStringColumn("name", "carlos"),
	}

	qb := NewQueryBuilder(query)
	result := qb.Where(conditions...).Apply()

	expected := "select * from accounts WHERE id = %v AND name = %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_NoConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Apply()

	expected := query + ";"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_WhereCondition(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByIntColumn("id", 1)).Apply()

	expected := "select * from accounts WHERE id = %v;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_OrConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(Or(ByIntColumn("id", 1), ByStringColumn("name", "john"))).Apply()

	expected := "select * from accounts WHERE (id = %v OR name = %s);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_MultipleConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByIntColumn("id", 1), ByStringColumn("name", "carlos")).Apply()

	expected := "select * from accounts WHERE id = %v AND name = %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_NestedConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)

	result := qb.Where(ByIntColumn("id", 1), Or(ByIntColumn("id", 2), ByStringColumn("name", "carlos"))).Apply()

	expected := "select * from accounts WHERE id = %v AND (id = %v OR name = %s);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}
