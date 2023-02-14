package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	parser := NewParser()

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parser.parse(line)
	}
}

type Parser struct {
	regexParsers []*regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		regexParsers: []*regexp.Regexp{
			unsafeCompileParser(rxBasis("\\\"?\\d\\.\\d\\.\\d\\.\\d{2}\\.(?P<id>\\d+) ")),
			unsafeCompileParser(rxBasis("(?P<id>\\d+)")),
			unsafeCompileParser(rxBasis("\\\"?(?P<id>[0-9][0-9]?) \\+? ")),
		},
	}
}

func unsafeCompileParser(regex string) *regexp.Regexp {
	fmt.Println(regex)

	parser, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}

	return parser
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
	rxDesc        = "(?P<desc>[ \\w\u00c4\u00e4\u00d6\u00f6\u00dc\u00fc\u00df\\-\\.\\\",\\)\\(\\d]*)" // umlauts and ß encoded in utf-8
	rxFloatNumber = "\"(?P<_2020>-?\\d+(\\.\\d+)*(,\\d+)?)\""
)

func rxNumber(name string) string {
	return fmt.Sprintf("(?P<%s>-?\\d+(\\.\\d+)*)", name)
}

func (p *Parser) parse(line string) Row {
	for _, parser := range p.regexParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		row := parseRow(parser, line, matches)

		if false {
			fmt.Printf("------------------------\n%s\n%+v\n\n", line, row)
		}
		return row
	}

	panic(fmt.Sprintf("No parser found for line '%s'", line))
}

type Row struct {
	id         string
	desc       string
	budget2020 float64
	budget2021 float64
	budget2022 float64
	budget2023 float64
	budget2024 float64
	budget2025 float64
}

func parseRow(parser *regexp.Regexp, line string, matches []string) Row {
	return Row{
		id:         parseString(parser, "id", matches),
		desc:       parseString(parser, "desc", matches),
		budget2020: parseBudget(parser, "_2020", matches),
		budget2021: parseBudget(parser, "_2021", matches),
		budget2022: parseBudget(parser, "_2022", matches),
		budget2023: parseBudget(parser, "_2023", matches),
		budget2024: parseBudget(parser, "_2024", matches),
		budget2025: parseBudget(parser, "_2025", matches),
	}
}

func parseBudget(parser *regexp.Regexp, matchLabel string, matches []string) float64 {
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

func parseString(parser *regexp.Regexp, matchLabel string, matches []string) string {
	return matches[parser.SubexpIndex(matchLabel)]
}
