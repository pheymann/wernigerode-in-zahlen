package financialplan_a

import (
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

func Decode(rows []model.RawCSVRow) model.FinancialPlanA {
	financialPlanA := model.FinancialPlanA{}

	for _, row := range rows {
		id := decoder.DecodeString(row.Regexp, "id", row.Matches)

		if row.Tpe == model.RowTypeUnitAccount {
			if class, ok := accountBalanceClassifier[id]; ok {
				financialPlanA.Balances = append(financialPlanA.Balances, decodeAccountBalance(row, id, class))
			} else if _, ok := accountClassifier[id]; ok {
				lastBalanceIndex := len(financialPlanA.Balances) - 1

				financialPlanA.Balances[lastBalanceIndex].Accounts = append(
					financialPlanA.Balances[lastBalanceIndex].Accounts,
					decodeAccount(row, id),
				)
			} else {
				// sub-account

			}
		} else {
			// sub-account unit
		}
	}

	return financialPlanA
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
