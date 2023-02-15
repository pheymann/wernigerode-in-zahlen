package financeplan_a

import (
	"fmt"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func EncodeAndWriteGroup(financePlans []model.FinancePlanACostCenter, metadata model.Metadata) {
	content := CSVHeader
	filepath := fmt.Sprintf(
		"assets/data/processed/%s/%s/%s/%s/",
		metadata.ProductClass,
		metadata.ProductDomain,
		metadata.ProductGroup,
		metadata.Product,
	)
	filename := fmt.Sprintf("%s.csv", metadata.FileName)

	for _, financePlan := range financePlans {
		content += toCSVRow(financePlan) + "\n"
	}

	writeFile(filepath, filename, content)
}
