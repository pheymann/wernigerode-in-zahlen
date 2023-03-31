package main

import (
	"flag"
	"os"

	"wernigerode-in-zahlen.de/internal/cmd/financialplanmerger"
	decodeTarget "wernigerode-in-zahlen.de/internal/pkg/decoder/targetfile"
	"wernigerode-in-zahlen.de/internal/pkg/io"
)

func main() {
	directory := flag.String("dir", "", "directory to read HTML files from")
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	financialPlanAFile, err := os.Open(*debugRootPath + *directory + "/financial_plan_a.json")
	if err != nil {
		panic(err)
	}
	defer financialPlanAFile.Close()

	merged := financialplanmerger.Merge(
		io.ReadCompleteFile(financialPlanAFile),
	)

	target := decodeTarget.Decode(financialPlanAFile, "data/processed")
	target.Name = "merged_financial_plan"
	target.Tpe = "json"

	io.WriteFile(target, merged)
}
