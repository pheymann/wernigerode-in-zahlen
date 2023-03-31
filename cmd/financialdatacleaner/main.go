package main

import (
	"flag"
	"os"
)

func main() {
	financialDataFilePath := flag.String("file", "", "financial data file")

	flag.Parse()

	if *financialDataFilePath == "" {
		panic("file is required")
	}

	financialDataFile, err := os.Open(*financialDataFilePath)
	if err != nil {
		panic(err)
	}
	defer financialDataFile.Close()

}
