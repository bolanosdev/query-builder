package querybuilder

func Sort(field string, direction ...SortDirection) SortField {
	if err := validateColumnName(field); err != nil {
		panic(err)
	}

	dir := SortAsc
	if len(direction) > 0 {
		dir = direction[0]
	}
	return SortField{
		field:     field,
		direction: dir,
	}
}

func (qb *QueryBuilder) SortBy(fields ...SortField) *QueryBuilder {
	qb.sortFields = append(qb.sortFields, fields...)
	return qb
}
