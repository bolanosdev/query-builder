package querybuilder_test

import (
	"testing"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_StringContains(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByStringColumn("name", "car", StringContains)).Apply()

	expected := "select * from accounts WHERE name LIKE %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 || values[0] != "%car%" {
		t.Errorf("Expected value '%%car%%', got %v", values)
	}
}

func TestQueryBuilder_StringStartsWith(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByStringColumn("name", "car", StringStartsWith)).Apply()

	expected := "select * from accounts WHERE name LIKE %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 || values[0] != "car%" {
		t.Errorf("Expected value 'car%%', got %v", values)
	}
}

func TestQueryBuilder_StringEndsWith(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByStringColumn("name", "los", StringEndsWith)).Apply()

	expected := "select * from accounts WHERE name LIKE %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 || values[0] != "%los" {
		t.Errorf("Expected value '%%los', got %v", values)
	}
}

func TestQueryBuilder_StringCaseInsensitive(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByStringColumn("name", "Carlos", StringExact, NonSensitive)).Apply()

	expected := "select * from accounts WHERE LOWER(name) = LOWER(%s);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 || values[0] != "Carlos" {
		t.Errorf("Expected value 'Carlos', got %v", values)
	}
}

func TestQueryBuilder_StringContainsCaseInsensitive(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByStringColumn("name", "CAR", StringContains, NonSensitive)).Apply()

	expected := "select * from accounts WHERE name ILIKE %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 || values[0] != "%CAR%" {
		t.Errorf("Expected value '%%CAR%%', got %v", values)
	}
}

func TestQueryBuilder_StringInClause(t *testing.T) {
	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	result := qb.Where(ByStringColumn("name", "carlos", "john")).Apply()

	expected := "select * from accounts WHERE name IN %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}

	stringSlice, ok := values[0].([]string)
	if !ok {
		t.Errorf("Expected []string, got %T", values[0])
	}

	if len(stringSlice) != 2 || stringSlice[0] != "carlos" || stringSlice[1] != "john" {
		t.Errorf("Expected [carlos, john], got %v", stringSlice)
	}
}
