package querybuilder

import (
	"fmt"
	"strings"
)

func (qb *QueryBuilder) Where(conditions ...QueryCondition) *QueryBuilder {
	for _, cond := range conditions {
		if cond.isGroup {
			var groupParts []string
			for _, groupCond := range cond.groupConds {
				// Expand slice values for IN clauses; otherwise single placeholder
				if strings.Contains(groupCond.condition, "IN $1") {
					switch v := groupCond.value.(type) {
					case []int:
						placeholders := make([]string, len(v))
						for i := range v {
							placeholders[i] = fmt.Sprintf("$%d", qb.argCounter+i)
						}
						expanded := strings.Replace(groupCond.condition, "IN $1", "IN ("+strings.Join(placeholders, ", ")+")", 1)
						groupParts = append(groupParts, expanded)
						for _, item := range v {
							qb.values = append(qb.values, item)
						}
						qb.argCounter += len(v)
					case []string:
						placeholders := make([]string, len(v))
						for i := range v {
							placeholders[i] = fmt.Sprintf("$%d", qb.argCounter+i)
						}
						expanded := strings.Replace(groupCond.condition, "IN $1", "IN ("+strings.Join(placeholders, ", ")+")", 1)
						groupParts = append(groupParts, expanded)
						for _, item := range v {
							qb.values = append(qb.values, item)
						}
						qb.argCounter += len(v)
					default:
						tempPlaceholder := fmt.Sprintf("$%d", qb.argCounter)
						condition := strings.Replace(groupCond.condition, "$1", tempPlaceholder, 1)
						groupParts = append(groupParts, condition)
						qb.values = append(qb.values, groupCond.value)
						qb.argCounter++
					}
				} else {
					// Handle date range queries with []string values (two placeholders)
					if v, ok := groupCond.value.([]string); ok && (strings.Contains(groupCond.condition, "$1") && strings.Contains(groupCond.condition, "$2")) {
						placeholder1 := fmt.Sprintf("$%d", qb.argCounter)
						placeholder2 := fmt.Sprintf("$%d", qb.argCounter+1)
						condition := strings.Replace(groupCond.condition, "$1", placeholder1, 1)
						condition = strings.Replace(condition, "$2", placeholder2, 1)
						groupParts = append(groupParts, condition)
						for _, item := range v {
							qb.values = append(qb.values, item)
						}
						qb.argCounter += len(v)
					} else {
						tempPlaceholder := fmt.Sprintf("$%d", qb.argCounter)
						condition := strings.Replace(groupCond.condition, "$1", tempPlaceholder, 1)
						groupParts = append(groupParts, condition)
						qb.values = append(qb.values, groupCond.value)
						qb.argCounter++
					}
				}
			}
			groupCondition := "(" + strings.Join(groupParts, " "+cond.groupOp+" ") + ")"
			qb.conditions = append(qb.conditions, groupCondition)
			qb.operators = append(qb.operators, "AND")
		} else {
			if strings.Contains(cond.condition, "IN $1") {
				switch v := cond.value.(type) {
				case []int:
					placeholders := make([]string, len(v))
					for i := range v {
						placeholders[i] = fmt.Sprintf("$%d", qb.argCounter+i)
					}
					condition := strings.Replace(cond.condition, "IN $1", "IN ("+strings.Join(placeholders, ", ")+")", 1)
					qb.conditions = append(qb.conditions, condition)
					qb.operators = append(qb.operators, "AND")
					for _, item := range v {
						qb.values = append(qb.values, item)
					}
					qb.argCounter += len(v)
				case []string:
					placeholders := make([]string, len(v))
					for i := range v {
						placeholders[i] = fmt.Sprintf("$%d", qb.argCounter+i)
					}
					condition := strings.Replace(cond.condition, "IN $1", "IN ("+strings.Join(placeholders, ", ")+")", 1)
					qb.conditions = append(qb.conditions, condition)
					qb.operators = append(qb.operators, "AND")
					for _, item := range v {
						qb.values = append(qb.values, item)
					}
					qb.argCounter += len(v)
				default:
					tempPlaceholder := fmt.Sprintf("$%d", qb.argCounter)
					condition := strings.Replace(cond.condition, "$1", tempPlaceholder, 1)
					qb.conditions = append(qb.conditions, condition)
					qb.operators = append(qb.operators, "AND")
					qb.values = append(qb.values, cond.value)
					qb.argCounter++
				}
			} else {
				// Handle date range queries with []string values (two placeholders)
				if v, ok := cond.value.([]string); ok && (strings.Contains(cond.condition, "$1") && strings.Contains(cond.condition, "$2")) {
					placeholder1 := fmt.Sprintf("$%d", qb.argCounter)
					placeholder2 := fmt.Sprintf("$%d", qb.argCounter+1)
					condition := strings.Replace(cond.condition, "$1", placeholder1, 1)
					condition = strings.Replace(condition, "$2", placeholder2, 1)
					qb.conditions = append(qb.conditions, condition)
					qb.operators = append(qb.operators, "AND")
					for _, item := range v {
						qb.values = append(qb.values, item)
					}
					qb.argCounter += len(v)
				} else {
					tempPlaceholder := fmt.Sprintf("$%d", qb.argCounter)
					condition := strings.Replace(cond.condition, "$1", tempPlaceholder, 1)
					qb.conditions = append(qb.conditions, condition)
					qb.operators = append(qb.operators, "AND")
					qb.values = append(qb.values, cond.value)
					qb.argCounter++
				}
			}
		}
	}
	return qb
}

func Or(conditions ...QueryCondition) QueryCondition {
	return QueryCondition{
		isGroup:    true,
		groupConds: conditions,
		groupOp:    "OR",
	}
}

func And(conditions ...QueryCondition) QueryCondition {
	return QueryCondition{
		isGroup:    true,
		groupConds: conditions,
		groupOp:    "AND",
	}
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitValue = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offsetValue = offset
	return qb
}
