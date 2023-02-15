package cleaner

import (
	"bufio"
	"fmt"
	"os"

	"wernigode-in-zahlen.de/internal/pkg/decoder/financeplan_a"
	"wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigode-in-zahlen.de/internal/pkg/decoder/rawcsv"
	encoder "wernigode-in-zahlen.de/internal/pkg/encoder/financeplan_a"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func CleanUp(filename string, file *os.File, debug bool) {
	metadataDecoder := metadata.NewMetadataDecoder()
	rawCSVDecoder := rawcsv.NewDecoder()
	financePlanACostCenterDecoder := financeplan_a.NewFinancePlanACostCenterDecoder()

	metadata := metadataDecoder.Decode(filename)

	scanner := bufio.NewScanner(file)
	costCenter := []model.FinancePlanACostCenter{}

	for scanner.Scan() {
		line := scanner.Text()

		tpe, matches, regex := rawCSVDecoder.Decode(line)
		financePlan := financePlanACostCenterDecoder.Decode(tpe, matches, regex)

		if debug {
			fmt.Printf("------------------------\n%s\n%+v\n\n", line, financePlan)
		}

		costCenter = append(costCenter, financePlan)
	}

	financePlanA := financeplan_a.Decode(costCenter)

	encoder.EncodeAndWriteGroup(financePlanA.Groups, metadata)
	encoder.EncodeAndWriteUnit(financePlanA.Units, metadata)
}
