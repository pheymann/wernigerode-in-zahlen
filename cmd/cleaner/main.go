package main

import (
	"flag"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
	decodeTarget "wernigode-in-zahlen.de/internal/pkg/decoder/targetfile"
	encodeFpa "wernigode-in-zahlen.de/internal/pkg/encoder/financeplan_a"
	encodeMeta "wernigode-in-zahlen.de/internal/pkg/encoder/metadata"
	writeFpa "wernigode-in-zahlen.de/internal/pkg/io/financeplan_a"
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

		metadata := cleaner.CleanUpMetadata(metadataFile)

		writeMeta.Write(encodeMeta.Encode(metadata), decodeTarget.Decode(metadataFile))
	}

	financePlan_a_file, err := os.Open(*directory + "/financial_plan_a.csv")
	if err != nil {
		panic(err)
	}

	defer financePlan_a_file.Close()

	financePlan_a := cleaner.CleanUpFinancePlanA(financePlan_a_file)

	writeFpa.WriteGroup(encodeFpa.EncodeGroup(financePlan_a.Groups), decodeTarget.Decode(financePlan_a_file))
	writeFpa.WriteUnit(encodeFpa.EncodeUnit(financePlan_a.Units), decodeTarget.Decode(financePlan_a_file))
}
