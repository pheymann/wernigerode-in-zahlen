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
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	financialPlanFile, err := os.Open(*directory + "/merged_financial_plan.json")
	if err != nil {
		panic(err)
	}
	defer financialPlanFile.Close()

	metadataFile, err := os.Open(*directory + "/metadata.json")
	if err != nil {
		panic(err)
	}
	defer metadataFile.Close()

	productHtml := htmlgenerator.GenerateProductHTML(
		io.ReadCompleteFile(financialPlanFile),
		io.ReadCompleteFile(metadataFile),
		model.BudgetYear2023,
		*debugRootPath,
	)

	target := decodeTarget.Decode(financialPlanFile, "html")
	target.Name = "product"
	target.Tpe = "html"

	io.WriteFile(target, productHtml)
}
