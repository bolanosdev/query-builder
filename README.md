# query-builder

A fluent query builder for Go that makes it easy to build SQL queries with type-safe matchers and clean syntax.

## Features

- 🔍 **Type-safe column matchers** - Int, String, and Date column matchers
- 🔗 **Logical grouping** - AND/OR conditions with proper nesting
- 📝 **String matching** - Exact, Contains, StartsWith, EndsWith with case sensitivity options
- 📅 **Date ranges** - Exact, After, Before, Between date comparisons
- 📊 **Sorting** - Single or multiple field sorting with ASC/DESC
- 📄 **Pagination** - Limit and Offset support
- ✅ **Well tested** - 89.1% code coverage with unit and integration tests

## Installation

```bash
go get github.com/bolanosdev/query-builder
```

## Quick Start

```go
import qb "github.com/bolanosdev/query-builder"

// Create a query builder
query := qb.NewQueryBuilder("SELECT * FROM users")

// Add conditions and modifiers
result := query.
    Where(qb.ByIntColumn("age", 18, 25, 30)).
    Sort(qb.SortBy("name")).
    Limit(10).
    Apply()

// Output: SELECT * FROM users WHERE age IN %v ORDER BY name LIMIT 10;
values := query.GetValues() // [[]int{18, 25, 30}]
```
## Warnings
**⚠️ Security Note:** `GetFormattedValues()` is for display/debugging only. Always use `GetValues()` with parameterized queries for database execution to prevent SQL injection.

**⚠️ Early Development Warning**: This library is in active development and has not been battle-tested in production environments. While it includes security features like column name validation and parameterized queries, please thoroughly test and review the generated SQL before using in production. Use at your own risk.


## API Reference

### Creating a Query Builder

```go
qb := qb.NewQueryBuilder("SELECT * FROM accounts")
```

### Adding Conditions

#### Where Clause

```go
// Single condition
qb.Where(qb.ByIntColumn("id", 1))
// → WHERE id = %v

// Multiple AND conditions
qb.Where(
    qb.ByIntColumn("id", 1),
    qb.ByStringColumn("name", "john"),
)
// → WHERE id = %v AND name = %s

// OR conditions using helper
qb.Where(qb.Or(
    qb.ByIntColumn("id", 1),
    qb.ByStringColumn("name", "john"),
))
// → WHERE (id = %v OR name = %s)

// Nested conditions
qb.Where(
    qb.ByIntColumn("status", 1),
    qb.Or(
        qb.ByStringColumn("role", "admin"),
        qb.ByStringColumn("role", "moderator"),
    ),
)
// → WHERE status = %v AND (role = %s OR role = %s)
```

### Column Matchers

#### Integer Columns

```go
// Single value
qb.ByIntColumn("id", 1)
// → id = %v

// Multiple values (IN clause)
qb.ByIntColumn("id", 1, 2, 3)
// → id IN %v
```

#### String Columns

```go
// Exact match (default)
qb.ByStringColumn("name", "john")
// → name = %s

// Multiple values (IN clause)
qb.ByStringColumn("name", "john", "jane")
// → name IN %s

// Contains (case-sensitive by default)
qb.ByStringColumn("name", "joh", qb.StringContains)
// → name LIKE %s (value: "%joh%")

// Starts with
qb.ByStringColumn("name", "joh", qb.StringStartsWith)
// → name LIKE %s (value: "joh%")

// Ends with
qb.ByStringColumn("name", "ohn", qb.StringEndsWith)
// → name LIKE %s (value: "%ohn")

// Case-insensitive
qb.ByStringColumn("name", "JOHN", qb.StringExact, qb.NonSensitive)
// → LOWER(name) = LOWER(%s)

qb.ByStringColumn("name", "joh", qb.StringContains, qb.NonSensitive)
// → name ILIKE %s (value: "%joh%")
```

**String Match Types:**
- `StringExact` - Exact match (default)
- `StringContains` - Contains substring
- `StringStartsWith` - Starts with substring
- `StringEndsWith` - Ends with substring

**String Sensitivity:**
- `Sensitive` - Case-sensitive (default)
- `NonSensitive` - Case-insensitive

#### Date Columns

```go
import "time"

date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

// Exact date (default)
qb.ByDateColumn("created_at", date)
qb.ByDateColumn("created_at", date, qb.DateExact)
// → DATE_TRUNC('day', created_at) = DATE_TRUNC('day', %s::timestamp)

// After date
qb.ByDateColumn("created_at", date, qb.DateAfter)
// → created_at > %s

// Before date
qb.ByDateColumn("created_at", date, qb.DateBefore)
// → created_at < %s

// Between dates
startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
qb.ByDateColumn("created_at", startDate, endDate, qb.DateBetween)
// → created_at BETWEEN '2024-01-01T00:00:00Z' AND '2024-12-31T00:00:00Z'
```

**Date Range Types:**
- `DateExact` - Exact date match (default)
- `DateAfter` - After date
- `DateBefore` - Before date
- `DateBetween` - Between two dates

### Sorting

