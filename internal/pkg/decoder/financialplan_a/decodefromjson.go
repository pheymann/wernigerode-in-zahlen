package financialplan_a

import (
	"encoding/json"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func DecodeFromJSON(financialPlanAJSON string) model.FinancialPlanA {
	var financialPlanA = model.FinancialPlanA{}
	json.Unmarshal([]byte(financialPlanAJSON), &financialPlanA)

	return financialPlanA
}
