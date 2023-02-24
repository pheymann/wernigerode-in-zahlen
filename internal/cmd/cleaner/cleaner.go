package cleaner

import (
	"bufio"
	"fmt"
	"os"

	decodeFpa2 "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
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

	metadata := metadataDecoder.Decode(metadataLines)

	return metadata
}

func CleanUpFinancialPlanA(financeplan_a_file *os.File) model.FinancialPlanA {
	rawCSVDecoder := rawcsv.NewDecoder()

	defer func() {
		if r := recover(); r != nil {
			rawCSVDecoder.Debug()
			fmt.Printf("\n%+v\n", r)
			os.Exit(2)
		}
	}()

	financePlan_a_Scanner := bufio.NewScanner(financeplan_a_file)

	var rawCSVRows = []model.RawCSVRow{}
	for financePlan_a_Scanner.Scan() {
		line := financePlan_a_Scanner.Text()

		rawCSVRows = append(rawCSVRows, rawCSVDecoder.Decode(line))
	}

	return decodeFpa2.Decode(rawCSVRows)
}
