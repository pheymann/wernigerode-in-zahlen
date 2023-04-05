package financialplan

import (
	"regexp"

	"wernigerode-in-zahlen.de/internal/pkg/decoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	fd "wernigerode-in-zahlen.de/internal/pkg/model/financialdata"
)

var (
	isAdminIncomeRegex  = regexp.MustCompile(`^(\d\.)+(\d{2}\.)+[^6]\d+$`)
	isAdminExpenseRegex = regexp.MustCompile(`^(\d\.)+(\d{2}\.)+[^7]\d+$`)

	isInvestmentIncomeRegex  = regexp.MustCompile(`^(\d\.)+(\d{2})+\/\d{4}\.6\d+$`)
	isInvestmentExpenseRegex = regexp.MustCompile(`^(\d\.)+(\d{2})+\/\d{4}\.7\d+$`)

	idRegex = regexp.MustCompile(`^(?P<id>\d\.\d\.\d\.\d{2}(\.\d{2})?).*$`)
)

func DecodeFromAccounts(accounts []fd.Account) model.FinancialPlanProduct {
	var setMetadata = false
	var financialPlan = model.NewFinancialPlanProduct()

	for _, account := range accounts {
		if !setMetadata {
			addID(financialPlan, account)
			setMetadata = true
		}

		if isAdminIncomeRegex.MatchString(account.ID) {
			updateAdministrationBalance(financialPlan, account, false)
		} else if isAdminExpenseRegex.MatchString(account.ID) {
			updateAdministrationBalance(financialPlan, account, true)
		} else if isInvestmentIncomeRegex.MatchString(account.ID) {
			updateInvestmentsBalance(financialPlan, account, false)
		} else if isInvestmentExpenseRegex.MatchString(account.ID) {
			updateInvestmentsBalance(financialPlan, account, true)
		} else {
			panic("Unknown account type: " + account.ID)
		}
	}

	return *financialPlan
}

func addID(plan *model.FinancialPlanProduct, someAccount fd.Account) {
	matches := idRegex.FindStringSubmatch(someAccount.ID)

	plan.ID = decoder.DecodeString(idRegex, "id", matches)
}

func updateAdministrationBalance(plan *model.FinancialPlanProduct, account fd.Account, isExpense bool) {
	forBudget(plan, account.Budget, func(year model.BudgetYear, value float64) {
		updateCashflow(plan, isExpense, year, signBudget(value, isExpense))
	})

	plan.AdministrationBalance.Accounts = append(plan.AdministrationBalance.Accounts, model.Account2{
		ID:          account.ID,
		ProductID:   account.ProductID,
		Description: account.Description,
		Budget:      account.Budget,
	})
}

func updateInvestmentsBalance(plan *model.FinancialPlanProduct, account fd.Account, isExpense bool) {
	forBudget(plan, account.Budget, func(year model.BudgetYear, value float64) {
		updateCashflow(plan, isExpense, year, signBudget(value, isExpense))
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

func updateCashflow(plan *model.FinancialPlanProduct, isExpense bool, year model.BudgetYear, value float64) {
	plan.CashFlow.Total[year] += value

	if isExpense {
		plan.CashFlow.Expenses[year] += value
	} else {
		plan.CashFlow.Income[year] += value
	}
}

func forBudget(plan *model.FinancialPlanProduct, budgets map[model.BudgetYear]float64, update func(model.BudgetYear, float64)) {
	for year, budget := range budgets {
		update(year, budget)
	}
}
