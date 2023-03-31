package financialdatacleaner

import (
	"encoding/csv"
	"log"
	"os"

	fd "wernigerode-in-zahlen.de/internal/pkg/decoder/financialdata"
)

func Cleanup(financialDataFile *os.File) map[string]string {
	csvReader := csv.NewReader(financialDataFile)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse financial data CSV", err)
	}

	productAccounts := fd.DecodeAccounts(rows)

	return map[string]string{}
}
