package financeplan_a

import (
	"fmt"
	"os"

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

func writeFile(filepath string, filename string, content string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.MkdirAll(filepath, 0700)
	}

	file, err := os.Create(filepath + filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	file.Sync()
}