```go
// Single field (ASC by default)
qb.SortBy(qb.SortBy("name"))
// → ORDER BY name

// Single field DESC
qb.SortBy(qb.SortBy("created_at", qb.SortDesc))
// → ORDER BY created_at DESC

// Multiple fields
qb.SortBy(
    qb.SortBy("name"),
    qb.SortBy("created_at", qb.SortDesc),
)
// → ORDER BY name, created_at DESC
```

**Sort Directions:**
- `SortAsc` - Ascending (default)
- `SortDesc` - Descending

### Pagination

```go
// Limit only
qb.Limit(10)
// → LIMIT 10

// Offset only
qb.Offset(20)
// → OFFSET 20

// Both
qb.Limit(10).Offset(20)
// → LIMIT 10 OFFSET 20
```

### Generating SQL

```go
// Get the SQL query with placeholders
sql := qb.Apply()
// Returns: "SELECT * FROM users WHERE id = %v;"

// Get parameter values (for database execution)
values := qb.GetValues()
// Returns: []any{1}

// ✅ CORRECT - Use with database/sql
db.Query(qb.Apply(), qb.GetValues()...)

// Get formatted values (for display/logging ONLY)
formatted := qb.GetFormattedValues()
// Returns: []any{"'john'"} for strings, etc.

// ✅ CORRECT - For debugging/logging
fmt.Println(fmt.Sprintf(qb.Apply(), qb.GetFormattedValues()...))
// Prints: "SELECT * FROM users WHERE name = 'john';"

// ❌ NEVER DO THIS - SQL injection risk!
query := fmt.Sprintf(qb.Apply(), qb.GetFormattedValues()...)
db.Query(query)  // Don't execute formatted queries!
```

## Complete Example

```go
package main

import (
    "fmt"
    "time"
    qb "github.com/bolanosdev/query-builder"
)

func main() {
    query := qb.NewQueryBuilder("SELECT * FROM orders")
    
    startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
    
    result := query.
        Where(
            qb.ByIntColumn("status", 1, 2), // status IN (1, 2)
            qb.Or(
                qb.ByStringColumn("customer", "VIP", qb.StringContains),
                qb.ByIntColumn("amount", 1000), // amount > threshold
            ),
            qb.ByDateColumn("created_at", startDate, endDate, qb.DateBetween),
        ).
        Sort(
            qb.SortBy("created_at", qb.SortDesc),
            qb.SortBy("amount", qb.SortDesc),
        ).
        Limit(50).
        Offset(0).
        Apply()
    
    values := query.GetValues()
    
    fmt.Println("SQL:", result)
    fmt.Println("Values:", values)
    
    // Use with database/sql
    rows, err := db.Query(result, values...)
}
```

**Output:**
```sql
SELECT * FROM orders 
WHERE status IN %v 
  AND (customer LIKE %s OR amount = %v) 
  AND created_at BETWEEN '2024-01-01T00:00:00Z' AND '2024-12-31T23:59:59Z' 
ORDER BY created_at DESC, amount DESC 
LIMIT 50 
OFFSET 0;
```

## Logical Grouping

Use `Or()` and `And()` helper functions to create grouped conditions:

```go
// (A OR B) AND C
qb.Where(
    qb.Or(conditionA, conditionB),
    conditionC,
)
// → WHERE (A OR B) AND C

// A AND (B OR C)
qb.Where(
    conditionA,
    qb.Or(conditionB, conditionC),
)
// → WHERE A AND (B OR C)

// (A AND B) OR (C AND D)
qb.Where(qb.Or(
    qb.And(conditionA, conditionB),
    qb.And(conditionC, conditionD),
))
// → WHERE ((A AND B) OR (C AND D))
```

## Testing

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration
```

**Test Coverage:** 89.1%
- 30 unit tests
- 15 integration tests

## Architecture

The library is organized into the following files:

- `query_builder.go` - Core types and Apply() logic
- `query_builder_conditions.go` - Where() and logical grouping (Or, And)
- `query_builder_matchers.go` - Column matchers (ByIntColumn, ByStringColumn, ByDateColumn)
- `query_builder_types.go` - Enums and constants
- `query_builder_sort.go` - Sorting functionality

## Why This Library?

### Problems Solved

1. **Type Safety** - Compile-time checking for match types and sort directions
2. **Clean Syntax** - Fluent API that reads like natural language
3. **No String Concatenation** - Avoid SQL injection vulnerabilities
4. **Reusable Conditions** - Build conditions independently and combine them
5. **Complex Queries Made Simple** - Handle nested OR/AND logic easily

### Design Decisions

- **Simplified API**: Removed `And()` and `Or()` methods from QueryBuilder - only use `Where()` for adding conditions
- **Explicit Prefixes**: All enums use prefixes (`String*`, `Date*`, `Sort*`) to avoid naming conflicts
- **Immutable Values**: Query builder maintains internal state, values are copied on retrieval
- **PostgreSQL Focus**: Optimized for PostgreSQL syntax (ILIKE, DATE_TRUNC, etc.)

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
