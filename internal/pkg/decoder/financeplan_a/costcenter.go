package financeplan_a

import (
	"regexp"
	"strconv"
	"strings"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type FinancePlanACostCenterDecoder struct {
}

func NewFinancePlanACostCenterDecoder() FinancePlanACostCenterDecoder {
	return FinancePlanACostCenterDecoder{}
}

func (p FinancePlanACostCenterDecoder) Decode(tpe model.CostCenterType, matches []string, parser *regexp.Regexp) model.FinancePlanACostCenter {
	return model.FinancePlanACostCenter{
		Id:         decoder.DecodeString(parser, "id", matches),
		Tpe:        tpe,
		Desc:       strings.TrimSpace(decoder.DecodeString(parser, "desc", matches)),
		Budget2020: decodeBudget(parser, "_2020", matches),
		Budget2021: decodeBudget(parser, "_2021", matches),
		Budget2022: decodeBudget(parser, "_2022", matches),
		Budget2023: decodeBudget(parser, "_2023", matches),
		Budget2024: decodeBudget(parser, "_2024", matches),
		Budget2025: decodeBudget(parser, "_2025", matches),
	}
}

func decodeBudget(parser *regexp.Regexp, matchLabel string, matches []string) float64 {
	// 123.456,78 -> 123456.78
	strNumber := strings.ReplaceAll(
		strings.ReplaceAll(
			matches[parser.SubexpIndex(matchLabel)],
			".",
			"",
		),
		",",
		".",
	)

	i, err := strconv.ParseFloat(strNumber, 64)
	if err != nil {
		panic(err)
	}
	return i
}
