package main

import (
	"flag"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/htmlgenerator"
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

	htmlgenerator.GenerateHTMLForProduct(financialPlanAFile, metadataFile)
}
