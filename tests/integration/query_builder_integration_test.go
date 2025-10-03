package querybuilder_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/bolanosdev/query-builder"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE accounts (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO accounts (id, name, created_at) VALUES 
		(1, 'carlos', '2024-01-01 10:00:00'),
		(2, 'john', '2024-01-02 11:00:00'),
		(3, 'jane', '2024-01-03 12:00:00'),
		(4, 'alice', '2024-01-04 13:00:00'),
		(5, 'bob', '2024-01-05 14:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	return db
}

func TestQueryBuilder_Integration_SingleCondition(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByIntColumn("id", 1))
	result := qb.Apply()
	values := qb.GetValues()

	execQuery := fmt.Sprintf(result, values...)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		if id != 1 || name != "carlos" {
			t.Errorf("Expected id=1, name=carlos, got id=%d, name=%s", id, name)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

func TestQueryBuilder_Integration_MultipleConditions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	conditions := []QueryCondition{
		ByIntColumn("id", 1),
		ByStringColumn("name", "carlos"),
	}

	qb := NewQueryBuilder(query)
	qb.Where(conditions...)
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		if id != 1 || name != "carlos" {
			t.Errorf("Expected id=1, name=carlos, got id=%d, name=%s", id, name)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

func TestQueryBuilder_Integration_InClause(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByIntColumn("id", 1, 2))
	result := qb.Apply()
	values := qb.GetValues()

	var execQuery string
	if len(values) > 0 {
		if intSlice, ok := values[0].([]int); ok {
			parts := make([]string, len(intSlice))
			for i, v := range intSlice {
				parts[i] = fmt.Sprintf("%d", v)
			}
			inClause := fmt.Sprintf("(%s)", strings.Join(parts, ", "))
			execQuery = fmt.Sprintf(result, inClause)
		}
	}

	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}
}

func TestQueryBuilder_Integration_NestedConditions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)

	qb.Where(ByIntColumn("id", 1), Or(ByIntColumn("id", 2), ByStringColumn("name", "carlos")))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	ids := make(map[int]bool)
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		ids[id] = true
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}

	if !ids[1] {
		t.Errorf("Expected id 1, got %v", ids)
	}
}

func TestQueryBuilder_Integration_StringContains(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByStringColumn("name", "ar", StringContains))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	execQuery = strings.ReplaceAll(execQuery, "ILIKE", "LIKE")
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	names := make(map[string]bool)
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names[name] = true
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}

	if !names["carlos"] {
		t.Errorf("Expected carlos, got %v", names)
	}
}

func TestQueryBuilder_Integration_StringStartsWith(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByStringColumn("name", "jo", StringStartsWith))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	execQuery = strings.ReplaceAll(execQuery, "ILIKE", "LIKE")
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	names := make(map[string]bool)
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names[name] = true
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}

	if !names["john"] {
		t.Errorf("Expected john, got %v", names)
	}
}

func TestQueryBuilder_Integration_CaseInsensitive(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByStringColumn("name", "CARLOS", StringExact, NonSensitive))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		if name != "carlos" {
			t.Errorf("Expected carlos, got %s", name)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

func TestQueryBuilder_Integration_StringInClause(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByStringColumn("name", "carlos", "john"))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	names := make(map[string]bool)
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names[name] = true
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}

	if !names["carlos"] || !names["john"] {
		t.Errorf("Expected carlos and john, got %v", names)
	}
}

func TestQueryBuilder_Integration_Where(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByIntColumn("id", 1))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		if id != 1 || name != "carlos" {
			t.Errorf("Expected id=1, name=carlos, got id=%d, name=%s", id, name)
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

func TestQueryBuilder_Integration_WhereWithOr(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(Or(ByIntColumn("id", 1), ByStringColumn("name", "john")))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	names := make(map[string]bool)
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names[name] = true
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}

	if !names["carlos"] || !names["john"] {
		t.Errorf("Expected carlos and john, got %v", names)
	}
}

func TestQueryBuilder_Integration_Sort(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.SortBy(Sort("name", SortDesc))
	result := qb.Apply()

	execQuery := result
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names = append(names, name)
	}

	expected := []string{"john", "jane", "carlos", "bob", "alice"}
	if len(names) != len(expected) {
		t.Errorf("Expected %d rows, got %d", len(expected), len(names))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected name[%d] = %s, got %s", i, expected[i], name)
		}
	}
}

func TestQueryBuilder_Integration_SortMultipleFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.SortBy(Sort("name"), Sort("id", SortDesc))
	result := qb.Apply()

	execQuery := result
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names = append(names, name)
	}

	if len(names) != 5 {
		t.Errorf("Expected 5 rows, got %d", len(names))
	}
}

func TestQueryBuilder_Integration_LimitOffset(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Limit(2).Offset(1)
	result := qb.Apply()

	execQuery := result
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	ids := make(map[int]bool)
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		ids[id] = true
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}

	if !ids[2] || !ids[3] {
		t.Errorf("Expected ids 2 and 3, got %v", ids)
	}
}

func TestQueryBuilder_Integration_CompleteQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	qb.Where(ByIntColumn("id", 1, 2, 3, 4, 5)).
		SortBy(Sort("name", SortDesc)).
		Limit(3).
		Offset(1)
	result := qb.Apply()
	values := qb.GetValues()

	var execQuery string
	if len(values) > 0 {
		if intSlice, ok := values[0].([]int); ok {
			parts := make([]string, len(intSlice))
			for i, v := range intSlice {
				parts[i] = fmt.Sprintf("%d", v)
			}
			inClause := fmt.Sprintf("(%s)", strings.Join(parts, ", "))
			execQuery = fmt.Sprintf(result, inClause)
		}
	}

	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names = append(names, name)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(names))
	}

	expected := []string{"jane", "carlos", "bob"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected name[%d] = %s, got %s", i, expected[i], name)
		}
	}
}

func TestQueryBuilder_Integration_DateAfter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	date := time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC)
	qb.Where(ByDateColumn("created_at", date, DateAfter))
	result := qb.Apply()
	values := qb.GetFormattedValues()

	execQuery := fmt.Sprintf(result, values...)
	fmt.Printf("Generated query: %s\n", execQuery)

	rows, err := db.Query(execQuery)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, execQuery)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var createdAt string
		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 rows (after 2024-01-02), got %d", count)
	}
}
