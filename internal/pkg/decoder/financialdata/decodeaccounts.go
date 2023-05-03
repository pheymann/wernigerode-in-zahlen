package financialdata

import (
	"wernigerode-in-zahlen.de/internal/pkg/decoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	fd "wernigerode-in-zahlen.de/internal/pkg/model/financialdata"
)

func DecodeAccounts(rows [][]string) map[string][]fd.Account {
	accounts := make(map[string][]fd.Account)

	for _, row := range rows {
		account := fd.Account{
			ID:          row[0],
			ProductID:   row[1],
			Description: row[2],
		}

		budget := make(map[string]float64)
		budget[model.BudgetYear2022] = decoder.DecodeGermanFloat(row[8])
		budget[model.BudgetYear2023] = decoder.DecodeGermanFloat(row[10])
		budget[model.BudgetYear2024] = decoder.DecodeGermanFloat(row[11])
		budget[model.BudgetYear2025] = decoder.DecodeGermanFloat(row[12])
		budget[model.BudgetYear2026] = decoder.DecodeGermanFloat(row[13])

		account.Budget = budget

		if accounts[account.ProductID] == nil {
			accounts[account.ProductID] = []fd.Account{account}
		} else {
			accounts[account.ProductID] = append(accounts[account.ProductID], account)
		}
	}

	return accounts
}
