package financialdata

import (
	"regexp"

	"wernigerode-in-zahlen.de/internal/pkg/decoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	fd "wernigerode-in-zahlen.de/internal/pkg/model/financialdata"
)

var (
	isFinancialPlanAccount = regexp.MustCompile(`^(\d\.)+(\d{2}\.)+[^45]\d+$`)
)

func DecodeAccounts(rows [][]string) map[string][]fd.Account {
	accounts := make(map[string][]fd.Account)

	for _, row := range rows {
		account := fd.Account{
			ID:          row[0],
			ProductID:   row[1],
			Description: row[2],
		}

		if !isFinancialPlanAccount.MatchString(account.ID) {
			continue
		}

		budget := make(map[string]float64)
		budget[model.BudgetYear2022] = decoder.DecodeFloat64(row[8])
		budget[model.BudgetYear2023] = decoder.DecodeFloat64(row[10])
		budget[model.BudgetYear2024] = decoder.DecodeFloat64(row[11])
		budget[model.BudgetYear2025] = decoder.DecodeFloat64(row[12])
		budget[model.BudgetYear2026] = decoder.DecodeFloat64(row[13])

		account.Budget = budget

		if accounts[account.ID] == nil {
			accounts[account.ID] = []fd.Account{account}
		} else {
			accounts[account.ID] = append(accounts[account.ID], account)
		}
	}

	return accounts
}
