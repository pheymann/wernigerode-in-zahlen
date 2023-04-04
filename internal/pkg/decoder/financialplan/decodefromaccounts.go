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

	metadataRegex = regexp.MustCompile(`^(?P<class>\d)\.(?P<domain>\d)\.(?P<group>\d)\.(?P<product>\d{2})(\.(?P<sub_product>\d{2}))?.*$`)
)

func DecodeFromAccounts(accounts []fd.Account) model.FinancialPlanProduct {
	var setMetadata = false
	var financialPlan = &model.FinancialPlanProduct{}

	for _, account := range accounts {
		if !setMetadata {
			addMetadata(financialPlan, account)
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

func addMetadata(plan *model.FinancialPlanProduct, someAccount fd.Account) {
	matches := metadataRegex.FindStringSubmatch(someAccount.ID)

	plan.ProductClass = decoder.DecodeString(metadataRegex, "class", matches)
	plan.ProductDomain = decoder.DecodeString(metadataRegex, "domain", matches)
	plan.ProductGroup = decoder.DecodeString(metadataRegex, "group", matches)
	plan.Product = decoder.DecodeString(metadataRegex, "product", matches)
	plan.SubProduct = decoder.DecodeOptString(metadataRegex, "sub_product", matches)
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
