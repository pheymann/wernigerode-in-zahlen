package financialplan

import (
	"wernigerode-in-zahlen.de/internal/pkg/model"
	fd "wernigerode-in-zahlen.de/internal/pkg/model/financialdata"
)

var ()

func DecodeFromAccounts(accounts []fd.Account) model.FinancialPlan {
	var financialPlan = model.FinancialPlan2{}

	for _, account := range accounts {
		account
	}

	return financialPlan
}
