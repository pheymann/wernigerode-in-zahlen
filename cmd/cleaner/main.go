package main

import (
	"flag"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
	decodeTarget "wernigode-in-zahlen.de/internal/pkg/decoder/targetfile"
	encodeFp "wernigode-in-zahlen.de/internal/pkg/encoder/financialplan"
	encodeMeta "wernigode-in-zahlen.de/internal/pkg/encoder/metadata"
	writeFp "wernigode-in-zahlen.de/internal/pkg/io/financialplan"
	writeMeta "wernigode-in-zahlen.de/internal/pkg/io/metadata"
)

func main() {
	directory := flag.String("dir", "", "directory to clean up")
	tpe := flag.String("type", "", "type of financial plan (product, department)")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}
	if *tpe == "" {
		panic("type is required")
	}

	if *tpe == "product" {
		metadataFile, err := os.Open(*directory + "/metadata.csv")
		if err != nil {
			panic(err)
		}

		defer metadataFile.Close()

		writeMeta.Write(
			encodeMeta.Encode(
				cleaner.CleanUpMetadata(metadataFile),
			),
			decodeTarget.Decode(metadataFile, "data/processed"),
		)

		financialPlanBFile, err := os.Open(*directory + "/financial_plan_b.csv")
		if err == nil {
			defer financialPlanBFile.Close()

			writeFp.Write(
				encodeFp.Encode(
					cleaner.CleanUpFinancialPlanB(financialPlanBFile),
				),
				decodeTarget.Decode(financialPlanBFile, "data/processed"),
			)
		}
	}

	financialPlanAFile, err := os.Open(*directory + "/financial_plan_a.csv")
	if err != nil {
		panic(err)
	}

	defer financialPlanAFile.Close()

	writeFp.Write(
		encodeFp.Encode(
			cleaner.CleanUpFinancialPlanA(financialPlanAFile),
		),
		decodeTarget.Decode(financialPlanAFile, "data/processed"),
	)
}
