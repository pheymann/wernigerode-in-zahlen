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
	separateLineParser           *regexp.Regexp
}

func NewDecoder() Decoder {
	return Decoder{
		groupCostCenterBudgetParsers: []*regexp.Regexp{
			regexp.MustCompile(rxBasis(`(?P<id>\d+)`)),
			regexp.MustCompile(rxBasis(`(?P<id>[0-9][0-9]?) \+? `)),
			regexp.MustCompile(
				fmt.Sprintf(
					`^"(?P<id>[0-9][0-9]?) \+? (?P<desc>[ %s\-\.,\)\(\d&]*)",%s,%s,%s,%s,%s,%s`,
					decoder.RxGermanLetter,
					rxFloatNumber,
					rxNumber("_2021"),
					rxNumber("_2022"),
					rxNumber("_2023"),
					rxNumber("_2024"),
					rxNumber("_2025"),
				),
			),
		},
		unitCostCenterBudgetParsers: []*regexp.Regexp{
			regexp.MustCompile(rxBasis(`\d\.\d\.\d\.\d{2}\.(?P<id>\d+) `)),
			regexp.MustCompile(
				fmt.Sprintf(
					`^"\d\.\d\.\d\.\d{2}\.(?P<id>\d+) (?P<desc>[ %s\-\.,\)\(\d&%%]*)",+%s,%s,%s,%s,%s,%s`,
					decoder.RxGermanLetter,
					rxFloatNumber,
					rxNumber("_2021"),
					rxNumber("_2022"),
					rxNumber("_2023"),
					rxNumber("_2024"),
					rxNumber("_2025"),
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^\d\.\d\.\d\.\d{2}(?P<id>/\d+\.\d+) (?P<desc>[ %s\-\.\)\(\d&%%]*),+%s,%s,%s,%s,%s,%s`,
					decoder.RxGermanLetter,
					rxFloatNumber,
					rxNumber("_2021"),
					rxNumber("_2022"),
					rxNumber("_2023"),
					rxNumber("_2024"),
					rxNumber("_2025"),
				),
			),
		},
		separateLineParser: regexp.MustCompile(
			fmt.Sprintf(
				`^"?(?P<desc>[ %s\.&\(\)/>]+)"?,+`,
				decoder.RxGermanLetter,
			),
		),
	}
}

func (d *Decoder) Debug() {
	fmt.Println("=== DEBUG rawcsv ===")
	fmt.Printf("%+v\n", d.groupCostCenterBudgetParsers)
	fmt.Printf("%+v\n", d.unitCostCenterBudgetParsers)
	fmt.Printf("%+v\n", d.separateLineParser)
}

type DecodeType = string

const (
	DecodeTypeAccount      DecodeType = "account"
	DecodeTypeUnit         DecodeType = "unit"
	DeocdeTypeSeparateLine DecodeType = "separate"
)

func (p *Decoder) Decode(line string) (DecodeType, []string, *regexp.Regexp) {
	for _, parser := range p.unitCostCenterBudgetParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return DecodeTypeUnit, matches, parser
	}

	for _, parser := range p.groupCostCenterBudgetParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return DecodeTypeAccount, matches, parser
	}

	matches := p.separateLineParser.FindStringSubmatch(line)
	if len(matches) > 0 {
		return DeocdeTypeSeparateLine, matches, p.separateLineParser
	}

	panic(fmt.Sprintf("No parser found for line '%s'", line))
}

func (p *Decoder) Decode2(line string) model.RawCSVRow {
	for _, regex := range p.unitCostCenterBudgetParsers {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return model.RawCSVRow{
			Tpe:     model.RowTypeUnitAccount,
			Matches: matches,
			Regexp:  regex,
		}
	}

	for _, regex := range p.groupCostCenterBudgetParsers {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return model.RawCSVRow{
			Tpe:     model.RowTypeOther,
			Matches: matches,
			Regexp:  regex,
		}
	}

	matches := p.separateLineParser.FindStringSubmatch(line)
	if len(matches) > 0 {
		return model.RawCSVRow{
			Tpe:     model.RowTypeSeparateLine,
			Matches: matches,
			Regexp:  p.separateLineParser,
		}
	}

	panic(fmt.Sprintf("No parser found for line '%s'", line))
}

func rxBasis(rxID string) string {
	return fmt.Sprintf(
		"^%s%s,+%s,%s,%s,%s,%s,%s",
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
	rxDesc = fmt.Sprintf(`(?P<desc>[ %s\-\.\)\(\d&%%]*)`, decoder.RxGermanLetter)
)

func rxNumber(name string) string {
	return fmt.Sprintf("(?P<%s>-?\\d+(\\.\\d+)*)", name)
}
