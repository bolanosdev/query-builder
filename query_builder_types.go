package querybuilder

import "time"

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

// StringOpts configures string matching behavior.
// Zero values default to StringExact and Sensitive.
type StringOpts struct {
    Match       StringMatchType
    Sensitivity StringSensitivityType
}

// Dates represents a date range with optional after, before, and exact date times.
// Priority: On > After/Before combination
// If On is set: query is exact date match (DATE_TRUNC)
// If After is zero and Before is set: query is "before Before"
// If Before is zero and After is set: query is "after After"
// If both After and Before are set: query is "between After and Before"
// If all are zero: returns empty condition
type Dates struct {
    On     time.Time // Exact date match (takes priority over After/Before)
    After  time.Time // Range start (after this date)
    Before time.Time // Range end (before this date)
}
