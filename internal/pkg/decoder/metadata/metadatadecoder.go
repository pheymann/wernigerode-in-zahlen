package metadata

import (
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type MetadataDecoder struct {
	regexParser *regexp.Regexp
}

func NewMetadataDecoder() MetadataDecoder {
	return MetadataDecoder{
		regexParser: regexp.MustCompile(rxFileClassification),
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
