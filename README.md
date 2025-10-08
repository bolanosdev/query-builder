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
builder := qb.NewQueryBuilder("SELECT * FROM users")

// Add conditions and modifiers, then commit to get query and values
query, values := builder.
    Where(qb.ByIntColumn("age", []int{18, 25, 30})).
    SortBy(qb.Sort("name")).
    Limit(10).
    Commit()

// query: SELECT * FROM users WHERE age IN ($1, $2, $3) ORDER BY name LIMIT 10;
// values: [18, 25, 30]

// Use directly with database/sql
rows, err := db.Query(query, values...)
```
## Warnings
**⚠️ Security Note:** Always use `Commit()` to get both query and values for parameterized queries. The library uses PostgreSQL-style `$n` placeholders to prevent SQL injection.

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
qb.Where(qb.ByIntColumn("id", []int{1}))
// → WHERE id = $1

// Multiple AND conditions
qb.Where(
    qb.ByIntColumn("id", []int{1}),
    qb.ByStringColumn("name", []string{"john"}),
)
// → WHERE id = $1 AND name = $2

// OR conditions using helper
qb.Where(qb.Or(
    qb.ByIntColumn("id", []int{1}),
    qb.ByStringColumn("name", []string{"john"}),
))
// → WHERE (id = $1 OR name = $2)

// Nested conditions
qb.Where(
    qb.ByIntColumn("status", []int{1}),
    qb.Or(
        qb.ByStringColumn("role", []string{"admin"}),
        qb.ByStringColumn("role", []string{"moderator"}),
    ),
)
// → WHERE status = $1 AND (role = $2 OR role = $3)
```

### Column Matchers

#### Integer Columns

```go
// Single value
qb.ByIntColumn("id", []int{1})
// → id = $1

// Multiple values (IN clause)
qb.ByIntColumn("id", []int{1, 2, 3})
// → id IN ($1, $2, $3)
```

#### String Columns

```go
// Exact match (default)
qb.ByStringColumn("name", []string{"john"})
// → name = $1

// Multiple values (IN clause)
qb.ByStringColumn("name", []string{"john", "jane"})
// → name IN ($1, $2)

// Contains (case-sensitive by default) - Simplified syntax
qb.ByStringColumn("name", []string{"joh"}, qb.StringContains)
// → name LIKE $1 (value: "%joh%")

// Starts with
qb.ByStringColumn("name", []string{"joh"}, qb.StringStartsWith)
// → name LIKE $1 (value: "joh%")

// Ends with
qb.ByStringColumn("name", []string{"ohn"}, qb.StringEndsWith)
// → name LIKE $1 (value: "%ohn")

// Case-insensitive exact match - Pass both match type and sensitivity
qb.ByStringColumn("name", []string{"JOHN"}, qb.StringExact, qb.NonSensitive)
// → LOWER(name) = LOWER($1)

// Case-insensitive contains
qb.ByStringColumn("name", []string{"joh"}, qb.StringContains, qb.NonSensitive)
// → LOWER(name) LIKE '%' || LOWER($1) || '%'

// Case-insensitive starts with
qb.ByStringColumn("name", []string{"joh"}, qb.StringStartsWith, qb.NonSensitive)
// → LOWER(name) LIKE LOWER($1) || '%'

// Case-insensitive ends with
qb.ByStringColumn("name", []string{"ohn"}, qb.StringEndsWith, qb.NonSensitive)
// → LOWER(name) LIKE '%' || LOWER($1)

// Alternative: Using StringOpts struct (for complex configurations)
qb.ByStringColumn("name", []string{"joh"}, qb.StringOpts{
    Match: qb.StringContains,
    Sensitivity: qb.NonSensitive,
})
// → LOWER(name) LIKE '%' || LOWER($1) || '%'
```

**Flexible Options:**
You can pass options in three ways:
1. **No options** - Defaults to exact match, case-sensitive
2. **Direct parameters** - `StringMatchType` and/or `StringSensitivityType`
3. **StringOpts struct** - For explicit configuration

```go
type StringOpts struct {
    Match       StringMatchType       // Optional, defaults to StringExact
    Sensitivity StringSensitivityType // Optional, defaults to Sensitive
}
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

// Exact date match (using On field)
qb.ByDateColumn("created_at", qb.Dates{On: date})
// → DATE_TRUNC('day', created_at) = DATE_TRUNC('day', $1::timestamp)

// After date (using After field)
qb.ByDateColumn("created_at", qb.Dates{After: date})
// → created_at > $1

// Before date (using Before field)
qb.ByDateColumn("created_at", qb.Dates{Before: date})
// → created_at < $1

// Between dates (using both After and Before)
afterDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
beforeDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
qb.ByDateColumn("created_at", qb.Dates{After: afterDate, Before: beforeDate})
// → created_at BETWEEN $1 AND $2
```

**Dates Structure:**
```go
type Dates struct {
    On     time.Time // Exact date match (takes priority over After/Before)
    After  time.Time // Range start (after this date)
    Before time.Time // Range end (before this date)
}
```

**Query Logic:**
- If `On` is set → Exact date match using DATE_TRUNC
- If only `After` is set → After query (`> $1`)
- If only `Before` is set → Before query (`< $1`)
- If both `After` and `Before` are set → Between query
- If all are zero → Returns empty condition

### Sorting

```go
// Single field (ASC by default)
qb.SortBy(qb.Sort("name"))
// → ORDER BY name

// Single field DESC
qb.SortBy(qb.Sort("created_at", qb.SortDesc))
// → ORDER BY created_at DESC

// Multiple fields
qb.SortBy(
    qb.Sort("name"),
    qb.Sort("created_at", qb.SortDesc),
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

// Offset only (applies default LIMIT 10)
qb.Offset(20)
// → LIMIT 10 OFFSET 20

// Both
qb.Limit(50).Offset(20)
// → LIMIT 50 OFFSET 20
```

**Note:** When using `Offset()` without an explicit `Limit()`, a default limit of 10 is automatically applied to prevent unbounded result sets.

### Generating SQL

```go
// Commit returns both query and values in one call
query, values := qb.Commit()
// query: "SELECT * FROM users WHERE id = $1;"
// values: []any{1}
rows, err := db.Query(query, values...)
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
    builder := qb.NewQueryBuilder("SELECT * FROM orders")
    
    startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
    
    query, values := builder.
        Where(
            qb.ByIntColumn("status", []int{1, 2}),
            qb.Or(
                qb.ByStringColumn("customer", []string{"VIP"}, qb.StringContains),
                qb.ByIntColumn("amount", []int{1000}),
            ),
            qb.ByDateColumn("created_at", qb.Dates{After: startDate, Before: endDate}),
        ).
        SortBy(
            qb.Sort("created_at", qb.SortDesc),
            qb.Sort("amount", qb.SortDesc),
        ).
        Limit(50).
        Offset(0).
        Commit()
    
    fmt.Println("SQL:", query)
    fmt.Println("Values:", values)
    
    // Use with database/sql
    rows, err := db.Query(query, values...)
}
```

**Output:**
```sql
SELECT * FROM orders 
WHERE status IN ($1, $2) 
  AND (customer LIKE $3 OR amount = $4) 
  AND created_at BETWEEN $5 AND $6
ORDER BY created_at DESC, amount DESC 
LIMIT 50 
OFFSET 0;
```

**Values:** `[1, 2, "%VIP%", 1000, "2024-01-01T00:00:00Z", "2024-12-31T00:00:00Z"]`

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




## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
