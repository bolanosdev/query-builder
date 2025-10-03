package querybuilder

type StringMatchType int

const (
	StringExact StringMatchType = iota
	StringContains
	StringStartsWith
	StringEndsWith
)

type StringSensitivityType int

const (
	Sensitive StringSensitivityType = iota
	NonSensitive
)

type DateRangeType int

const (
	DateExact DateRangeType = iota
	DateAfter
	DateBefore
	DateBetween
)

type SortDirection int

const (
	SortAsc SortDirection = iota
	SortDesc
)
