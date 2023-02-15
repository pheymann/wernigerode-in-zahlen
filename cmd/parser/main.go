package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	debug = false
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	metadataParser := NewFinancePlanMetadataParser()
	rawCSVParser := NewRawCSVParser()
	financePlanParser := NewCostCenterFinancePlanParser()

	metadata := metadataParser.parse(filename)

	scanner := bufio.NewScanner(file)
	costCenterFinancePlans := []CostCenterFinancePlan{}

	for scanner.Scan() {
		line := scanner.Text()

		tpe, matches, regex := rawCSVParser.parse(line)
		financePlan := financePlanParser.parse(tpe, matches, regex)

		if debug {
			fmt.Printf("------------------------\n%s\n%+v\n\n", line, financePlan)
		}

		costCenterFinancePlans = append(costCenterFinancePlans, financePlan)
	}

	costCenterGroups, perGroupCostCenterUnits := separateCostCenterUnits(costCenterFinancePlans)

	writeCostCenterGroupFinancePlansAsCSV(costCenterGroups, metadata)
	writeCostCenterUnitFinancePlansAsCSV(perGroupCostCenterUnits, metadata)
}

type FinancePlanMetadata struct {
	department    string
	productClass  string
	productDomain string
	productGroup  string
	product       string
	fileName      string
	fileType      string
}

type FinancePlanMetadataParser struct {
	regexParser *regexp.Regexp
}

const (
	rxFileClassification = "^assets/data/raw/(?P<department>\\d+)/(?P<product_class>\\d+)/(?P<product_domain>\\d+)/(?P<product_group>\\d+)/(?P<product>\\d+)/(?P<file_name>\\w+)\\.(?P<file_type>\\w+)"
)

func NewFinancePlanMetadataParser() FinancePlanMetadataParser {
	return FinancePlanMetadataParser{
		regexParser: compileParser(rxFileClassification),
	}
}

func (p FinancePlanMetadataParser) parse(filename string) FinancePlanMetadata {
	matches := p.regexParser.FindStringSubmatch(filename)

	return FinancePlanMetadata{
		department:    parseString(p.regexParser, "department", matches),
		productClass:  parseString(p.regexParser, "product_class", matches),
		productDomain: parseString(p.regexParser, "product_domain", matches),
		productGroup:  parseString(p.regexParser, "product_group", matches),
		product:       parseString(p.regexParser, "product", matches),
		fileName:      parseString(p.regexParser, "file_name", matches),
		fileType:      parseString(p.regexParser, "file_type", matches),
	}
}

type RawCSVParser struct {
	groupCostCenterBudgetParsers []*regexp.Regexp
	unitCostCenterBudgetParsers  []*regexp.Regexp
}

func NewRawCSVParser() RawCSVParser {
	return RawCSVParser{
		groupCostCenterBudgetParsers: []*regexp.Regexp{
			compileParser(rxBasis("(?P<id>\\d+)")),
			compileParser(rxBasis("\\\"?(?P<id>[0-9][0-9]?) \\+? ")),
		},
		unitCostCenterBudgetParsers: []*regexp.Regexp{
			compileParser(rxBasis("\\\"?\\d\\.\\d\\.\\d\\.\\d{2}\\.(?P<id>\\d+) ")),
		},
	}
}

func compileParser(regex string) *regexp.Regexp {
	if debug {
		fmt.Println(regex)
	}

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

func (p *RawCSVParser) parse(line string) (CostCenterType, []string, *regexp.Regexp) {
	for _, parser := range p.unitCostCenterBudgetParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return CostCenterUnit, matches, parser
	}

	for _, parser := range p.groupCostCenterBudgetParsers {
		matches := parser.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return CostCenterGroup, matches, parser
	}

	panic(fmt.Sprintf("No parser found for line '%s'", line))
}

type CostCenterType = string

const (
	CostCenterGroup CostCenterType = "group"
	CostCenterUnit  CostCenterType = "unit"
)

type CostCenterFinancePlan struct {
	id         string
	tpe        CostCenterType
	desc       string
	budget2020 float64
	budget2021 float64
	budget2022 float64
	budget2023 float64
	budget2024 float64
	budget2025 float64
}

type CostCenterFinancePlanParser struct {
}

func NewCostCenterFinancePlanParser() CostCenterFinancePlanParser {
	return CostCenterFinancePlanParser{}
}

func (p CostCenterFinancePlanParser) parse(tpe CostCenterType, matches []string, parser *regexp.Regexp) CostCenterFinancePlan {
	return CostCenterFinancePlan{
		id:         parseString(parser, "id", matches),
		tpe:        tpe,
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

// separate CostCenter groups from units
func separateCostCenterUnits(financePlans []CostCenterFinancePlan) ([]CostCenterFinancePlan, map[string][]CostCenterFinancePlan) {
	var groups []CostCenterFinancePlan
	perGroupUnits := make(map[string][]CostCenterFinancePlan)

	currentCostCenterGroupID := ""
	currentGroupUnits := []CostCenterFinancePlan{}

	for _, financePlan := range financePlans {
		if financePlan.tpe == CostCenterGroup {
			perGroupUnits[currentCostCenterGroupID] = currentGroupUnits
			groups = append(groups, financePlan)

			currentCostCenterGroupID = financePlan.id
			currentGroupUnits = []CostCenterFinancePlan{}
		} else {
			currentGroupUnits = append(currentGroupUnits, financePlan)
		}
	}

	return groups, perGroupUnits
}

func (financePlan CostCenterFinancePlan) toCSV() string {
	return fmt.Sprintf(
		"%s,%s,%f,%f,%f,%f,%f,%f",
		financePlan.id,
		financePlan.desc,
		financePlan.budget2020,
		financePlan.budget2021,
		financePlan.budget2022,
		financePlan.budget2023,
		financePlan.budget2024,
		financePlan.budget2025,
	)
}

const (
	CSVHeader = "id,desc,_2020,_2021,_2022,_2023,_2024,_2025\n"
)

func writeCostCenterGroupFinancePlansAsCSV(financePlans []CostCenterFinancePlan, metadata FinancePlanMetadata) {
	content := CSVHeader
	filepath := fmt.Sprintf("assets/data/processed/%s/%s/%s/%s/", metadata.productClass, metadata.productDomain, metadata.productGroup, metadata.product)
	filename := "data.csv"

	for _, financePlan := range financePlans {
		content += financePlan.toCSV() + "\n"
	}

	writeFile(filepath, filename, content)
}

func writeCostCenterUnitFinancePlansAsCSV(financePlans map[string][]CostCenterFinancePlan, metadata FinancePlanMetadata) {
	for costCenterGroup, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		content := CSVHeader
		filepath := fmt.Sprintf("assets/data/processed/%s/%s/%s/%s/%s/", metadata.productClass, metadata.productDomain, metadata.productGroup, metadata.product, costCenterGroup)
		filename := "data.csv"

		for _, financePlan := range financePlans {
			content += financePlan.toCSV() + "\n"
		}

		writeFile(filepath, filename, content)
	}
}

func writeFile(filepath string, filename string, content string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.MkdirAll(filepath, 0700)
	}

	file, err := os.Create(filepath + filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	file.Sync()
}
