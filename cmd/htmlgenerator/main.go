package main

import (
	"bufio"
	"flag"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/htmlgenerator"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func main() {
	directory := flag.String("dir", "", "directory to read HTML files from")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	financialPlanAFile, err := os.Open(*directory + "financial_plan_a.json")
	if err != nil {
		panic(err)
	}
	defer financialPlanAFile.Close()

	metadataFile, err := os.Open(*directory + "metadata.json")
	if err != nil {
		panic(err)
	}
	defer metadataFile.Close()

	productHtml := htmlgenerator.GenerateHTMLForProduct(
		readCompleteFile(financialPlanAFile),
		readCompleteFile(metadataFile),
		model.BudgetYear2023,
	)

	productHtmlFile, err := os.Create("product.html")
	if err != nil {
		panic(err)
	}
	defer productHtmlFile.Close()

	productHtmlFile.WriteString(productHtml)
}

func readCompleteFile(file *os.File) string {
	scanner := bufio.NewScanner(file)

	var content = ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	return content
}
