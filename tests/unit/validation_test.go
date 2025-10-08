package querybuilder_test

import (
	"testing"
	"time"

	. "github.com/bolanosdev/query-builder"
)

func TestValidation_InvalidColumnName_IntColumn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid column name")
		}
	}()

	ByIntColumn("id; DROP TABLE users--", []int{1})
}

func TestValidation_InvalidColumnName_StringColumn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid column name")
		}
	}()

	ByStringColumn("name' OR '1'='1", []string{"john"}, StringOpts{})
}

func TestValidation_InvalidColumnName_DateColumn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid column name")
		}
	}()

	ByDateColumn("created_at/*comment*/", Dates{After: time.Now()})
}

func TestValidation_InvalidColumnName_SortField(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid column name")
		}
	}()

	Sort("id WHERE 1=1")
}

func TestValidation_ValidColumnNames(t *testing.T) {
	validColumns := []string{
		"id",
		"user_name",
		"users.id",
		"table_name.column_name",
		"_internal_id",
		"column123",
	}

	for _, col := range validColumns {
		// Should not panic
		ByIntColumn(col, []int{1})
		ByStringColumn(col, []string{"value"}, StringOpts{})
		Sort(col)
	}
}

func TestValidation_InvalidColumnNames(t *testing.T) {
	invalidColumns := []string{
		"id; DROP TABLE",
		"id' OR '1'='1",
		"id WHERE 1=1",
		"id/*comment*/",
		"id--comment",
		"id OR 1=1",
		"",
		"123column",   // starts with number
		"column name", // contains space
	}

	for _, col := range invalidColumns {
		t.Run("IntColumn_"+col, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for invalid column name: %s", col)
				}
			}()
			ByIntColumn(col, []int{1})
		})

		t.Run("StringColumn_"+col, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for invalid column name: %s", col)
				}
			}()
			ByStringColumn(col, []string{"value"}, StringOpts{})
		})

		t.Run("SortField_"+col, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for invalid column name: %s", col)
				}
			}()
			Sort(col)
		})
	}
}
