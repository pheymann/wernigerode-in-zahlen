package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	htmlgenerator "wernigode-in-zahlen.de/internal/cmd/producthtmlgenerator"
	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	decodeTarget "wernigode-in-zahlen.de/internal/pkg/decoder/targetfile"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
)

var (
	subProductDirRegex = regexp.MustCompile(`assets/data/processed/\d+/\d+/\d+/\d+/\d+/\d+$`)
)

func main() {
	directory := flag.String("dir", "", "directory to read HTML files from")
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	financialPlanFile, err := os.Open(*debugRootPath + *directory + "/merged_financial_plan.json")
	if err != nil {
		panic(err)
	}
	defer financialPlanFile.Close()

	metadataFile, err := os.Open(*debugRootPath + *directory + "/metadata.json")
	if err != nil {
		panic(err)
	}
	defer metadataFile.Close()

	var subProductData = []html.ProductData{}
	errWalk := filepath.Walk(*debugRootPath+*directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() && subProductDirRegex.MatchString(path) {
			fmt.Printf("Read %s\n", path)

			financialPlanFile, err := os.Open(path + "/merged_financial_plan.json")
			if err != nil {
				panic(err)
			}
			defer financialPlanFile.Close()

			metadataFile, err := os.Open(path + "/metadata.json")
			if err != nil {
				panic(err)
			}
			defer metadataFile.Close()

			financialPlan := fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanFile))
			metadata := metaDecoder.DecodeFromJSON(io.ReadCompleteFile(metadataFile))

			subProductData = append(subProductData, html.ProductData{
				FinancialPlan: financialPlan,
				Metadata:      metadata,
			})
			return nil
		}

		return nil
	})

	if errWalk != nil {
		panic(errWalk)
	}

	productHtml := htmlgenerator.Generate(
		io.ReadCompleteFile(financialPlanFile),
		io.ReadCompleteFile(metadataFile),
		subProductData,
		model.BudgetYear2023,
		*debugRootPath,
	)

	target := decodeTarget.Decode(financialPlanFile, "html")
	target.Name = "product"
	target.Tpe = "html"

	io.WriteFile(target, productHtml)
}
