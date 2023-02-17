package metadata

import (
	"fmt"

	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Write(metadataJSON string, metadata model.Metadata) {
	filepath := fmt.Sprintf(
		"assets/data/processed/%s/%s/%s/%s/",
		metadata.ProductClass,
		metadata.ProductDomain,
		metadata.ProductGroup,
		metadata.Product,
	)
	filename := fmt.Sprintf("%s.csv", metadata.FileName)

	io.WriteFile(filepath, filename, metadataJSON)
}
