package financialplanmerger

import (
	fpDecoder "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	fpEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/financialplan"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
)

func Merge(fpaJSON string) string {
	fpa := fpDecoder.DecodeFromJSON(fpaJSON)

	fixBudgetSigns(&fpa)

	return fpEncoder.Encode(fpa)
}

func fixBudgetSigns(fp *model.FinancialPlan) {
	for balanceIndex, balance := range fp.Balances {
		for accountIndex, account := range balance.Accounts {
			accountClass := html.ClassifyAccount(account)
			fp.Balances[balanceIndex].Accounts[accountIndex].Budgets = updateBudges(account.Budgets, accountClass)

			for subIndex, sub := range account.Subs {
				fp.Balances[balanceIndex].Accounts[accountIndex].Subs[subIndex].Budgets = updateBudges(sub.Budgets, accountClass)

				for unitIndex, unit := range sub.Units {
					fp.Balances[balanceIndex].Accounts[accountIndex].Subs[subIndex].Units[unitIndex].Budgets = updateBudges(unit.Budgets, accountClass)
				}
			}
		}
	}
}

func updateBudges(budgets map[string]float64, class html.CashflowClass) map[string]float64 {
	for year, budget := range budgets {
		if class == html.CashflowClassExpenses {
			budgets[year] = -budget
		}
	}
	return budgets
}
