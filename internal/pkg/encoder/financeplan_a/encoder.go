package financeplan_a

import (
	"fmt"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

const (
	CSVHeader = "id;desc;_2020;_2021;_2022;_2023;_2024;_2025\n"
)

func toCSVRow(financePlan model.FinancePlanACostCenter) string {
	return fmt.Sprintf(
		"%s;%s;%f;%f;%f;%f;%f;%f",
		financePlan.Id,
		financePlan.Desc,
		financePlan.Budget2020,
		financePlan.Budget2021,
		financePlan.Budget2022,
		financePlan.Budget2023,
		financePlan.Budget2024,
		financePlan.Budget2025,
	)
}

func EncodeGroup(financePlans []model.FinancePlanACostCenter) string {
	content := CSVHeader

	for _, financePlan := range financePlans {
		content += toCSVRow(financePlan) + "\n"
	}

	return content
}

func EncodeUnit(financePlans map[string][]model.FinancePlanACostCenter) map[string]string {
	groupBasedUnits := map[string]string{}

	for costCenterGroup, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		content := CSVHeader

		for _, financePlan := range financePlans {
			content += toCSVRow(financePlan) + "\n"
		}
		groupBasedUnits[costCenterGroup] = content
	}

	return groupBasedUnits
}
