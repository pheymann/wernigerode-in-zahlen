package financialplan_b

import (
	"fmt"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	accountBalanceClassifier = map[string]model.AccountClass{
		"60": model.AccountClassOneOff,
	}
	accountClassifier = map[string]model.AccountClass{
		"10": model.AccountClassOneOff,
		"20": model.AccountClassOneOff,
		"40": model.AccountClassOneOff,
		"50": model.AccountClassOneOff,
	}
)

type rowTpeState = string

const (
	rowTpeStateBalance     rowTpeState = "balance"
	rowTpeStateAccount     rowTpeState = "account"
	rowTpeStateSubAccount  rowTpeState = "sub"
	rowTpeStateUnitAccount rowTpeState = "unit"
)

func DecodeFromCSV(rows []model.RawCSVRow) model.FinancialPlan {
	financialPlanB := &model.FinancialPlan{}

	var lastTpe rowTpeState = ""

	for _, row := range rows {
		if row.Tpe == model.RowTypeSeparateLine {
			switch lastTpe {
			case rowTpeStateBalance:
				financialPlanB.UpdateLastAccountBalance(func(balance model.AccountBalance) model.AccountBalance {
					balance.Desc = updateDesc(balance.Desc, row.Regexp, row.Matches)

					return balance
				})

			case rowTpeStateAccount:
				financialPlanB.UpdateLastAccount(func(account model.Account) model.Account {
					account.Desc = updateDesc(account.Desc, row.Regexp, row.Matches)

					return account
				})

			case rowTpeStateSubAccount:
				financialPlanB.UpdateLastSubAccount(func(subAccount model.SubAccount) model.SubAccount {
					subAccount.Desc = updateDesc(subAccount.Desc, row.Regexp, row.Matches)

					return subAccount
				})

			case rowTpeStateUnitAccount:
				financialPlanB.UpdateLastUnitAccount(func(unitAccount model.UnitAccount) model.UnitAccount {
					unitAccount.Desc = updateDesc(unitAccount.Desc, row.Regexp, row.Matches)

					return unitAccount
				})

			default:
				panic(fmt.Sprintf("unknown RawCSVRow type '%s'", lastTpe))
			}
		} else if row.Tpe == model.RowTypeIgnore {
			continue
		} else {
			id := decoder.DecodeString(row.Regexp, "id", row.Matches)

			if row.Tpe == model.RowTypeOneOff {
				financialPlanB.AddAccountBalance(decodeAccountBalance(row, id, model.AccountClassOneOff, model.AccountBalance{}))
			} else if row.Tpe == model.RowTypeOther {
				if class, ok := accountBalanceClassifier[id]; ok {
					lastTpe = rowTpeStateBalance

					financialPlanB.UpdateLastAccountBalance(func(balance model.AccountBalance) model.AccountBalance {
						balance.Budgets = map[model.BudgetYear]float64{
							model.BudgetYear2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
							model.BudgetYear2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
							model.BudgetYear2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
							model.BudgetYear2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
							model.BudgetYear2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
							model.BudgetYear2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
						}
						balance.Class = class

						return balance
					})
				} else if _, ok := accountClassifier[id]; ok {
					lastTpe = rowTpeStateAccount

					financialPlanB.AddAccount(decodeAccount(row, id, model.Account{}))
				} else {
					lastTpe = rowTpeStateSubAccount

					financialPlanB.AddSubAccount(decodeSubAccount(row, id))
				}
			} else {
				// UnitAccount
				lastTpe = rowTpeStateUnitAccount

				financialPlanB.AddUnitAccount(decodeUnitAccount(row, id))
			}
		}
	}

	financialPlanB.RemoveLastAccountBalance()
	return *financialPlanB
}

func updateDesc(original string, regex *regexp.Regexp, matches []string) string {
	return fmt.Sprintf("%s %s", original, decoder.DecodeString(regex, "desc", matches))
}

func decodeAccountBalance(row model.RawCSVRow, id string, class model.AccountClass, balance model.AccountBalance) model.AccountBalance {
	return model.AccountBalance{
		Id:       id,
		Class:    class,
		Desc:     decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Accounts: balance.Accounts,
	}
}

func decodeAccount(row model.RawCSVRow, id string, account model.Account) model.Account {
	return model.Account{
		Id:   id,
		Desc: decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budgets: map[model.BudgetYear]float64{
			model.BudgetYear2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
			model.BudgetYear2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
			model.BudgetYear2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
			model.BudgetYear2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
			model.BudgetYear2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
			model.BudgetYear2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
		},
		Subs: account.Subs,
	}
}

func decodeSubAccount(row model.RawCSVRow, id string) model.SubAccount {
	return model.SubAccount{
		Id:   id,
		Desc: decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budgets: map[model.BudgetYear]float64{
			model.BudgetYear2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
			model.BudgetYear2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
			model.BudgetYear2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
			model.BudgetYear2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
			model.BudgetYear2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
			model.BudgetYear2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
		},
	}
}

func decodeUnitAccount(row model.RawCSVRow, id string) model.UnitAccount {
	return model.UnitAccount{
		Id:   id,
		Desc: decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budgets: map[model.BudgetYear]float64{
			model.BudgetYear2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
			model.BudgetYear2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
			model.BudgetYear2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
			model.BudgetYear2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
			model.BudgetYear2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
			model.BudgetYear2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
		},
	}
}
