package cleaner

import (
	"bufio"
	"fmt"
	"os"

	decodeFpa "wernigode-in-zahlen.de/internal/pkg/decoder/financeplan_a"
	decodeMeta "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigode-in-zahlen.de/internal/pkg/decoder/rawcsv"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func CleanUp(metadataFile *os.File, financeplan_a_file *os.File, debug bool) (model.Metadata, model.FinancePlanA) {
	metadataScanner := bufio.NewScanner(metadataFile)
	metadataLines := []string{}

	for metadataScanner.Scan() {
		metadataLines = append(metadataLines, metadataScanner.Text())
	}

	metadataDecoder := decodeMeta.NewMetadataDecoder()
	metadata := metadataDecoder.Decode(metadataLines)

	rawCSVDecoder := rawcsv.NewDecoder()
	financePlanACostCenterDecoder := decodeFpa.NewFinancePlanACostCenterDecoder()
	costCenter := []model.FinancePlanACostCenter{}

	financePlan_a_Scanner := bufio.NewScanner(financeplan_a_file)

	for financePlan_a_Scanner.Scan() {
		line := financePlan_a_Scanner.Text()

		tpe, matches, regex := rawCSVDecoder.Decode(line)
		financePlan := financePlanACostCenterDecoder.Decode(tpe, matches, regex)

		if debug {
			fmt.Printf("------------------\n%s\n%+v\n", line, financePlan)
		}
		costCenter = append(costCenter, financePlan)
	}

	financePlan_a := decodeFpa.Decode(costCenter)

	return metadata, financePlan_a
}
