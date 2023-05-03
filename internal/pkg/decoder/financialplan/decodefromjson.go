package financialplan

import (
	"encoding/json"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func DecodeFromJSON2(financialPlanJSON string) model.FinancialPlanCity {
	var financialPlan = model.FinancialPlanCity{}
	json.Unmarshal([]byte(financialPlanJSON), &financialPlan)

	return financialPlan
}
