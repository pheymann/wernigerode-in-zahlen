package financialplanmerger

import (
	"fmt"

	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	fpEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/financialplan"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func Merge(fpaJSON string, fpbJSONOpt shared.Option[string]) string {
	fpa := fpDecoder.DecodeFromJSON(fpaJSON)
	fpbOpt := shared.Map(fpbJSONOpt, func(fpbJSON string) model.FinancialPlan {
		return fpDecoder.DecodeFromJSON(fpbJSON)
	})

	if fpbOpt.IsSome {
		valueLimits := fpbToAboveValueLimits(fpbOpt.Value)

		for balanceIndex, fpaBalance := range fpa.Balances {
			for accountIndex, fpaAccount := range fpaBalance.Accounts {
				for subIndex, fpaSubAccount := range fpaAccount.Subs {
					for unitIndex, fpaUnit := range fpaSubAccount.Units {
						for index, valueLimit := range valueLimits {
							if valueLimit != nil && fpaUnit.Id == valueLimit.ID {
								fpa.Balances[balanceIndex].Accounts[accountIndex].Subs[subIndex].Units[unitIndex].AboveValueLimit = valueLimit
								valueLimits[index] = nil
								break
							}
						}
					}
				}
			}
		}

		var errorMessages string = "Not all value limits were used.\n"
		var errorCounter = 0
		for _, valueLimit := range valueLimits {
			if valueLimit == nil {
				continue
			}
			errorMessages += fmt.Sprintf("%+v\n", valueLimit)
			errorCounter++
		}
		if errorCounter > 0 {
			panic(errorMessages)
		}
	}

	fixBudgetSigns(&fpa)

	return fpEncoder.Encode(fpa)
}

func fpbToAboveValueLimits(fpb model.FinancialPlan) []*model.AboveValueLimit {
	var aboveValueLimits = []*model.AboveValueLimit{}

	for _, balance := range fpb.Balances {
		category := balance.Desc
		for _, account := range balance.Accounts {
			subCategory := account.Desc

			for _, subAccount := range account.Subs {
				for _, unit := range subAccount.Units {
					aboveValueLimits = append(aboveValueLimits, unitToAboveValueLimit(unit.Id, category, subCategory))
				}
			}
		}
	}

	return aboveValueLimits
}

func unitToAboveValueLimit(id string, category string, subCategory string) *model.AboveValueLimit {
	return &model.AboveValueLimit{
		ID:          id,
		Category:    category,
		SubCategory: subCategory,
	}
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
