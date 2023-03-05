package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	fpaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	productDirRegex = regexp.MustCompile(`^assets/data/processed/\d+/\d+/\d+/\d+/\d+$`)
)

func main() {
	department := flag.String("department", "", "department to generate a HTML file from")

	flag.Parse()

	if *department == "" {
		panic("department is required")
	}

	var productData = []ProductData{}

	err := filepath.Walk("assets/data/processed/"+*department, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() && productDirRegex.MatchString(path) {
			fmt.Printf("Read %s\n", path)

			financialPlanAFile, err := os.Open(path + "/financial_plan_a.json")
			if err != nil {
				panic(err)
			}
			defer financialPlanAFile.Close()

			metadataFile, err := os.Open(path + "/metadata.json")
			if err != nil {
				panic(err)
			}
			defer metadataFile.Close()

			financialPlanA := fpaDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanAFile))
			metadata := metaDecoder.DecodeFromJSON(io.ReadCompleteFile(metadataFile))

			productData = append(productData, ProductData{
				FinancialPlanA: financialPlanA,
				Metadata:       metadata,
			})
			return nil
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	// for _, product := range productData {

	// }

	fmt.Printf("%+v\n", productData)
}

type ProductData struct {
	FinancialPlanA model.FinancialPlanA
	Metadata       model.Metadata
}
