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

func DecodeFromAccounts(accounts []fd.Account) model.FinancialPlan2 {
	var setMetadata = false
	var financialPlan = &model.FinancialPlan2{}

	for _, account := range accounts {
		if !setMetadata {
			addMetadata(financialPlan, account)
			setMetadata = true
		}

		if isAdminIncomeRegex.MatchString(account.ID) {
			updateAdministrationBalance(financialPlan, account)
		} else if isAdminExpenseRegex.MatchString(account.ID) {
			account.Budget = negateBudget(account.Budget)
			updateAdministrationBalance(financialPlan, account)
		} else if isInvestmentIncomeRegex.MatchString(account.ID) {
			updateInvestmentsBalance(financialPlan, account)
		} else if isInvestmentExpenseRegex.MatchString(account.ID) {
			account.Budget = negateBudget(account.Budget)
			updateInvestmentsBalance(financialPlan, account)
		} else {
			panic("Unknown account type: " + account.ID)
		}
	}

	return *financialPlan
}

func addMetadata(plan *model.FinancialPlan2, someAccount fd.Account) {
	matches := metadataRegex.FindStringSubmatch(someAccount.ID)

	plan.ProductClass = decoder.DecodeString(metadataRegex, "class", matches)
	plan.ProductDomain = decoder.DecodeString(metadataRegex, "domain", matches)
	plan.ProductGroup = decoder.DecodeString(metadataRegex, "group", matches)
	plan.Product = decoder.DecodeString(metadataRegex, "product", matches)
	plan.SubProduct = decoder.DecodeOptString(metadataRegex, "sub_product", matches)
}

func updateAdministrationBalance(plan *model.FinancialPlan2, account fd.Account) {
	for budgetYear, value := range account.Budget {
		plan.AdministrationBalance.Budget[budgetYear] += value
	}

	plan.AdministrationBalance.Accounts = append(plan.AdministrationBalance.Accounts, model.Account2{
		ID:          account.ID,
		ProductID:   account.ProductID,
		Description: account.Description,
		Budget:      account.Budget,
	})
}

func updateInvestmentsBalance(plan *model.FinancialPlan2, account fd.Account) {
	for budgetYear, value := range account.Budget {
		plan.InvestmentsBalance.Budget[budgetYear] += value
	}

	plan.InvestmentsBalance.Accounts = append(plan.InvestmentsBalance.Accounts, model.Account2{
		ID:          account.ID,
		ProductID:   account.ProductID,
		Description: account.Description,
		Budget:      account.Budget,
	})
}

func negateBudget(budget map[string]float64) map[string]float64 {
	for budgetYear, value := range budget {
		budget[budgetYear] = -value
	}

	return budget
}
