package main

import (
	"flag"
	"os"

	"wernigerode-in-zahlen.de/internal/cmd/metadatacleaner"
	decodeTarget "wernigerode-in-zahlen.de/internal/pkg/decoder/targetfile"
	encodeMeta "wernigerode-in-zahlen.de/internal/pkg/encoder/metadata"
	writeMeta "wernigerode-in-zahlen.de/internal/pkg/io/metadata"
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
				metadatacleaner.Cleanup(metadataFile),
			),
			decodeTarget.Decode(metadataFile, "data/processed"),
		)
	}
}
