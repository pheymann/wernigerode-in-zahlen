package metadata

import (
	"fmt"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type MetadataDecoder struct {
	regexParser       *regexp.Regexp
	departmentRegex   *regexp.Regexp
	productClassRegex *regexp.Regexp
}

func NewMetadataDecoder() MetadataDecoder {
	return MetadataDecoder{
		regexParser: regexp.MustCompile(rxFileClassification),
		departmentRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Dezernat/( )+Fachbereich (?P<department>\\d+),(?P<department_name>[ %s]+),+verantwortlich: (?P<accountable>[ %s]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
		productClassRegex: regexp.MustCompile(
			fmt.Sprintf(
				"^Produktklasse (?P<product_class>\\d+),(?P<product_class_name>[ %s]+),+verantwortlich: (?P<accountable>[ %s]+)",
				decoder.RxGermanLetter,
				decoder.RxGermanLetter,
			),
		),
	}
}

const (
	rxFileClassification = "^assets/data/raw/(?P<department>\\d+)/(?P<product_class>\\d+)/(?P<product_domain>\\d+)/(?P<product_group>\\d+)/(?P<product>\\d+)/(?P<file_name>\\w+)\\.(?P<file_type>\\w+)"
)

func (p MetadataDecoder) Decode(filename string) model.Metadata {
	matches := p.regexParser.FindStringSubmatch(filename)

	return model.Metadata{
		Department:    decoder.DecodeString(p.regexParser, "department", matches),
		ProductClass:  decoder.DecodeString(p.regexParser, "product_class", matches),
		ProductDomain: decoder.DecodeString(p.regexParser, "product_domain", matches),
		ProductGroup:  decoder.DecodeString(p.regexParser, "product_group", matches),
		Product:       decoder.DecodeString(p.regexParser, "product", matches),
		FileName:      decoder.DecodeString(p.regexParser, "file_name", matches),
		FileType:      decoder.DecodeString(p.regexParser, "file_type", matches),
	}
}

func (p MetadataDecoder) DecodeV2(lines []string) model.Metadata {
	metadata := &model.Metadata{}

	for _, line := range lines {
		if p.decodeDepartment(metadata, line) {
			continue
		}
		if p.decodeProductClass(metadata, line) {
			continue
		}

		fmt.Printf(">>> %+v\n", *metadata)

		panic(fmt.Sprintf("No parser found for line '%s'", line))
	}

	return *metadata
}

func (p MetadataDecoder) decodeDepartment(metadata *model.Metadata, line string) bool {
	matches := p.departmentRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.Department = decoder.DecodeString(p.departmentRegex, "department", matches)

	return true
}

func (p MetadataDecoder) decodeProductClass(metadata *model.Metadata, line string) bool {
	matches := p.productClassRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false
	}

	metadata.ProductClass = decoder.DecodeString(p.productClassRegex, "product_class", matches)

	return true
}
