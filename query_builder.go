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

func (qb *QueryBuilder) Commit() (string, []any) {
    query := qb.baseQuery

    if len(qb.conditions) > 0 {
        whereClause := " WHERE " + qb.conditions[0]
        for i := 1; i < len(qb.conditions); i++ {
            whereClause += " " + qb.operators[i] + " " + qb.conditions[i]
        }
        query += whereClause
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

    // Apply default limit of 10 if offset is set but limit is not
    limitToApply := qb.limitValue
    if qb.offsetValue >= 0 && qb.limitValue < 0 {
        limitToApply = 10
    }

    if limitToApply >= 0 {
        query += fmt.Sprintf(" LIMIT %d", limitToApply)
    }

    if qb.offsetValue >= 0 {
        query += fmt.Sprintf(" OFFSET %d", qb.offsetValue)
    }

    return query + ";", qb.values
}
