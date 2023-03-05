package main

import (
	"flag"
	"os"

	htmlgenerator "wernigode-in-zahlen.de/internal/cmd/producthtmlgenerator"
	decodeTarget "wernigode-in-zahlen.de/internal/pkg/decoder/targetfile"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func main() {
	directory := flag.String("dir", "", "directory to read HTML files from")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	financialPlanAFile, err := os.Open(*directory + "/financial_plan_a.json")
	if err != nil {
		panic(err)
	}
	defer financialPlanAFile.Close()

	metadataFile, err := os.Open(*directory + "/metadata.json")
	if err != nil {
		panic(err)
	}
	defer metadataFile.Close()

	productHtml := htmlgenerator.GenerateHTMLForProduct(
		io.ReadCompleteFile(financialPlanAFile),
		io.ReadCompleteFile(metadataFile),
		model.BudgetYear2023,
	)

	productHtmlFile, err := os.Create("product.html")
	if err != nil {
		panic(err)
	}
	defer productHtmlFile.Close()

	target := decodeTarget.Decode(financialPlanAFile, "html")
	target.Name = "product"
	target.Tpe = "html"

	io.WriteFile(target, productHtml)
}
