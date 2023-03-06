package rawcsv

import (
	"fmt"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type Decoder struct {
	oneOffBudgetParsers          []*regexp.Regexp
	groupCostCenterBudgetParsers []*regexp.Regexp
	unitCostCenterBudgetParsers  []*regexp.Regexp
	separateLineParsers          []*regexp.Regexp
	ignoreLineParsers            []*regexp.Regexp
}

func NewDecoder() Decoder {
	return Decoder{
		oneOffBudgetParsers: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^(?P<id>\d+)[ ]*(?P<desc>[ %s\-\.\)\(\d&€><]*),,,,,,+`,
					decoder.RxGermanLetter,
				),
			),
		},
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
		separateLineParsers: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^"(?P<desc>[ %s\.&\(\)/>,]+)",+`,
					decoder.RxGermanLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^(?P<desc>[ %s\.&\(\)/>€]+),+`,
					decoder.RxGermanLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^,"?(?P<desc>[ %s\.&\(\)/>]+)"?,+`,
					decoder.RxGermanLetter,
				),
			),
		},
		ignoreLineParsers: []*regexp.Regexp{
			regexp.MustCompile(`^"",,,,,+`),
		},
	}
}

func (d *Decoder) Debug() {
	fmt.Println("=== DEBUG rawcsv ===")
	fmt.Printf("%+v\n", d.oneOffBudgetParsers)
	fmt.Printf("%+v\n", d.groupCostCenterBudgetParsers)
	fmt.Printf("%+v\n", d.unitCostCenterBudgetParsers)
	fmt.Printf("%+v\n", d.separateLineParsers)
}

type DecodeType = string

const (
	DecodeTypeOneOffBudget DecodeType = "one-off"
	DecodeTypeAccount      DecodeType = "account"
	DecodeTypeUnit         DecodeType = "unit"
	DeocdeTypeSeparateLine DecodeType = "separate"
)

func (p *Decoder) Decode(line string) model.RawCSVRow {
	for _, regex := range p.oneOffBudgetParsers {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return model.RawCSVRow{
			Tpe:     model.RowTypeOneOff,
			Matches: matches,
			Regexp:  regex,
		}
	}

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

	for _, regex := range p.separateLineParsers {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return model.RawCSVRow{
			Tpe:     model.RowTypeSeparateLine,
			Matches: matches,
			Regexp:  regex,
		}
	}

	for _, regex := range p.ignoreLineParsers {
		if regex.MatchString(line) {
			return model.RawCSVRow{
				Tpe: model.RowTypeIgnore,
			}
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
