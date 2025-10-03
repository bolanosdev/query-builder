package querybuilder

import (
	"fmt"
	"time"
)

func ByIntColumn(column string, values ...int) QueryCondition {
	if err := validateColumnName(column); err != nil {
		panic(err)
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

func ByStringColumn(column string, values ...any) QueryCondition {
	if err := validateColumnName(column); err != nil {
		panic(err)
	}

	if len(values) == 0 {
		return QueryCondition{}
	}

	var stringValues []string
	mt := StringExact
	sensitivity := Sensitive

	for _, val := range values {
		switch v := val.(type) {
		case string:
			stringValues = append(stringValues, v)
		case StringMatchType:
			mt = v
		case StringSensitivityType:
			sensitivity = v
		}
	}

	if len(stringValues) == 0 {
		return QueryCondition{}
	}

	if len(stringValues) > 1 {
		return QueryCondition{
			condition:   fmt.Sprintf("%s IN $1", column),
			value:       stringValues,
			placeholder: "%s",
		}
	}

	value := stringValues[0]
	caseSensitive := sensitivity == Sensitive

	var condition string
	var actualValue string
	var operator string

	if caseSensitive {
		operator = "="
		if mt != StringExact {
			operator = "LIKE"
		}
	} else {
		operator = "ILIKE"
	}

	switch mt {
	case StringExact:
		if caseSensitive {
			condition = fmt.Sprintf("%s = $1", column)
		} else {
			condition = fmt.Sprintf("LOWER(%s) = LOWER($1)", column)
		}
		actualValue = value
	case StringContains:
		condition = fmt.Sprintf("%s %s $1", column, operator)
		actualValue = fmt.Sprintf("%%%s%%", value)
	case StringStartsWith:
		condition = fmt.Sprintf("%s %s $1", column, operator)
		actualValue = fmt.Sprintf("%s%%", value)
	case StringEndsWith:
		condition = fmt.Sprintf("%s %s $1", column, operator)
		actualValue = fmt.Sprintf("%%%s", value)
	}

	return QueryCondition{
		condition:   condition,
		value:       actualValue,
		placeholder: "%s",
	}
}

func ByDateColumn(column string, values ...any) QueryCondition {
	if err := validateColumnName(column); err != nil {
		panic(err)
	}

	if len(values) == 0 {
		return QueryCondition{}
	}

	var dates []time.Time
	rangeType := DateExact

	for _, val := range values {
		switch v := val.(type) {
		case time.Time:
			dates = append(dates, v)
		case DateRangeType:
			rangeType = v
		}
	}

	if len(dates) == 0 {
		return QueryCondition{}
	}

	var condition string
	var value any
	var placeholder string

	switch rangeType {
	case DateExact:
		if len(dates) >= 1 {
			condition = fmt.Sprintf("DATE_TRUNC('day', %s) = DATE_TRUNC('day', $1::timestamp)", column)
			value = dates[0].Format(time.RFC3339)
			placeholder = "%s"
		}
	case DateAfter:
		if len(dates) >= 1 {
			condition = fmt.Sprintf("%s > $1", column)
			value = dates[0].Format(time.RFC3339)
			placeholder = "%s"
		}
	case DateBefore:
		if len(dates) >= 1 {
			condition = fmt.Sprintf("%s < $1", column)
			value = dates[0].Format(time.RFC3339)
			placeholder = "%s"
		}
	case DateBetween:
		if len(dates) >= 2 {
			condition = fmt.Sprintf("%s BETWEEN '%s' AND '%s'", column, dates[0].Format(time.RFC3339), dates[1].Format(time.RFC3339))
			value = nil
			placeholder = ""
		}
	}

	return QueryCondition{
		condition:   condition,
		value:       value,
		placeholder: placeholder,
	}
}
