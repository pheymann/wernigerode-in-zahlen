package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
)

func main() {
	directory := flag.String("dir", "", "directory to read HTML files from")

	flag.Parse()

	if *directory == "" {
		panic("directory is required")
	}

	tpeExtractionRegex := regexp.MustCompile(`^assets/data/processed/(?P<type>(\d+/)+)`)
	departmentDetectionRegex := regexp.MustCompile(`^\d+/$`)
	productDetectionRegex := regexp.MustCompile(`^\d+/\d+/\d+/\d+/\d+/$`)

	matches := tpeExtractionRegex.FindStringSubmatch(*directory)
	tpe := decoder.DecodeString(tpeExtractionRegex, "type", matches)

	if departmentDetectionRegex.MatchString(tpe) {
		println("=== DEPARTMENT ===")
		financialPlan_a, err := os.Open(*directory + "financial_plan_a.csv")
		if err != nil {
			panic(err)
		}
		defer financialPlan_a.Close()

		// parse file

	} else if productDetectionRegex.MatchString(tpe) {
		println("product")
	} else {
		panic(fmt.Sprintf("unknown type: %s", tpe))
	}
}
