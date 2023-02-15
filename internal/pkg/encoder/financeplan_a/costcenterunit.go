package financeplan_a

import (
	"fmt"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func EncodeAndWriteUnit(financePlans map[string][]model.FinancePlanACostCenter, metadata model.Metadata) {
	for costCenterGroup, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		content := CSVHeader
		filepath := fmt.Sprintf(
			"assets/data/processed/%s/%s/%s/%s/%s/",
			metadata.ProductClass,
			metadata.ProductDomain,
			metadata.ProductGroup,
			metadata.Product,
			costCenterGroup,
		)
		filename := fmt.Sprintf("%s.csv", metadata.FileName)

		for _, financePlan := range financePlans {
			content += toCSVRow(financePlan) + "\n"
		}

		writeFile(filepath, filename, content)
	}
}
