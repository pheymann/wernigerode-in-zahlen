package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"wernigode-in-zahlen.de/internal/cmd/departmenthtmlgenerator"
	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	compressedEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/compresseddepartment"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
	html "wernigode-in-zahlen.de/internal/pkg/model/html"
)

var (
	productDirRegex = regexp.MustCompile(`assets/data/processed/\d+/\d+/\d+/\d+/\d+$`)
)

func main() {
	year := model.BudgetYear2023

	department := flag.String("department", "", "department to generate a HTML file from")
	departmentName := flag.String("name", "", "department name")
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	if *department == "" {
		panic("department is required")
	}

	if *departmentName == "" {
		panic("department name is required")
	}

	financialPlanFile, err := os.Open(*debugRootPath + "assets/data/processed/" + *department + "/financial_plan_a.json")
	if err != nil {
		panic(err)
	}
	defer financialPlanFile.Close()

	financialPlan := fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanFile))

	var productData = []html.ProductData{}
	errWalk := filepath.Walk(*debugRootPath+"assets/data/processed/"+*department, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() && productDirRegex.MatchString(path) {
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

			cashflowTotalFile, err := os.Open(path + "/cashflow.txt")
			if err != nil {
				panic(err)
			}
			defer cashflowTotalFile.Close()

			financialPlan := fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanFile))
			metadata := metaDecoder.DecodeFromJSON(io.ReadCompleteFile(metadataFile))
			cashflowTotal, err := strconv.ParseFloat(io.ReadCompleteFile(cashflowTotalFile), 64)
			if err != nil {
				panic(err)
			}

			productData = append(productData, html.ProductData{
				FinancialPlan: financialPlan,
				Metadata:      metadata,
				CashflowTotal: cashflowTotal,
			})
			return nil
		}

		return nil
	})

	if errWalk != nil {
		panic(errWalk)
	}

	compressed := &model.CompressedDepartment{
		DepartmentName: *departmentName,
		ID:             *department,
	}

	io.WriteFile(
		model.TargetFile{
			Path: "assets/html/" + *department + "/",
			Name: "department",
			Tpe:  "html",
		},
		departmenthtmlgenerator.Generate(financialPlan, productData, compressed, year, *debugRootPath),
	)

	io.WriteFile(
		model.TargetFile{
			Path: "assets/data/processed/" + *department + "/",
			Name: "compressed",
			Tpe:  "json",
		},
		compressedEncoder.Encode(*compressed),
	)
}
