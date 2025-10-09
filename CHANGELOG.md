# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2025-10-08

### Changed
- **BREAKING**: Refactored `Dates` struct field names from `Start`/`End` to `After`/`Before` for better clarity
- **BREAKING**: Replaced PostgreSQL-specific `ILIKE` with database-agnostic `LOWER()` function for case-insensitive string matching
- Changed date exact match from `DATE_TRUNC` to `DATE()` for SQLite compatibility
- Changed date between query from `BETWEEN` to `>= AND <=` for more explicit inclusive range queries
- Updated all case-insensitive string patterns to use `LOWER()` with `LIKE` instead of `ILIKE`

### Added
- Default limit of 10 when using `Offset()` without explicit `Limit()` to prevent unbounded result sets
- Expanded test database from 5 to 50 diverse users for comprehensive testing
- Added 9 comprehensive string integration tests covering all match types and case sensitivity

### Fixed
- Date range queries now work correctly with SQLite and other databases
- Case-insensitive string matching now works across all databases (SQLite, PostgreSQL, MySQL, etc.)
- Date query placeholder handling for range queries with two parameters

### Improved
- Refactored integration tests to use helper functions, reducing code duplication
- All date integration tests now use consistent patterns with helper functions
- Sort integration tests updated to verify results with larger dataset

## [0.1.0] - 2025-10-03

### Added
- Initial release of query-builder
- Type-safe column matchers for Int, String, and Date columns
- Logical grouping with AND/OR conditions
- String matching with Exact, Contains, StartsWith, EndsWith options
- Case-sensitive and case-insensitive string matching
- Date range comparisons (Exact, After, Before, Between)
- Sorting functionality with ASC/DESC support
- Pagination with Limit and Offset
- Column name validation for security
- Parameterized query support
- Comprehensive test suite with 89.1% code coverage
- 30 unit tests and 15 integration tests

[0.1.2]: https://github.com/bolanosdev/query-builder/releases/tag/v0.1.2
[0.1.0]: https://github.com/bolanosdev/query-builder/releases/tag/v0.1.0
