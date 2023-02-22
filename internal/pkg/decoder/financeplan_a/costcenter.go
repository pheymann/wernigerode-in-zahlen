package financeplan_a

import (
	"regexp"
	"strings"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type Decoder struct {
}

func New() Decoder {
	return Decoder{}
}

func (p Decoder) Decode(tpe model.CostCenterType, matches []string, parser *regexp.Regexp) model.FinancePlanACostCenter {
	return model.FinancePlanACostCenter{
		Id:         decoder.DecodeString(parser, "id", matches),
		Tpe:        tpe,
		Desc:       strings.TrimSpace(decoder.DecodeString(parser, "desc", matches)),
		Budget2020: decoder.DecodeBudget(parser, "_2020", matches),
		Budget2021: decoder.DecodeBudget(parser, "_2021", matches),
		Budget2022: decoder.DecodeBudget(parser, "_2022", matches),
		Budget2023: decoder.DecodeBudget(parser, "_2023", matches),
		Budget2024: decoder.DecodeBudget(parser, "_2024", matches),
		Budget2025: decoder.DecodeBudget(parser, "_2025", matches),
	}
}
