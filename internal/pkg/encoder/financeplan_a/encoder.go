package financeplan_a

import (
	"fmt"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	CSVHeader = []string{"id", "desc", "_2020", "_2021", "_2022", "_2023", "_2024", "_2025"}
)

func toCSVRow(financePlan model.FinancePlanACostCenter) []string {
	return []string{
		financePlan.Id,
		financePlan.Desc,
		fmt.Sprintf("%f", financePlan.Budget2020),
		fmt.Sprintf("%f", financePlan.Budget2021),
		fmt.Sprintf("%f", financePlan.Budget2022),
		fmt.Sprintf("%f", financePlan.Budget2023),
		fmt.Sprintf("%f", financePlan.Budget2024),
		fmt.Sprintf("%f", financePlan.Budget2025),
	}
}

func EncodeGroup(financePlans []model.FinancePlanACostCenter) [][]string {
	var content = [][]string{CSVHeader}

	for _, financePlan := range financePlans {
		content = append(content, toCSVRow(financePlan))
	}

	return content
}

func EncodeUnit(financePlans map[string][]model.FinancePlanACostCenter) map[string][][]string {
	groupBasedUnits := map[string][][]string{}

	for costCenterGroup, financePlans := range financePlans {
		if len(financePlans) == 0 {
			continue
		}

		var content = [][]string{CSVHeader}

		for _, financePlan := range financePlans {
			content = append(content, toCSVRow(financePlan))
		}
		groupBasedUnits[costCenterGroup] = content
	}

	return groupBasedUnits
}
