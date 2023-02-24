package financialplan_a

import (
	"encoding/json"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Encode(financialPlanA model.FinancialPlanA) string {
	bytes, err := json.MarshalIndent(financialPlanA, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
