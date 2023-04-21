package main

import (
	"flag"
	"html/template"
	"os"

	"wernigerode-in-zahlen.de/internal/cmd/htmlgenerator"
	fpDecoder "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	"wernigerode-in-zahlen.de/internal/pkg/io"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func main() {
	budgetYear := model.BudgetYear2023

	debugRootPath := flag.String("debug-root-path", "", "Debug: root path")

	flag.Parse()

	financialCityDataFile, err := os.Open(*debugRootPath + "assets/data/processed/financial_data.json")
	if err != nil {
		panic(err)
	}
	defer financialCityDataFile.Close()

	financialCityData := fpDecoder.DecodeFromJSON2(io.ReadCompleteFile(financialCityDataFile))

	overviewTmpl := template.Must(template.ParseFiles(*debugRootPath + "assets/html/templates/overview.template.html"))
	overviewFile, overviewContent := htmlgenerator.GenerateOverview(financialCityData, budgetYear, overviewTmpl)

	overviewFile.Path = *debugRootPath + overviewFile.Path

	io.WriteFile(overviewFile, overviewContent)

	departmentTmpl := template.Must(template.ParseFiles(*debugRootPath + "assets/html/templates/department.template.html"))
	departmentPairs := htmlgenerator.GenerateDepartments(financialCityData, budgetYear, departmentTmpl)

	for _, pair := range departmentPairs {
		pair.First.Path = *debugRootPath + pair.First.Path

		io.WriteFile(pair.First, pair.Second)
	}
}
