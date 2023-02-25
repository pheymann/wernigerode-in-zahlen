package metadata

import (
	"fmt"
	"regexp"
	"strings"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type MetadataDecoder struct {
	departmentRegex                *regexp.Regexp
	productClassRegex              *regexp.Regexp
	productDomainRegex             *regexp.Regexp
	productGroupRegex              *regexp.Regexp
	productRegex                   []*regexp.Regexp
	descriptionDetectionRegex      *regexp.Regexp
	missionAndTargetDetectionRegex *regexp.Regexp
	missionAndTargetRegex          []*regexp.Regexp
	objectivesDetectionRegex       *regexp.Regexp
	servicesDetectionRegex         *regexp.Regexp
}

func NewMetadataDecoder() MetadataDecoder {
	return MetadataDecoder{
		departmentRegex: regexp.MustCompile(
			fmt.Sprintf(
				`^Dezernat/( )+Fachbereich (?P<department>\d+),+(?P<department_name>[ %s\-]+),+verantwortlich:,*( )*(?P<accountable>[ %s\-]+)`,
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productClassRegex: regexp.MustCompile(
			fmt.Sprintf(
				`^Produktklasse (?P<product_class>\d+),+(?P<product_class_name>[ %s-]+),+verantwortlich:,*( )*(?P<accountable>[ %s-]+)`,
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productDomainRegex: regexp.MustCompile(
			fmt.Sprintf(
				`^Produktbereich (?P<product_domain>\d+\.\d+),+(?P<product_domain_name>[ %s-]+),+zuständig:,*( )*(?P<responsible>[ %s-]+)`,
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productGroupRegex: regexp.MustCompile(
			fmt.Sprintf(
				`^Produktgruppe (?P<product_group>\d+\.\d+\.\d+),+(?P<product_group_name>[ %s\-/]+),+Produktart:,*( )*(?P<desc>[ %s-]+)`,
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produkt (?P<product>\d+\.\d+\.\d+\.\d+),+(?P<product_name>[ %s\-]+),+Rechtsbindung:,*( )*(?P<legal_requirement>[ %s-]+)`,
					decoder.RxGermanLetter,
					decoder.RxGermanLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produkt (?P<product>\d+\.\d+\.\d+\.\d+),+"(?P<product_name>[ %s\-,]+)",+Rechtsbindung:,*( )*(?P<legal_requirement>[ %s-]+)`,
					decoder.RxGermanLetter,
					decoder.RxGermanLetter,
				),
			),
		},
		descriptionDetectionRegex:      regexp.MustCompile("^Beschreibung,+"),
		missionAndTargetDetectionRegex: regexp.MustCompile("^Auftrag,+Zielgruppe,+"),
		missionAndTargetRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^"(?P<mission>[ %s,;:\-]*)",+"(?P<target>[ %s,:\-]*)"`,
					decoder.RxGermanLetter,
					decoder.RxGermanLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^"(?P<mission>[ %s,;:\-]*)",+(?P<target>[ %s:\-]*)`,
					decoder.RxGermanLetter,
					decoder.RxGermanLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^(?P<mission>[ %s;:\-]*),+"(?P<target>[ %s,:\-]*)"`,
					decoder.RxGermanLetter,
					decoder.RxGermanLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^(?P<mission>[ %s;:\-]*),+(?P<target>[ %s:\-]*)`,
					decoder.RxGermanLetter,
					decoder.RxGermanLetter,
				),
			),
		},
		objectivesDetectionRegex: regexp.MustCompile("^Ziele,+"),
		servicesDetectionRegex:   regexp.MustCompile("^Leistung,+"),
	}
}

func (d MetadataDecoder) Debug() {
	fmt.Println("=== Debug Metadata ===")
	fmt.Printf("Department: %+v\n", d.departmentRegex)
	fmt.Printf("Product Class: %+v\n", d.productClassRegex)
	fmt.Printf("Product Domain: %+v\n", d.productDomainRegex)
	fmt.Printf("Product Group: %+v\n", d.productGroupRegex)
	fmt.Printf("Product; %+v\n", d.productRegex)
	fmt.Printf("Mission/Target: %+v\n", d.missionAndTargetRegex)
}

func (p MetadataDecoder) DecodeFromCSV(lines []string) model.Metadata {
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
	for _, line := range lines[dropToLine:] {
		if p.descriptionDetectionRegex.MatchString(line) {
			state = "description"
			continue
		}
		if p.missionAndTargetDetectionRegex.MatchString(line) {
			state = "missionAndTarget"
			continue
		}
		if p.objectivesDetectionRegex.MatchString(line) {
			state = "objectives"
			continue
		}
		if p.servicesDetectionRegex.MatchString(line) {
			state = "services"
			continue
		}

		if state == "description" {
			metadata.Description = strings.Join(
				[]string{
					metadata.Description,
					descriptionCleanupRegex.ReplaceAllString(line, "$1"),
				},
				" ",
			)
			continue
		}
		if state == "missionAndTarget" {
			mission, target := p.decodeMissionAndTargets(line)

			metadata.Mission = strings.Join(
				[]string{
					metadata.Mission,
					strings.TrimSpace(strings.Trim(mission, ",")),
				},
				" ",
			)
			metadata.Target = strings.Join(
				[]string{
					metadata.Target,
					strings.TrimSpace(strings.Trim(target, ",")),
				},
				"",
			)
			continue
		}
		if state == "objectives" {
			metadata.Objectives = strings.Join(
				[]string{
					metadata.Objectives,
					objectivesCleanupRegex.ReplaceAllString(line, "$1"),
				},
				" ",
			)
			continue
		}
		if state == "services" {
			metadata.Services = strings.Join(
				[]string{
					metadata.Services,
					strings.ReplaceAll(strings.ReplaceAll(line, "\"", ""), ",,,", ""),
				},
				" ",
			)
			continue
		}

		panic(fmt.Sprintf("No parser found for line '%s'\nmetadata: %+v", line, *metadata))
	}

	metadata.Validate()

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
	for _, regex := range p.productRegex {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		metadata.Product = model.Product{
			ID:               decoder.DecodeString(regex, "product", matches),
			Name:             decoder.DecodeString(regex, "product_name", matches),
			LegalRequirement: decoder.DecodeString(regex, "legal_requirement", matches),
		}

		return true
	}

	return false
}

var (
	descriptionCleanupRegex = regexp.MustCompile(
		fmt.Sprintf("\"([ %s,-]+)\",+", decoder.RxGermanLetter),
	)
)

func (p MetadataDecoder) decodeMissionAndTargets(line string) (string, string) {
	for _, regex := range p.missionAndTargetRegex {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		return decoder.DecodeString(regex, "mission", matches), decoder.DecodeString(regex, "target", matches)
	}

	panic(fmt.Sprintf("Expected mission and targets but got '%s'.", line))
}

var (
	objectivesCleanupRegex = regexp.MustCompile(
		fmt.Sprintf("\"([ %s,-]+)\",+", decoder.RxGermanLetter),
	)
)
