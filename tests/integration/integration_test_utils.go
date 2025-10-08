package querybuilder_test

import (
	"database/sql"
	"fmt"
	"testing"

	. "github.com/bolanosdev/query-builder"
)

type User struct {
	ID        int
	Name      string
	CreatedAt string
}

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

	// Insert 50 diverse users for comprehensive testing
	_, err = db.Exec(`
		INSERT INTO accounts (id, name, created_at) VALUES 
		(1, 'carlos', '2024-01-01 10:00:00'),
		(2, 'john', '2024-01-02 11:00:00'),
		(3, 'jane', '2024-01-03 12:00:00'),
		(4, 'alice', '2024-01-04 13:00:00'),
		(5, 'bob', '2024-01-05 14:00:00'),
		(6, 'charlie', '2024-01-06 15:00:00'),
		(7, 'david', '2024-01-07 16:00:00'),
		(8, 'emma', '2024-01-08 17:00:00'),
		(9, 'frank', '2024-01-09 18:00:00'),
		(10, 'grace', '2024-01-10 19:00:00'),
		(11, 'henry', '2024-01-11 20:00:00'),
		(12, 'isabel', '2024-01-12 21:00:00'),
		(13, 'jack', '2024-01-13 22:00:00'),
		(14, 'kate', '2024-01-14 23:00:00'),
		(15, 'liam', '2024-01-15 10:00:00'),
		(16, 'maria', '2024-01-16 11:00:00'),
		(17, 'nathan', '2024-01-17 12:00:00'),
		(18, 'olivia', '2024-01-18 13:00:00'),
		(19, 'peter', '2024-01-19 14:00:00'),
		(20, 'quinn', '2024-01-20 15:00:00'),
		(21, 'rachel', '2024-01-21 16:00:00'),
		(22, 'samuel', '2024-01-22 17:00:00'),
		(23, 'tina', '2024-01-23 18:00:00'),
		(24, 'uma', '2024-01-24 19:00:00'),
		(25, 'victor', '2024-01-25 20:00:00'),
		(26, 'wendy', '2024-01-26 21:00:00'),
		(27, 'xavier', '2024-01-27 22:00:00'),
		(28, 'yara', '2024-01-28 23:00:00'),
		(29, 'zack', '2024-01-29 10:00:00'),
		(30, 'anna', '2024-01-30 11:00:00'),
		(31, 'brian', '2024-01-31 12:00:00'),
		(32, 'clara', '2024-02-01 13:00:00'),
		(33, 'daniel', '2024-02-02 14:00:00'),
		(34, 'elena', '2024-02-03 15:00:00'),
		(35, 'felix', '2024-02-04 16:00:00'),
		(36, 'gina', '2024-02-05 17:00:00'),
		(37, 'hugo', '2024-02-06 18:00:00'),
		(38, 'iris', '2024-02-07 19:00:00'),
		(39, 'james', '2024-02-08 20:00:00'),
		(40, 'karen', '2024-02-09 21:00:00'),
		(41, 'lucas', '2024-02-10 22:00:00'),
		(42, 'maya', '2024-02-11 23:00:00'),
		(43, 'noah', '2024-02-12 10:00:00'),
		(44, 'oscar', '2024-02-13 11:00:00'),
		(45, 'paula', '2024-02-14 12:00:00'),
		(46, 'quentin', '2024-02-15 13:00:00'),
		(47, 'rosa', '2024-02-16 14:00:00'),
		(48, 'steve', '2024-02-17 15:00:00'),
		(49, 'tracy', '2024-02-18 16:00:00'),
		(50, 'ursula', '2024-02-19 17:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	return db
}

func mapUserIDs(users []User) []int {
	ids := make([]int, len(users))
	for i, user := range users {
		ids[i] = user.ID
	}
	return ids
}

func fetchOneUser(rows *sql.Rows) (User, error) {
	var u User
	if !rows.Next() {
		return u, fmt.Errorf("no rows")
	}
	if err := rows.Scan(&u.ID, &u.Name, &u.CreatedAt); err != nil {
		return u, err
	}
	return u, nil
}

func fetchAllUsers(rows *sql.Rows) ([]User, error) {
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func executeSortQuery(t *testing.T, db *sql.DB, sorts ...SortField) []User {
	t.Helper()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.SortBy(sorts...).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	users, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	return users
}

func executeDateQuery(t *testing.T, db *sql.DB, column string, dates Dates) []User {
	t.Helper()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.Where(ByDateColumn(column, dates)).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	users, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	return users
}

func executeWhereQuery(t *testing.T, db *sql.DB, conditions ...QueryCondition) []User {
	t.Helper()

	query := "select * from accounts"
	qb := NewQueryBuilder(query)
	finalQuery, values := qb.Where(conditions...).Commit()

	rows, err := db.Query(finalQuery, values...)
	if err != nil {
		t.Fatalf("Query failed: %v\nQuery: %s", err, finalQuery)
	}
	defer rows.Close()

	users, err := fetchAllUsers(rows)
	if err != nil {
		t.Fatalf("Failed to fetch users: %v", err)
	}

	return users
}
