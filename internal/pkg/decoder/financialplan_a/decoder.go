package financialplan_a

import (
	"fmt"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	accountBalanceClassifier = map[string]model.AccountClass{
		"16": model.AccountClassAdministration,
		"34": model.AccountClassInvestments,
	}
	accountClassifier = map[string]model.AccountClass{
		"8":  model.AccountClassAdministration,
		"15": model.AccountClassAdministration,
		"17": model.AccountClassInvestments,
		"25": model.AccountClassInvestments,
		"26": model.AccountClassInvestments,
		"33": model.AccountClassInvestments,
	}
)

type rowTpeState = string

const (
	rowTpeStateBalance     rowTpeState = "balance"
	rowTpeStateAccount     rowTpeState = "account"
	rowTpeStateSubAccount  rowTpeState = "sub"
	rowTpeStateUnitAccount rowTpeState = "unit"
)

func Decode(rows []model.RawCSVRow) model.FinancialPlanA {
	financialPlanA := &model.FinancialPlanA{}

	var lastTpe rowTpeState = ""

	for _, row := range rows {
		id := decoder.DecodeString(row.Regexp, "id", row.Matches)

		if row.Tpe == model.RowTypeOther {
			if class, ok := accountBalanceClassifier[id]; ok {
				lastTpe = rowTpeStateBalance

				financialPlanA.UpdateLastAccountBalance(func(_ model.AccountBalance) model.AccountBalance {
					return decodeAccountBalance(row, id, class)
				})
			} else if _, ok := accountClassifier[id]; ok {
				lastTpe = rowTpeStateAccount

				financialPlanA.UpdateLastAccount(func(_ model.Account) model.Account {
					return decodeAccount(row, id)
				})
			} else {
				lastTpe = rowTpeStateSubAccount

				financialPlanA.UpdateLastSubAccount(func(_ model.SubAccount) model.SubAccount {
					return decodeSubAccount(row, id)
				})
			}
		} else if row.Tpe == model.RowTypeUnitAccount {
			lastTpe = rowTpeStateUnitAccount

			financialPlanA.UpdateLastUnitAccount(func(_ model.UnitAccount) model.UnitAccount {
				return decodeUnitAccount(row, id)
			})
		} else {
			// separate line
			switch lastTpe {
			case rowTpeStateBalance:
				financialPlanA.UpdateLastAccountBalance(func(balance model.AccountBalance) model.AccountBalance {
					balance.Desc = updateDesc(balance.Desc, row.Regexp, row.Matches)

					return balance
				})

			case rowTpeStateAccount:
				financialPlanA.UpdateLastAccount(func(account model.Account) model.Account {
					account.Desc = updateDesc(account.Desc, row.Regexp, row.Matches)

					return account
				})

			case rowTpeStateSubAccount:
				financialPlanA.UpdateLastSubAccount(func(subAccount model.SubAccount) model.SubAccount {
					subAccount.Desc = updateDesc(subAccount.Desc, row.Regexp, row.Matches)

					return subAccount
				})

			case rowTpeStateUnitAccount:
				financialPlanA.UpdateLastUnitAccount(func(unitAccount model.UnitAccount) model.UnitAccount {
					unitAccount.Desc = updateDesc(unitAccount.Desc, row.Regexp, row.Matches)

					return unitAccount
				})

			default:
				panic(fmt.Sprintf("unknown RawCSVRow type '%s'", lastTpe))
			}
		}
	}

	return *financialPlanA
}

func updateDesc(original string, regex *regexp.Regexp, matches []string) string {
	return fmt.Sprintf("%s %s", original, decoder.DecodeString(regex, "desc", matches))
}

func decodeAccountBalance(row model.RawCSVRow, id string, class model.AccountClass) model.AccountBalance {
	return model.AccountBalance{
		Id:         id,
		Class:      class,
		Desc:       decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budget2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
		Budget2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
		Budget2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
		Budget2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
		Budget2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
		Budget2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
	}
}

func decodeAccount(row model.RawCSVRow, id string) model.Account {
	return model.Account{
		Id:         id,
		Desc:       decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budget2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
		Budget2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
		Budget2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
		Budget2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
		Budget2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
		Budget2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
	}
}

func decodeSubAccount(row model.RawCSVRow, id string) model.SubAccount {
	return model.SubAccount{
		Id:         id,
		Desc:       decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budget2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
		Budget2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
		Budget2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
		Budget2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
		Budget2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
		Budget2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
	}
}

func decodeUnitAccount(row model.RawCSVRow, id string) model.UnitAccount {
	return model.UnitAccount{
		Id:         id,
		Desc:       decoder.DecodeString(row.Regexp, "desc", row.Matches),
		Budget2020: decoder.DecodeBudget(row.Regexp, "_2020", row.Matches),
		Budget2021: decoder.DecodeBudget(row.Regexp, "_2021", row.Matches),
		Budget2022: decoder.DecodeBudget(row.Regexp, "_2022", row.Matches),
		Budget2023: decoder.DecodeBudget(row.Regexp, "_2023", row.Matches),
		Budget2024: decoder.DecodeBudget(row.Regexp, "_2024", row.Matches),
		Budget2025: decoder.DecodeBudget(row.Regexp, "_2025", row.Matches),
	}
}
