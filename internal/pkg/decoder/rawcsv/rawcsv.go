package rawcsv

import (
	"fmt"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type Decoder struct {
	groupCostCenterBudgetParsers []*regexp.Regexp
	unitCostCenterBudgetParsers  []*regexp.Regexp
}

func NewDecoder() Decoder {
	return Decoder{
		groupCostCenterBudgetParsers: []*regexp.Regexp{
			regexp.MustCompile(rxBasis("(?P<id>\\d+)")),
			regexp.MustCompile(rxBasis("\\\"?(?P<id>[0-9][0-9]?) \\+? ")),
		},
		unitCostCenterBudgetParsers: []*regexp.Regexp{
			regexp.MustCompile(rxBasis("\\\"?\\d\\.\\d\\.\\d\\.\\d{2}\\.(?P<id>\\d+) ")),
		},
	}
}

func (p *Decoder) Decode(line string) (model.CostCenterType, []string, *regexp.Regexp) {
	for _, parser := range p.unitCostCenterBudgetParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return model.CostCenterUnit, matches, parser
	}

	for _, parser := range p.groupCostCenterBudgetParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return model.CostCenterGroup, matches, parser
	}

	panic(fmt.Sprintf("No parser found for line '%s'", line))
}

func rxBasis(rxID string) string {
	return fmt.Sprintf(
		"^%s%s,%s,%s,%s,%s,%s,%s",
		rxID,
		rxDesc,
		rxFloatNumber,
		rxNumber("_2021"),
		rxNumber("_2022"),
		rxNumber("_2023"),
		rxNumber("_2024"),
		rxNumber("_2025"),
	)
}

const (
	rxFloatNumber = "\"(?P<_2020>-?\\d+(\\.\\d+)*(,\\d+)?)\""
)

var (
	rxDesc = fmt.Sprintf("(?P<desc>[ %s\\-\\.\\\",\\)\\(\\d]*)", decoder.RxGermanLetter)
)

func rxNumber(name string) string {
	return fmt.Sprintf("(?P<%s>-?\\d+(\\.\\d+)*)", name)
}
