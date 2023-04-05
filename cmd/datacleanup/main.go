package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"wernigerode-in-zahlen.de/internal/cmd/datacleanup"
	writeMeta "wernigerode-in-zahlen.de/internal/pkg/io/metadata"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

var (
	productDirRegex = regexp.MustCompile(`^.+\d/\d/\d/\d/\d{2}(/\d{2})?$`)
)

func main() {
	debugRootPath := flag.String("debug-root-path", "", "Debug: root path")

	flag.Parse()

	// collect metadata.csv files
	var metadataFiles = []*os.File{}

	fmt.Println(*debugRootPath + "assets/data/raw")
	errWalk := filepath.Walk(*debugRootPath+"assets/data/raw", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() && productDirRegex.MatchString(path) {
			fmt.Printf("Read %s\n", path)

			metadataFile, err := os.Open(path + "/metadata.csv")
			if err != nil {
				panic(err)
			}

			metadataFiles = append(metadataFiles, metadataFile)
		}

		return nil
	})

	if errWalk != nil {
		panic(errWalk)
	}

	financialDataFile, err := os.Open(*debugRootPath + "assets/data/financial_data.csv")
	if err != nil {
		panic(err)
	}
	defer financialDataFile.Close()

	writeMeta.Write(
		datacleanup.Cleanup(financialDataFile, metadataFiles),
		model.TargetFile{
			Path: *debugRootPath + "assets/data/processed/",
			Name: "financial_data",
			Tpe:  "json",
		},
	)

	for _, file := range metadataFiles {
		file.Close()
	}
}
