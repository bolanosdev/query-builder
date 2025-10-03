package querybuilder

import (
	"fmt"
	"regexp"
)

var columnNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`)

func validateColumnName(column string) error {
	if column == "" || !columnNameRegex.MatchString(column) {
		return fmt.Errorf("invalid column name: %s (must contain only letters, numbers, underscores, and dots)", column)
	}
	return nil
}
