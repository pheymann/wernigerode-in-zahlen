package filewriter

import (
	"fmt"

	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func WriteGroup(financePlans string, metadata model.Metadata) {
	filepath := fmt.Sprintf(
		"assets/data/processed/%s/%s/%s/%s/",
		metadata.ProductClass,
		metadata.ProductDomain,
		metadata.ProductGroup,
		metadata.Product,
	)
	filename := fmt.Sprintf("%s.csv", metadata.FileName)

	io.WriteFile(filepath, filename, financePlans)
}

func WriteUnit(financePlans map[string]string, metadata model.Metadata) {
	for costCenterGroup, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		filepath := fmt.Sprintf(
			"assets/data/processed/%s/%s/%s/%s/%s/",
			metadata.ProductClass,
			metadata.ProductDomain,
			metadata.ProductGroup,
			metadata.Product,
			costCenterGroup,
		)
		filename := fmt.Sprintf("%s.csv", metadata.FileName)

		io.WriteFile(filepath, filename, financePlans)
	}
}
