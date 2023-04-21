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
	file, content := htmlgenerator.GenerateOverview(financialCityData, budgetYear, overviewTmpl)

	file.Path = *debugRootPath + file.Path

	io.WriteFile(file, content)
}
