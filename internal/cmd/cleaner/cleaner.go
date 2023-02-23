package cleaner

import (
	"bufio"
	"fmt"
	"os"

	decodeFpa "wernigode-in-zahlen.de/internal/pkg/decoder/financeplan_a"
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

func CleanUpFinancePlanA(financeplan_a_file *os.File) model.FinancePlanA {
	rawCSVDecoder := rawcsv.NewDecoder()
	financePlanACostCenterDecoder := decodeFpa.New()
	costCenter := []model.FinancePlanACostCenter{}

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		rawCSVDecoder.Debug()
	// 		fmt.Printf("\n%+v\n", r)
	// 		os.Exit(2)
	// 	}
	// }()

	financePlan_a_Scanner := bufio.NewScanner(financeplan_a_file)

	var rawCSVRows = []model.RawCSVRow{}
	for financePlan_a_Scanner.Scan() {
		line := financePlan_a_Scanner.Text()

		fmt.Println("RawCSVRow")
		rawCSVRows = append(rawCSVRows, rawCSVDecoder.Decode2(line))

		tpe, matches, regex := rawCSVDecoder.Decode(line)

		switch tpe {
		case rawcsv.DecodeTypeAccount:
			financePlan := financePlanACostCenterDecoder.Decode(model.CostCenterGroup, matches, regex)
			costCenter = append(costCenter, financePlan)
		case rawcsv.DecodeTypeUnit:
			financePlan := financePlanACostCenterDecoder.Decode(model.CostCenterUnit, matches, regex)
			costCenter = append(costCenter, financePlan)
		case rawcsv.DeocdeTypeSeparateLine:
			separateLine := matches[regex.SubexpIndex("desc")]
			costCenter[len(costCenter)-1].Desc = costCenter[len(costCenter)-1].Desc + " " + separateLine
		}
	}

	fmt.Println("FinancialPlan")
	fmt.Printf("%+v\n", decodeFpa2.Decode(rawCSVRows))

	financePlan_a := decodeFpa.Decode(costCenter)

	return financePlan_a
}
