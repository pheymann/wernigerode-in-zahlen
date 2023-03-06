package model

import "regexp"

type RowType = string

const (
	RowTypeOneOff       RowType = "oneoff"
	RowTypeUnitAccount  RowType = "unit"
	RowTypeSeparateLine RowType = "separate"
	RowTypeOther        RowType = "other"
	RowTypeIgnore       RowType = "ignore"
)

type RawCSVRow struct {
	Tpe     RowType
	Matches []string
	Regexp  *regexp.Regexp
}
