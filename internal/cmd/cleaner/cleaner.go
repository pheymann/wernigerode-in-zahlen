package cleaner

import (
	"bufio"
	"fmt"
	"os"

	decodeFpa "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
	decodeFpb "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_b"
	decodeMeta "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigode-in-zahlen.de/internal/pkg/decoder/rawcsv"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func CleanUpMetadata(metadataFile *os.File) model.Metadata {
	metadataScanner := bufio.NewScanner(metadataFile)
	metadataLines := []string{}

	for metadataScanner.Scan() {
		metadataLines = append(metadataLines, metadataScanner.Text())
	}

	metadataDecoder := decodeMeta.NewMetadataDecoder()

	defer func() {
		if r := recover(); r != nil {
			metadataDecoder.Debug()
			fmt.Printf("\n%+v\n", r)
			os.Exit(1)
		}
	}()

	metadata := metadataDecoder.DecodeFromCSV(metadataLines)

	return metadata
}

func CleanUpFinancialPlanA(financialPlaAFile *os.File) model.FinancialPlan {
	rawCSVDecoder := rawcsv.NewDecoder()

	defer func() {
		if r := recover(); r != nil {
			rawCSVDecoder.Debug()
			fmt.Printf("\n%+v\n", r)
			os.Exit(2)
		}
	}()

	financialPlanAScanner := bufio.NewScanner(financialPlaAFile)

	var rawCSVRows = []model.RawCSVRow{}
	for financialPlanAScanner.Scan() {
		line := financialPlanAScanner.Text()

		rawCSVRows = append(rawCSVRows, rawCSVDecoder.Decode(line, false))
	}

	return decodeFpa.DecodeFromCSV(rawCSVRows)
}

func CleanUpFinancialPlanB(financialPlaAFile *os.File) model.FinancialPlan {
	rawCSVDecoder := rawcsv.NewDecoder()

	defer func() {
		if r := recover(); r != nil {
			rawCSVDecoder.Debug()
			fmt.Printf("\n%+v\n", r)
			os.Exit(2)
		}
	}()

	financialPlanAScanner := bufio.NewScanner(financialPlaAFile)

	var rawCSVRows = []model.RawCSVRow{}
	for financialPlanAScanner.Scan() {
		line := financialPlanAScanner.Text()

		rawCSVRows = append(rawCSVRows, rawCSVDecoder.Decode(line, true))
	}

	return decodeFpb.DecodeFromCSV(rawCSVRows)
}
