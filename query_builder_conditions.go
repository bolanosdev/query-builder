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
				tempPlaceholder := fmt.Sprintf("$%d", qb.argCounter)
				condition := strings.Replace(groupCond.condition, "$1", tempPlaceholder, 1)
				groupParts = append(groupParts, condition)
				qb.args = append(qb.args, groupCond.placeholder)
				qb.values = append(qb.values, groupCond.value)
				qb.argCounter++
			}
			groupCondition := "(" + strings.Join(groupParts, " "+cond.groupOp+" ") + ")"
			qb.conditions = append(qb.conditions, groupCondition)
			qb.operators = append(qb.operators, "AND")
		} else {
			tempPlaceholder := fmt.Sprintf("$%d", qb.argCounter)
			condition := strings.Replace(cond.condition, "$1", tempPlaceholder, 1)
			qb.conditions = append(qb.conditions, condition)
			qb.operators = append(qb.operators, "AND")
			qb.args = append(qb.args, cond.placeholder)
			qb.values = append(qb.values, cond.value)
			qb.argCounter++
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
