package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"wernigode-in-zahlen.de/internal/cmd/departmenthtmlgenerator"
	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	compressedEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/compresseddepartment"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
	html "wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

var (
	productDirRegex = regexp.MustCompile(`assets/data/processed/\d+/\d+/\d+/\d+/\d+$`)
)

func main() {
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

	var productData = []html.ProductData{}
	errWalk := filepath.Walk(*debugRootPath+"assets/data/processed/"+*department, func(path string, info os.FileInfo, err error) error {
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

			var financialPlanBJSONOpt = shared.None[string]()
			financialPlanBFile, err := os.Open(path + "/financial_plan_b.json")
			if err == nil {
				defer financialPlanAFile.Close()

				financialPlanBJSONOpt = shared.Some(io.ReadCompleteFile(financialPlanBFile))
			}

			metadataFile, err := os.Open(path + "/metadata.json")
			if err != nil {
				panic(err)
			}
			defer metadataFile.Close()

			financialPlanA := fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanAFile))
			financialPlanBOpt := shared.Map(financialPlanBJSONOpt, func(financialPlanBJSON string) model.FinancialPlan {
				return fpDecoder.DecodeFromJSON(financialPlanBJSON)
			})
			metadata := metaDecoder.DecodeFromJSON(io.ReadCompleteFile(metadataFile))

			productData = append(productData, html.ProductData{
				FinancialPlanA:    financialPlanA,
				FinancialPlanBOpt: financialPlanBOpt,
				Metadata:          metadata,
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
		departmenthtmlgenerator.GenerateDepartmentHTML(productData, compressed, *debugRootPath),
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
