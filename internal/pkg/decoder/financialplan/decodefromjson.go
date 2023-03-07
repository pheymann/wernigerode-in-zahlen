package financialplan

import (
	"encoding/json"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

func DecodeFromJSON(financialPlanJSON string) model.FinancialPlan {
	var financialPlan = model.FinancialPlan{}
	json.Unmarshal([]byte(financialPlanJSON), &financialPlan)

	return financialPlan
}
