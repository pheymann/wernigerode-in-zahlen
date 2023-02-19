package cleaner

import (
	"bufio"
	"os"

	decodeFpa "wernigode-in-zahlen.de/internal/pkg/decoder/financeplan_a"
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

	metadataDecoder.Debug()

	metadata := metadataDecoder.Decode(metadataLines)

	return metadata
}

func CleanUpFinancePlanA(financeplan_a_file *os.File) model.FinancePlanA {
	rawCSVDecoder := rawcsv.NewDecoder()
	financePlanACostCenterDecoder := decodeFpa.NewFinancePlanACostCenterDecoder()
	costCenter := []model.FinancePlanACostCenter{}

	rawCSVDecoder.Debug()

	financePlan_a_Scanner := bufio.NewScanner(financeplan_a_file)

	for financePlan_a_Scanner.Scan() {
		line := financePlan_a_Scanner.Text()

		tpe, matches, regex := rawCSVDecoder.Decode(line)

		switch tpe {
		case rawcsv.DecodeTypeGroup:
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

	financePlan_a := decodeFpa.Decode(costCenter)

	return financePlan_a
}
