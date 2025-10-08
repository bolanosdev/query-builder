package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_StringEqual(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	names := []string{"carlos"}

	result, values := qb.Where(ByStringColumn("name", names)).Commit()

	expected := "select * from accounts WHERE name = $1;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
	if values[0] != "carlos" {
		t.Errorf("Expected [carlos], got %v", values)
	}
}

func TestQueryBuilder_StringIn(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	names := []string{"carlos", "john"}

	result, values := qb.Where(ByStringColumn("name", names)).Commit()

	expected := "select * from accounts WHERE name IN ($1, $2);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
	if values[0] != "carlos" || values[1] != "john" {
		t.Errorf("Expected [carlos, john], got %v", values)
	}
}

func TestQueryBuilder_StringContains(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"car"}, StringContains)).Commit()

	expected := "select * from accounts WHERE name LIKE '%' || $1 || '%';"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "car" {
		t.Errorf("Expected value '%%car%%', got %v", values)
	}
}

func TestQueryBuilder_StringStartsWith(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"car"}, StringStartsWith)).Commit()

	expected := "select * from accounts WHERE name LIKE $1 || '%';"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "car" {
		t.Errorf("Expected value 'car%%', got %v", values)
	}
}

func TestQueryBuilder_StringEndsWith(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"los"}, StringEndsWith)).Commit()

	expected := "select * from accounts WHERE name LIKE '%' || $1;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "los" {
		t.Errorf("Expected value '%%los', got %v", values)
	}
}

func TestQueryBuilder_StringCaseInsensitive(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"Carlos"}, StringExact, NonSensitive)).Commit()

	expected := "select * from accounts WHERE LOWER(name) = LOWER($1);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "Carlos" {
		t.Errorf("Expected value 'Carlos', got %v", values)
	}
}

func TestQueryBuilder_StringContainsCaseInsensitive(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"CAR"}, StringContains, NonSensitive)).Commit()

	expected := "select * from accounts WHERE LOWER(name) LIKE '%' || LOWER($1) || '%';"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "CAR" {
		t.Errorf("Expected value 'CAR', got %v", values)
	}
}

func TestQueryBuilder_StringStartsWithCaseInsensitive(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"CAR"}, StringStartsWith, NonSensitive)).Commit()

	expected := "select * from accounts WHERE LOWER(name) LIKE LOWER($1) || '%';"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "CAR" {
		t.Errorf("Expected value 'CAR', got %v", values)
	}
}

func TestQueryBuilder_StringEndsWithCaseInsensitive(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result, values := qb.Where(ByStringColumn("name", []string{"LOS"}, StringEndsWith, NonSensitive)).Commit()

	expected := "select * from accounts WHERE LOWER(name) LIKE '%' || LOWER($1);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 || values[0] != "LOS" {
		t.Errorf("Expected value 'LOS', got %v", values)
	}
}
