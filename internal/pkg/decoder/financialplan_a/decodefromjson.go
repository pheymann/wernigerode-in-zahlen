package financialplan_a

import (
	"encoding/json"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func DecodeFromJSON(financialPlanAJSON string) model.FinancialPlan {
	var financialPlanA = model.FinancialPlan{}
	json.Unmarshal([]byte(financialPlanAJSON), &financialPlanA)

	return financialPlanA
}
