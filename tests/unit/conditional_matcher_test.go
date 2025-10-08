package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_Basic(t *testing.T) {
	query := "select * from accounts"
	conditions := []QueryCondition{
		ByIntColumn("id", []int{1}),
		ByStringColumn("name", []string{"carlos"}, StringOpts{}),
	}

	qb := NewQueryBuilder(query)
	result, _ := qb.Where(conditions...).Commit()

	expected := "select * from accounts WHERE id = $1 AND name = $2;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_NoConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Commit()

	expected := query + ";"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_WhereCondition(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Where(ByIntColumn("id", []int{1})).Commit()

	expected := "select * from accounts WHERE id = $1;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_OrConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Where(Or(ByIntColumn("id", []int{1}), ByStringColumn("name", []string{"john"}, StringOpts{}))).Commit()

	expected := "select * from accounts WHERE (id = $1 OR name = $2);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_MultipleConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, _ := qb.Where(ByIntColumn("id", []int{1}), ByStringColumn("name", []string{"carlos"}, StringOpts{})).Commit()

	expected := "select * from accounts WHERE id = $1 AND name = $2;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_NestedConditions(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)

	result, _ := qb.Where(ByIntColumn("id", []int{1}), Or(ByIntColumn("id", []int{2}), ByStringColumn("name", []string{"carlos"}, StringOpts{}))).Commit()

	expected := "select * from accounts WHERE id = $1 AND (id = $2 OR name = $3);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}
