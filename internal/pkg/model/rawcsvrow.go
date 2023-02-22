package model

import "regexp"

type RowType = string

const (
	RowTypeUnitAccount RowType = "unit"
	RowTypeOther       RowType = "other"
)

type RawCSVRow struct {
	Tpe     RowType
	Matches []string
	Regexp  *regexp.Regexp
}
