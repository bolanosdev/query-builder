package querybuilder

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	baseQuery   string
	conditions  []string
	operators   []string
	args        []any
	values      []any
	argCounter  int
	limitValue  int
	offsetValue int
	sortFields  []SortField
}

type QueryCondition struct {
	condition   string
	value       any
	placeholder string
	isGroup     bool
	groupConds  []QueryCondition
	groupOp     string
}

type SortField struct {
	field     string
	direction SortDirection
}

func NewQueryBuilder(query string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:   strings.TrimSpace(query),
		conditions:  []string{},
		operators:   []string{},
		args:        []any{},
		values:      []any{},
		argCounter:  1,
		limitValue:  -1,
		offsetValue: -1,
		sortFields:  []SortField{},
	}
}

func (qb *QueryBuilder) Apply() string {
	query := qb.baseQuery

	if len(qb.conditions) > 0 {
		whereClause := " WHERE " + qb.conditions[0]
		for i := 1; i < len(qb.conditions); i++ {
			whereClause += " " + qb.operators[i] + " " + qb.conditions[i]
		}
		query += whereClause

		for i, placeholder := range qb.args {
			tempPlaceholder := fmt.Sprintf("$%d", i+1)
			query = strings.Replace(query, tempPlaceholder, placeholder.(string), 1)
		}
	}

	if len(qb.sortFields) > 0 {
		var sortParts []string
		for _, field := range qb.sortFields {
			sortStr := field.field
			if field.direction == SortDesc {
				sortStr += " DESC"
			}
			sortParts = append(sortParts, sortStr)
		}
		query += " ORDER BY " + strings.Join(sortParts, ", ")
	}

	if qb.limitValue >= 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limitValue)
	}

	if qb.offsetValue >= 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offsetValue)
	}

	return query + ";"
}

func (qb *QueryBuilder) GetValues() []any {
	return qb.values
}

func (qb *QueryBuilder) GetFormattedValues() []any {
	formatted := make([]any, len(qb.values))
	for i, val := range qb.values {
		switch v := val.(type) {
		case string:
			formatted[i] = fmt.Sprintf("'%s'", v)
		case []string:
			parts := make([]string, len(v))
			for j, str := range v {
				parts[j] = fmt.Sprintf("'%s'", str)
			}
			formatted[i] = fmt.Sprintf("(%s)", strings.Join(parts, ", "))
		default:
			formatted[i] = val
		}
	}
	return formatted
}
