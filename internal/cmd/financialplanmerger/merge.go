package financialplanmerger

import (
	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	fpEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/financialplan"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func Merge(fpaJSON string, fpbJSONOpt shared.Option[string]) string {
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
