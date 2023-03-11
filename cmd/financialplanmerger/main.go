package main

import (
	"flag"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/financialplanmerger"
	decodeTarget "wernigode-in-zahlen.de/internal/pkg/decoder/targetfile"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/shared"
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

	financialPlanBJSONOpt := shared.Option[string]{IsSome: false}
	financialPlanBFile, err := os.Open(*debugRootPath + *directory + "/financial_plan_b.json")
	if err == nil {
		financialPlanBJSONOpt.ToSome(io.ReadCompleteFile(financialPlanBFile))

		defer financialPlanBFile.Close()
	}

	merged := financialplanmerger.Merge(
		io.ReadCompleteFile(financialPlanAFile),
		financialPlanBJSONOpt,
	)

	target := decodeTarget.Decode(financialPlanAFile, "data/processed")
	target.Name = "merged_financial_plan"
	target.Tpe = "json"

	io.WriteFile(target, merged)
}
