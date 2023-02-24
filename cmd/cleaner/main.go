package main

import (
	"flag"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
	decodeTarget "wernigode-in-zahlen.de/internal/pkg/decoder/targetfile"
	encodeFpa "wernigode-in-zahlen.de/internal/pkg/encoder/financialplan_a"
	encodeMeta "wernigode-in-zahlen.de/internal/pkg/encoder/metadata"
	writeFpa "wernigode-in-zahlen.de/internal/pkg/io/financialplan_a"
	writeMeta "wernigode-in-zahlen.de/internal/pkg/io/metadata"
)

func main() {
	directory := flag.String("dir", "", "directory to clean up")
	existsMetadata := flag.Bool("metadata", false, "exists metadata file")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	if *existsMetadata {
		metadataFile, err := os.Open(*directory + "/metadata.csv")
		if err != nil {
			panic(err)
		}

		defer metadataFile.Close()

		writeMeta.Write(
			encodeMeta.Encode(
				cleaner.CleanUpMetadata(metadataFile),
			),
			decodeTarget.Decode(metadataFile),
		)
	}

	financePlan_a_file, err := os.Open(*directory + "/financial_plan_a.csv")
	if err != nil {
		panic(err)
	}

	defer financePlan_a_file.Close()

	writeFpa.Write(
		encodeFpa.Encode(
			cleaner.CleanUpFinancialPlanA(financePlan_a_file),
		),
		decodeTarget.Decode(financePlan_a_file),
	)
}
