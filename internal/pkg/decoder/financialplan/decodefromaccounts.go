package financialplan

import (
	"regexp"

	"wernigerode-in-zahlen.de/internal/pkg/decoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	fd "wernigerode-in-zahlen.de/internal/pkg/model/financialdata"
)

var (
	adminAccountIdRegex      = regexp.MustCompile(`^(\d\.)+(\d{2}\.)+(?P<id>\d+)$`)
	investmentAccountIdRegex = regexp.MustCompile(`^(\d\.)+(\d{2}\.?)+\/(?P<investment_class>\d{4})\.(?P<id>\d+)$`)

	idRegex = regexp.MustCompile(`^(?P<id>\d\.\d\.\d\.\d{2}(\.\d{2})?)[^\d].+$`)
)

func DecodeFromAccounts(accounts []fd.Account) model.FinancialPlanProduct {
	var setMetadata = false
	var financialPlan = model.NewFinancialPlanProduct()

	for _, account := range accounts {
		var isAdminAccount = false
		var matches = investmentAccountIdRegex.FindStringSubmatch(account.ID)
		if matches == nil {
			isAdminAccount = true
			matches = adminAccountIdRegex.FindStringSubmatch(account.ID)
			if matches == nil {
				panic("could not find id in account id: " + account.ID)
			}
		}

		var id = ""
		var investmentClass = ""
		if !isAdminAccount {
			id = decoder.DecodeString(investmentAccountIdRegex, "id", matches)
			investmentClass = decoder.DecodeString(investmentAccountIdRegex, "investment_class", matches)
		} else {
			id = decoder.DecodeString(adminAccountIdRegex, "id", matches)
		}

		if !setMetadata {
			addID(financialPlan, account)
			setMetadata = true
		}

		switch id[0] {
		case '4':
			// ignore result accounts
			continue

		case '5':
			// ignore result accounts
			continue

		case '6':
			if isAdminAccount {
				updateAdministrationBalance(financialPlan, account, false)
			} else {
				if !isFinancialActivity(investmentClass) {
					updateInvestmentsBalance(financialPlan, account, false)
				}
			}

		case '7':
			if isAdminAccount {
				updateAdministrationBalance(financialPlan, account, true)
			} else {
				if !isFinancialActivity(investmentClass) {
					updateInvestmentsBalance(financialPlan, account, true)
				}
			}

		case '8':
			// ignore correction accounts
			continue

		default:
			panic("unknown id: " + id)
		}
	}

	return *financialPlan
}

func addID(plan *model.FinancialPlanProduct, someAccount fd.Account) {
	matches := idRegex.FindStringSubmatch(someAccount.ID)

	plan.ID = decoder.DecodeString(idRegex, "id", matches)
}

func isFinancialActivity(investmentClass string) bool {
	return investmentClass[0] == '9'
}

func updateAdministrationBalance(plan *model.FinancialPlanProduct, account fd.Account, isExpense bool) {
	forBudget(plan, account.Budget, func(year model.BudgetYear, value float64) {
		updateCashflow(plan, isExpense, true, year, signBudget(value, isExpense))
	})

	tpe := model.Account2TypeIncome
	if isExpense {
		tpe = model.Account2TypeExpense
	}

	plan.AdministrationBalance.Accounts = append(plan.AdministrationBalance.Accounts, model.Account2{
		ID:          account.ID,
		ProductID:   account.ProductID,
		Description: account.Description,
		Type:        tpe,
		Budget:      account.Budget,
	})
}

func updateInvestmentsBalance(plan *model.FinancialPlanProduct, account fd.Account, isExpense bool) {
	forBudget(plan, account.Budget, func(year model.BudgetYear, value float64) {
		updateCashflow(plan, isExpense, false, year, signBudget(value, isExpense))
	})

	plan.InvestmentsBalance.Accounts = append(plan.InvestmentsBalance.Accounts, model.Account2{
		ID:          account.ID,
		ProductID:   account.ProductID,
		Description: account.Description,
		Budget:      account.Budget,
	})
}

func signBudget(value float64, isExpense bool) float64 {
	if isExpense {
		return -value
	}

	return value
}

func updateCashflow(plan *model.FinancialPlanProduct, isExpense bool, isAdmin bool, year model.BudgetYear, value float64) {
	plan.Cashflow.Total[year] += value

	if isAdmin {
		plan.AdministrationBalance.Cashflow.Total[year] += value
	} else {
		plan.InvestmentsBalance.Cashflow.Total[year] += value
	}

	if isExpense {
		plan.Cashflow.Expenses[year] += value

		if isAdmin {
			plan.AdministrationBalance.Cashflow.Expenses[year] += value
		} else {
			plan.InvestmentsBalance.Cashflow.Expenses[year] += value
		}
	} else {
		plan.Cashflow.Income[year] += value

		if isAdmin {
			plan.AdministrationBalance.Cashflow.Income[year] += value
		} else {
			plan.InvestmentsBalance.Cashflow.Income[year] += value
		}
	}
}

func forBudget(plan *model.FinancialPlanProduct, budgets map[model.BudgetYear]float64, update func(model.BudgetYear, float64)) {
	for year, budget := range budgets {
		update(year, budget)
	}
}
