package main

import (
	"fmt"
	"os"

	"wernigode-in-zahlen.de/internal/cmd/cleaner"
	encodeFpa "wernigode-in-zahlen.de/internal/pkg/encoder/financeplan_a"
	encodeMeta "wernigode-in-zahlen.de/internal/pkg/encoder/metadata"
)

var (
	debug = false
)

func main() {
	directory := os.Args[1]
	metadataFile, err := os.Open(directory + "/metadata.csv")
	if err != nil {
		panic(err)
	}

	defer metadataFile.Close()

	financePlan_a_file, err := os.Open(directory + "/data.csv")
	if err != nil {
		panic(err)
	}

	defer financePlan_a_file.Close()

	metadata, financePlan_a := cleaner.CleanUp(metadataFile, financePlan_a_file, debug)

	fmt.Printf("%s\n%s\n%s", encodeMeta.Encode(metadata), encodeFpa.EncodeGroup(financePlan_a.Groups, metadata), encodeFpa.EncodeUnit(financePlan_a.Units, metadata))
}
