package financialplan_a

import (
	"encoding/json"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Encode(financialPlan model.FinancialPlan) string {
	bytes, err := json.MarshalIndent(financialPlan, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func Encode2(financialPlan model.FinancialPlan2) string {
	bytes, err := json.MarshalIndent(financialPlan, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
