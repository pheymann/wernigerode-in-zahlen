package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"wernigode-in-zahlen.de/internal/cmd/overviewhtmlgenerator"
	compressedDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/compresseddepartment"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func main() {
	departmentIDsRaw := flag.String("departments", "", "list of departments")
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	if *departmentIDsRaw == "" {
		panic("list of departments is required")
	}

	departmentIDs := strings.Split(*departmentIDsRaw, ",")
	if len(departmentIDs) == 0 {
		panic("empty list of departments")
	}

	var departments = []model.CompressedDepartment{}
	for _, departmentID := range departmentIDs {
		compressedFile, err := os.Open(fmt.Sprintf("%sassets/data/processed/%s/compressed.json", *debugRootPath, departmentID))
		if err != nil {
			panic(err)
		}
		defer compressedFile.Close()

		departments = append(departments, compressedDecoder.Decode(io.ReadCompleteFile(compressedFile)))
	}

	io.WriteFile(
		model.TargetFile{
			Path: "assets/html/",
			Name: "index",
			Tpe:  "html",
		},
		overviewhtmlgenerator.Generate(departments, *debugRootPath),
	)
}
