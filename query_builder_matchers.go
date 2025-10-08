package querybuilder

import (
	"fmt"
	"time"
)

func ByIntColumn(column string, values []int) QueryCondition {
	if err := validateColumnName(column); err != nil {
		panic(err)
	}

	if len(values) == 0 {
		return QueryCondition{}
	}

	if len(values) == 1 {
		return QueryCondition{
			condition:   fmt.Sprintf("%s = $1", column),
			value:       values[0],
			placeholder: "%v",
		}
	}
	return QueryCondition{
		condition:   fmt.Sprintf("%s IN $1", column),
		value:       values,
		placeholder: "%v",
	}
}

func ByStringColumn(column string, values []string, options ...any) QueryCondition {
	if err := validateColumnName(column); err != nil {
		panic(err)
	}

	if len(values) == 0 {
		return QueryCondition{}
	}

	// Parse options - can be StringOpts struct, StringMatchType, or both StringMatchType and StringSensitivityType
	var mt StringMatchType = StringExact
	var sensitivity StringSensitivityType = Sensitive

	for _, opt := range options {
		switch v := opt.(type) {
		case StringOpts:
			if v.Match != 0 {
				mt = v.Match
			}
			if v.Sensitivity != 0 {
				sensitivity = v.Sensitivity
			}
		case StringMatchType:
			mt = v
		case StringSensitivityType:
			sensitivity = v
		}
	}

	if len(values) > 1 {
		return QueryCondition{
			condition:   fmt.Sprintf("%s IN $1", column),
			value:       values,
			placeholder: "%s",
		}
	}

	value := values[0]
	caseSensitive := sensitivity == Sensitive

	var condition string
	var actualValue string

	switch mt {
	case StringExact:
		if caseSensitive {
			condition = fmt.Sprintf("%s = $1", column)
		} else {
			condition = fmt.Sprintf("LOWER(%s) = LOWER($1)", column)
		}
		actualValue = value
	case StringContains:
		if caseSensitive {
			condition = fmt.Sprintf("%s LIKE '%%' || $1 || '%%'", column)
		} else {
			condition = fmt.Sprintf("LOWER(%s) LIKE '%%' || LOWER($1) || '%%'", column)
		}
		actualValue = value
	case StringStartsWith:
		if caseSensitive {
			condition = fmt.Sprintf("%s LIKE $1 || '%%'", column)
		} else {
			condition = fmt.Sprintf("LOWER(%s) LIKE LOWER($1) || '%%'", column)
		}
		actualValue = value
	case StringEndsWith:
		if caseSensitive {
			condition = fmt.Sprintf("%s LIKE '%%' || $1", column)
		} else {
			condition = fmt.Sprintf("LOWER(%s) LIKE '%%' || LOWER($1)", column)
		}
		actualValue = value
	}

	return QueryCondition{
		condition:   condition,
		value:       actualValue,
		placeholder: "%s",
	}
}

func ByDateColumn(column string, dates Dates) QueryCondition {
	if err := validateColumnName(column); err != nil {
		panic(err)
	}

	// Check if all dates are zero (empty)
	if dates.On.IsZero() && dates.After.IsZero() && dates.Before.IsZero() {
		return QueryCondition{}
	}

	var condition string
	var value any
	var placeholder string

	// Priority: On field takes precedence
	if !dates.On.IsZero() {
		// Exact date match using DATE function (SQLite compatible)
		condition = fmt.Sprintf("DATE(%s) = DATE($1)", column)
		value = dates.On.Format(time.RFC3339)
		placeholder = "%s"
		return QueryCondition{
			condition:   condition,
			value:       value,
			placeholder: placeholder,
		}
	}

	// Determine query type based on which dates are set
	hasAfter := !dates.After.IsZero()
	hasBefore := !dates.Before.IsZero()

	if hasAfter && hasBefore {
		// Both dates set: BETWEEN query with two placeholders (inclusive)
		condition = fmt.Sprintf("%s >= $1 AND %s <= $2", column, column)
		value = []string{dates.After.Format(time.RFC3339), dates.Before.Format(time.RFC3339)}
		placeholder = "%s"
	} else if hasAfter && !hasBefore {
		// Only after set: AFTER query (exclusive)
		condition = fmt.Sprintf("%s > $1", column)
		value = dates.After.Format(time.RFC3339)
		placeholder = "%s"
	} else if !hasAfter && hasBefore {
		// Only before set: BEFORE query (exclusive)
		condition = fmt.Sprintf("%s < $1", column)
		value = dates.Before.Format(time.RFC3339)
		placeholder = "%s"
	}

	return QueryCondition{
		condition:   condition,
		value:       value,
		placeholder: placeholder,
	}
}
