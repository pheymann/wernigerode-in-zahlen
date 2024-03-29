package metadata

import (
	"fmt"
	"regexp"
	"strings"

	"wernigerode-in-zahlen.de/internal/pkg/decoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

type MetadataDecoder struct {
	departmentRegex                []*regexp.Regexp
	productClassRegex              []*regexp.Regexp
	productDomainRegex             []*regexp.Regexp
	productGroupRegex              []*regexp.Regexp
	productRegex                   []*regexp.Regexp
	subProductRegex                []*regexp.Regexp
	descriptionDetectionRegex      *regexp.Regexp
	missionAndTargetDetectionRegex *regexp.Regexp
	missionAndTargetRegex          []*regexp.Regexp
	objectivesDetectionRegex       *regexp.Regexp
	servicesDetectionRegex         *regexp.Regexp
}

func NewMetadataDecoder() MetadataDecoder {
	return MetadataDecoder{
		departmentRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Dezernat/( )+Fachbereich (?P<department>\d+),+(?P<department_name>[ %s]+),+verantwortlich:,*( )*(?P<accountable>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
		},
		productClassRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produktklasse (?P<product_class>\d+),+(?P<product_class_name>[ %s]+),+verantwortlich:,*( )*(?P<accountable>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
		},
		productDomainRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produktbereich \d+\.(?P<product_domain>\d+),+"(?P<product_domain_name>[ %s,]+)",+zuständig:,*( )*(?P<responsible>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produktbereich \d+\.(?P<product_domain>\d+),+(?P<product_domain_name>[ %s]+),+zuständig:,*( )*(?P<responsible>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
		},
		productGroupRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produktgruppe \d+\.\d+\.(?P<product_group>\d+),+"(?P<product_group_name>[ %s,]+)",+Produktart:,*( )*(?P<desc>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produktgruppe \d+\.\d+\.(?P<product_group>\d+),+(?P<product_group_name>[ %s]+),+Produktart:,*( )*(?P<desc>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
		},
		productRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produkt \d+\.\d+\.\d+\.(?P<product>\d+),+(?P<product_name>[ %s]+),+Rechtsbindung:,*( )*(?P<legal_requirement>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^Produkt \d+\.\d+\.\d+\.(?P<product>\d+),+"(?P<product_name>[ %s,]+)",+Rechtsbindung:,*( )*(?P<legal_requirement>[ %s]+)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
		},
		subProductRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^Unterprodukt \d+\.\d+\.\d+\.\d+\.(?P<sub_product>\d+),+"(?P<sub_product_name>[ %s,]+)",+`,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^Unterprodukt \d+\.\d+\.\d+\.\d+\.(?P<sub_product>\d+),+(?P<sub_product_name>[ %s]+),+`,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
		},
		descriptionDetectionRegex:      regexp.MustCompile("^Beschreibung,+"),
		missionAndTargetDetectionRegex: regexp.MustCompile("^Auftrag,+Zielgruppe,+"),
		missionAndTargetRegex: []*regexp.Regexp{
			regexp.MustCompile(
				fmt.Sprintf(
					`^"(?P<mission>[ %s,]*)",+"(?P<target>[ %s,]*)"`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^"(?P<mission>[ %s,]*)",+(?P<target>[ %s]*)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^(?P<mission>[ %s]*),+"(?P<target>[ %s,]*)"`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
				),
			),
			regexp.MustCompile(
				fmt.Sprintf(
					`^(?P<mission>[ %s]*),+(?P<target>[ %s]*)`,
					decoder.RxGermanPlusSpecialLetter,
					decoder.RxGermanPlusSpecialLetter,
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
	fmt.Printf("Sub-Product; %+v\n", d.subProductRegex)
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
	if !p.decodeProduct(metadata, lines[4]) {
		panic(fmt.Sprintf("Expected product but got '%s'.\nregex: %v", lines[4], p.productRegex))
	}

	var dropToLine = 5
	if p.decodeSubProduct(metadata, lines[dropToLine]) {
		dropToLine++
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
					strings.TrimSpace(strings.Trim(descriptionCleanupRegex.ReplaceAllString(line, "$1"), ",")),
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
	return forMatchingRegex(p.departmentRegex, line, func(regex *regexp.Regexp, matches []string) {
		metadata.Department = model.Department{
			ID:          decoder.DecodeString(regex, "department", matches),
			Name:        decoder.DecodeString(regex, "department_name", matches),
			Accountable: decoder.DecodeString(regex, "accountable", matches),
		}
	})
}

func (p MetadataDecoder) decodeProductClass(metadata *model.Metadata, line string) bool {
	return forMatchingRegex(p.productClassRegex, line, func(regex *regexp.Regexp, matches []string) {
		metadata.ProductClass = model.ProductClass{
			ID:          decoder.DecodeString(regex, "product_class", matches),
			Name:        decoder.DecodeString(regex, "product_class_name", matches),
			Accountable: decoder.DecodeString(regex, "accountable", matches),
		}
	})
}

func (p MetadataDecoder) decodeProductDomain(metadata *model.Metadata, line string) bool {
	return forMatchingRegex(p.productDomainRegex, line, func(regex *regexp.Regexp, matches []string) {
		metadata.ProductDomain = model.ProductDomain{
			ID:          decoder.DecodeString(regex, "product_domain", matches),
			Name:        decoder.DecodeString(regex, "product_domain_name", matches),
			Responsible: decoder.DecodeString(regex, "responsible", matches),
		}
	})
}

func (p MetadataDecoder) decodeProductGroup(metadata *model.Metadata, line string) bool {
	return forMatchingRegex(p.productGroupRegex, line, func(regex *regexp.Regexp, matches []string) {
		metadata.ProductGroup = model.ProductGroup{
			ID:   decoder.DecodeString(regex, "product_group", matches),
			Name: decoder.DecodeString(regex, "product_group_name", matches),
			Desc: decoder.DecodeString(regex, "desc", matches),
		}
	})
}

func (p MetadataDecoder) decodeProduct(metadata *model.Metadata, line string) bool {
	return forMatchingRegex(p.productRegex, line, func(regex *regexp.Regexp, matches []string) {
		metadata.Product = model.Product{
			ID:               decoder.DecodeString(regex, "product", matches),
			Name:             decoder.DecodeString(regex, "product_name", matches),
			LegalRequirement: decoder.DecodeString(regex, "legal_requirement", matches),
		}
	})
}

func (p MetadataDecoder) decodeSubProduct(metadata *model.Metadata, line string) bool {
	return forMatchingRegex(p.subProductRegex, line, func(regex *regexp.Regexp, matches []string) {
		metadata.SubProduct = &model.SubProduct{
			ID:   decoder.DecodeString(regex, "sub_product", matches),
			Name: decoder.DecodeString(regex, "sub_product_name", matches),
		}
	})
}

var (
	descriptionCleanupRegex = regexp.MustCompile(
		fmt.Sprintf("\"([ %s,-]+)\",+", decoder.RxGermanLetter),
	)
)

func (p MetadataDecoder) decodeMissionAndTargets(line string) (string, string) {
	var mission = ""
	var target = ""

	matched := forMatchingRegex(p.missionAndTargetRegex, line, func(regex *regexp.Regexp, matches []string) {
		mission = decoder.DecodeString(regex, "mission", matches)
		target = decoder.DecodeString(regex, "target", matches)
	})

	if !matched {
		panic(fmt.Sprintf("Expected mission and targets but got '%s'.", line))
	}
	return mission, target
}

var (
	objectivesCleanupRegex = regexp.MustCompile(
		fmt.Sprintf("\"([ %s,-]+)\",+", decoder.RxGermanLetter),
	)
)

func forMatchingRegex(regex []*regexp.Regexp, line string, callback func(regex *regexp.Regexp, matches []string)) bool {
	for _, regex := range regex {
		matches := regex.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		callback(regex, matches)

		return true
	}

	return false
}
