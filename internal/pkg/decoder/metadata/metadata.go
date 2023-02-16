package metadata

import (
	"fmt"
	"regexp"
	"strings"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type MetadataDecoder struct {
	departmentRegex       *regexp.Regexp
	productClassRegex     *regexp.Regexp
	productDomainRegex    *regexp.Regexp
	productGroupRegex     *regexp.Regexp
	productRegex          *regexp.Regexp
	descriptionRegex      *regexp.Regexp
	missionAndTargetRegex *regexp.Regexp
}

func NewMetadataDecoder() MetadataDecoder {
	return MetadataDecoder{
		departmentRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Dezernat/( )+Fachbereich (?P<department>\\d+),(?P<department_name>[ %s-]+),+verantwortlich: (?P<accountable>[ %s-]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productClassRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Produktklasse (?P<product_class>\\d+),+(?P<product_class_name>[ %s-]+),+verantwortlich: (?P<accountable>[ %s-]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productDomainRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Produktbereich (?P<product_domain>\\d+\\.\\d+),+(?P<product_domain_name>[ %s-]+),+zust\u00e4ndig: +(?P<responsible>[ %s-]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productGroupRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Produktgruppe (?P<product_group>\\d+\\.\\d+\\.\\d+),+(?P<product_group_name>[ %s-]+),+Produktart: +(?P<desc>[ %s-]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Produkt (?P<product>\\d+\\.\\d+\\.\\d+\\.\\d+),+(?P<product_name>[ %s-]+),+Rechtsbindung: +(?P<legal_requirement>[ %s-]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		descriptionRegex:      regexp.MustCompile("^Beschreibung,+"),
		missionAndTargetRegex: regexp.MustCompile("^Auftrag,+Zielgruppe,+"),
	}
}

func (p MetadataDecoder) Decode(lines []string) model.Metadata {
	metadata := &model.Metadata{}

	if !p.decodeDepartment(metadata, lines[0]) {
		panic(fmt.Sprintf("Expected department but got '%s'", lines[0]))
	}
	if !p.decodeProductClass(metadata, lines[1]) {
		panic(fmt.Sprintf("Expected product class but got '%s'", lines[1]))
	}
	if !p.decodeProductDomain(metadata, lines[2]) {
		panic(fmt.Sprintf("Expected product domain but got '%s'.\nregex: %v", lines[2], p.productDomainRegex))
	}
	if !p.decodeProductGroup(metadata, lines[3]) {
		panic(fmt.Sprintf("Expected product group but got '%s'.\nregex: %v", lines[3], p.productGroupRegex))
	}

	var dropToLine = 4
	if p.decodeProduct(metadata, lines[4]) {
		dropToLine = 5
	}

	var state = ""
	var content = []string{}
	for _, line := range lines[dropToLine:] {
		if p.descriptionRegex.MatchString(line) {
			state = "description"
			continue
		}
		if p.missionAndTargetRegex.MatchString(line) {
			metadata.Description = strings.Join(content, "")
			content = []string{}

			state = "missionAndTarget"
			continue
		}

		if state == "description" {
			content = append(content, descriptionCleanupRegex.ReplaceAllString(line, "$1"))
			continue
		}

		fmt.Printf(">>> %+v\n", *metadata)
		fmt.Printf("=== %+v\n", content)
		panic(fmt.Sprintf("No parser found for line '%s'", line))
	}

	return *metadata
}

func (p MetadataDecoder) decodeDepartment(metadata *model.Metadata, line string) bool {
	matches := p.departmentRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.Department = model.Department{
		ID:          decoder.DecodeString(p.departmentRegex, "department", matches),
		Name:        decoder.DecodeString(p.departmentRegex, "department_name", matches),
		Accountable: decoder.DecodeString(p.departmentRegex, "accountable", matches),
	}

	return true
}

func (p MetadataDecoder) decodeProductClass(metadata *model.Metadata, line string) bool {
	matches := p.productClassRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.ProductClass = model.ProductClass{
		ID:          decoder.DecodeString(p.productClassRegex, "product_class", matches),
		Name:        decoder.DecodeString(p.productClassRegex, "product_class_name", matches),
		Accountable: decoder.DecodeString(p.productClassRegex, "accountable", matches),
	}

	return true
}

func (p MetadataDecoder) decodeProductDomain(metadata *model.Metadata, line string) bool {
	matches := p.productDomainRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.ProductDomain = model.ProductDomain{
		ID:          decoder.DecodeString(p.productDomainRegex, "product_domain", matches),
		Name:        decoder.DecodeString(p.productDomainRegex, "product_domain_name", matches),
		Responsible: decoder.DecodeString(p.productDomainRegex, "responsible", matches),
	}

	return true
}

func (p MetadataDecoder) decodeProductGroup(metadata *model.Metadata, line string) bool {
	matches := p.productGroupRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.ProductGroup = model.ProductGroup{
		ID:   decoder.DecodeString(p.productGroupRegex, "product_group", matches),
		Name: decoder.DecodeString(p.productGroupRegex, "product_group_name", matches),
		Desc: decoder.DecodeString(p.productGroupRegex, "desc", matches),
	}

	return true
}

func (p MetadataDecoder) decodeProduct(metadata *model.Metadata, line string) bool {
	matches := p.productRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.Product = model.Product{
		ID:               decoder.DecodeString(p.productRegex, "product", matches),
		Name:             decoder.DecodeString(p.productRegex, "product_name", matches),
		LegalRequirement: decoder.DecodeString(p.productRegex, "legal_requirement", matches),
	}

	return true
}

var (
	descriptionCleanupRegex = regexp.MustCompile(
		fmt.Sprintf("\"([ %s,-]+)\",+", decoder.RxGermanLetter),
	)
)
