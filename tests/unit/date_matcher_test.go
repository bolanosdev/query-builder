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
	result := qb.Where(ByDateColumn("created_at", date, DateAfter)).Apply()

	expected := "select * from events WHERE created_at > %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}

func TestQueryBuilder_DateBefore(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	result := qb.Where(ByDateColumn("created_at", date, DateBefore)).Apply()

	expected := "select * from events WHERE created_at < %s;"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_DateBetween(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	result := qb.Where(ByDateColumn("created_at", start, end, DateBetween)).Apply()

	expected := "select * from events WHERE created_at BETWEEN '2024-01-01T00:00:00Z' AND '2024-12-31T23:59:59Z';"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}

func TestQueryBuilder_DateExact(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	result := qb.Where(ByDateColumn("created_at", date, DateExact)).Apply()

	expected := "select * from events WHERE DATE_TRUNC('day', created_at) = DATE_TRUNC('day', %s::timestamp);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}

	values := qb.GetValues()
	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}
}

func TestQueryBuilder_DateDefault(t *testing.T) {
	query := "select * from events"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	result := qb.Where(ByDateColumn("created_at", date)).Apply()

	expected := "select * from events WHERE DATE_TRUNC('day', created_at) = DATE_TRUNC('day', %s::timestamp);"
	if result != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, result)
	}
}
