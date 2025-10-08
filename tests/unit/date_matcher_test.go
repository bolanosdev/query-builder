package querybuilder_test

import (
	"testing"
	"time"

	. "github.com/bolanosdev/query-builder"
)

func TestQueryBuilder_DateAfter(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, values := qb.Where(ByDateColumn("created_at", Dates{After: date})).Commit()

	expected := "select * from events WHERE created_at > $1;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}

func TestQueryBuilder_DateBefore(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	result, values := qb.Where(ByDateColumn("created_at", Dates{Before: date})).Commit()

	expected := "select * from events WHERE created_at < $1;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}

func TestQueryBuilder_DateBetween(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	after := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	before := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	result, values := qb.Where(ByDateColumn("created_at", Dates{After: after, Before: before})).Commit()

	expected := "select * from events WHERE created_at >= $1 AND created_at <= $2;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
}

func TestQueryBuilder_DateExact(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	result, values := qb.Where(ByDateColumn("created_at", Dates{On: date})).Commit()

	expected := "select * from events WHERE DATE(created_at) = DATE($1);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}
