package product

import (
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/text/message"
	encodeHtml "wernigerode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	plan model.FinancialPlanProduct,
	tableData []html.AccountTableData,
	year model.BudgetYear,
	p *message.Printer,
) html.Product {
	var sections = balanceDataToSections(plan, year, p)
	subProductSection := subProductsToSection(plan, year, p)
	if subProductSection != nil {
		sections = append(sections, *subProductSection)
	}

	return html.Product{
		Meta:            plan.Metadata,
		BalanceSections: sections,
		Copy: html.ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",
			Year:     year,

			IntroCashflowTotal: fmt.Sprintf("Das Produkt - %s - wird in %s", plan.Metadata.Description, year),
			IntroDescription:   encodeIntroDescription(plan.Cashflow.Total[year], plan.Metadata),

			CashflowTotal: encodeHtml.EncodeBudget(plan.Cashflow.Total[year], p),

			Accounts: encodeAccountCopy(tableData, p),

			MetaDepartment:    "Fachbereich",
			MetaProductClass:  "Produktklasse",
			MetaProductDomain: "Produktbereich",
			MetaProductGroup:  "Produktgruppe",
			MetaProduct:       "Produkt",
			MetaAccountable:   "Verantwortlich",
			MetaResponsible:   "Zuständig",
			MetaMission:       "Aufgabe",
			MetaTargets:       "Zielgruppe",
			MetaServices:      "Dienstleistungen",
			MetaGrouping:      "Klassifizierung",

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>`,
		},
		CSS: html.ProductCSS{
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(plan.Cashflow.Total[year]),
		},
	}
}

func encodeIntroDescription(cashflowTotal float64, meta model.Metadata) string {
	var expenseEarnCopy = "kosten"
	if cashflowTotal >= 0 {
		expenseEarnCopy = "einbringen"
	}

	return fmt.Sprintf(
		`%s. Unten werden die verschiedenen Ausgaben aufgeführt, die das Budget ausmachen. Jede Ausgabe hat dann nochmal eine
		Auflistung von Kostenstellen. Ganz am Ende dieser Seite findest du noch eine Beschreibung des Produkts.`,
		expenseEarnCopy,
	)
}

func balanceDataToSections(plan model.FinancialPlanProduct, year model.BudgetYear, p *message.Printer) []html.BalanceSection {
	return []html.BalanceSection{
		balanceToSection(plan.AdministrationBalance, year, p),
		balanceToSection(plan.InvestmentsBalance, year, p),
	}
}

func balanceToSection(balance model.AccountBalance2, year model.BudgetYear, p *message.Printer) html.BalanceSection {
	adminAccountsSplit := splitAccountsByType(balance.Accounts)
	adminAccountsIncome := adminAccountsSplit.First
	adminAccountsExpenses := adminAccountsSplit.Second

	hasIncome := shared.IsUnequal(balance.Cashflow.Income[year], 0.0)
	hasExpenses := shared.IsUnequal(balance.Cashflow.Expenses[year], 0.0)

	section := html.BalanceSection{
		ID: strings.ReplaceAll("balance-"+uuid.New().String(), "-", ""),

		HasIncomeAndExpenses: hasIncome && hasExpenses,

		HasIncome:            hasIncome,
		HasMoreThanOneIncome: len(adminAccountsIncome) > 1,
		IncomeCashflowTotal:  balance.Cashflow.Income[year],
		Income:               dataPointsToChartJSDataset(adminAccountsIncome, year),

		HasExpenses:           hasExpenses,
		HasMoreThanOneExpense: len(adminAccountsExpenses) > 1,
		ExpensesCashflowTotal: balance.Cashflow.Expenses[year],
		Expenses:              dataPointsToChartJSDataset(adminAccountsExpenses, year),

		Copy: html.BalanceSectionCopy{
			Header:                encodeBalanceSectionHeader(balance, year, p),
			IncomeCashflowTotal:   "Einnahmen: " + encodeHtml.EncodeBudget(balance.Cashflow.Income[year], p),
			ExpensesCashflowTotal: "Ausgaben: " + encodeHtml.EncodeBudget(balance.Cashflow.Expenses[year], p),
		},
		CSS: html.BalanceSectionCSS{
			CashflowTotalClass: encodeHtml.EncodeCSSCashflowClass(balance.Cashflow.Total[year]),
		},
	}

	return section
}

func splitAccountsByType(accounts []model.Account2) shared.Pair[[]model.Account2, []model.Account2] {
	var income = make([]model.Account2, 0)
	var expenses = make([]model.Account2, 0)

	for _, account := range accounts {
		if account.Type == model.Account2TypeIncome {
			income = append(income, account)
		} else if account.Type == model.Account2TypeExpense {
			expenses = append(expenses, account)
		}
	}

	return shared.NewPair(income, expenses)
}

func subProductsToSection(plan model.FinancialPlanProduct, year model.BudgetYear, p *message.Printer) *html.BalanceSection {
	if len(plan.SubProducts) == 0 {
		return nil
	}

	section := html.BalanceSection{
		ID: strings.ReplaceAll("balance-sub-product-"+uuid.New().String(), "-", "_"),
		Income: html.ChartJSDataset{
			ID:           strings.ReplaceAll("chartjs-"+uuid.New().String(), "-", "_"),
			DatasetLabel: "Einnahmen",
		},
		Expenses: html.ChartJSDataset{
			ID:           strings.ReplaceAll("chartjs-"+uuid.New().String(), "-", "_"),
			DatasetLabel: "Ausgaben",
		},
	}

	section.IncomeID = template.JS(section.Income.ID)
	section.ExpensesID = template.JS(section.Expenses.ID)

	for _, plan := range plan.SubProducts {
		link := fmt.Sprintf("%s/product.html", plan.Metadata.SubProduct.ID)

		if plan.Cashflow.Total[year] < 0 {
			section.Expenses.Labels = append(section.Expenses.Labels, plan.Metadata.SubProduct.Name)
			section.Expenses.Data = append(section.Expenses.Data, plan.Cashflow.Total[year])
			section.ExpensesSubProductLinks = append(section.ExpensesSubProductLinks, link)
		} else {
			section.Income.Labels = append(section.Income.Labels, plan.Metadata.SubProduct.Name)
			section.Income.Data = append(section.Income.Data, plan.Cashflow.Total[year])
			section.IncomeSubProductLinks = append(section.IncomeSubProductLinks, link)
		}
	}

	if shared.IsUnequal(plan.Cashflow.Expenses[year], 0.0) {
		section.HasExpenses = true
		section.HasMoreThanOneExpense = len(section.ExpensesSubProductLinks) > 1
		section.ExpensesCashflowTotal = plan.Cashflow.Expenses[year]
		section.HasExpensesSubProductLinks = true
	}
	if shared.IsUnequal(plan.Cashflow.Income[year], 0.0) {
		section.HasIncome = true
		section.HasMoreThanOneIncome = len(section.IncomeSubProductLinks) > 1
		section.IncomeCashflowTotal = plan.Cashflow.Income[year]
		section.HasIncomeSubProductLinks = true
	}

	if section.HasIncome && section.HasExpenses {
		section.HasIncomeAndExpenses = true
	}

	section.Copy = html.BalanceSectionCopy{
		Header:                encodeSubProductBalanceSectionHeader(plan.Cashflow.Total[year], p),
		IncomeCashflowTotal:   "Einnahmen: " + encodeHtml.EncodeBudget(plan.Cashflow.Income[year], p),
		ExpensesCashflowTotal: "Ausgaben: " + encodeHtml.EncodeBudget(plan.Cashflow.Expenses[year], p),
	}

	section.CSS = html.BalanceSectionCSS{
		CashflowTotalClass: encodeHtml.EncodeCSSCashflowClass(plan.Cashflow.Total[year]),
	}

	return &section
}

func encodeAccountCopy(data []html.AccountTableData, p *message.Printer) []html.ProductAccountCopy {
	var accountCopy = []html.ProductAccountCopy{}

	for _, row := range data {
		accountCopy = append(accountCopy, html.ProductAccountCopy{
			Name:          row.Name,
			CashflowTotal: encodeHtml.EncodeBudget(row.CashflowTotal, p),
		})
	}

	sort.Slice(accountCopy, func(i, j int) bool {
		return accountCopy[i].Name < accountCopy[j].Name
	})
	return accountCopy
}

func dataPointsToChartJSDataset(accounts []model.Account2, year model.BudgetYear) html.ChartJSDataset {
	var labels = []string{}
	var data = []float64{}

	for _, account := range accounts {
		if shared.IsUnequal(account.Budget[year], 0) {
			labels = append(labels, account.Description)
			data = append(data, account.Budget[year])
		}
	}

	return html.ChartJSDataset{
		ID:           strings.ReplaceAll("chartjs-"+uuid.New().String(), "-", "_"),
		Labels:       labels,
		DatasetLabel: "Budget",
		Data:         data,
	}
}

func encodeBalanceSectionHeader(balance model.AccountBalance2, year model.BudgetYear, p *message.Printer) template.HTML {
	return template.HTML(fmt.Sprintf(
		`- %s <span class="%s">%s</span>`,
		encodeAccountClass(balance.Type, balance.Cashflow.Total[year]),
		encodeHtml.EncodeCSSCashflowClass(balance.Cashflow.Total[year]),
		encodeHtml.EncodeBudget(balance.Cashflow.Total[year], p),
	))
}

func encodeSubProductBalanceSectionHeader(cashflowTotal float64, p *message.Printer) template.HTML {
	return template.HTML(fmt.Sprintf(
		`- %s <span class="%s">%s</span>`,
		"Darin enthalten sind die folgenden Unter-Produkte: ",
		encodeHtml.EncodeCSSCashflowClass(cashflowTotal),
		encodeHtml.EncodeBudget(cashflowTotal, p),
	))
}

func encodeAccountClass(tpe model.AccountBalance2Type, cashflowTotal float64) string {
	switch tpe {
	case model.AccountBalance2TypeAdministration:
		if cashflowTotal >= 0 {
			return "Die Verwaltung erwirtschaftet"
		}
		return "Die Verwaltung kostet"

	case model.AccountBalance2TypeInvestments:
		if cashflowTotal >= 0 {
			return "Investitionen erwirtschaften"
		}
		return "Investitionen kosten"

	default:
		panic(fmt.Sprintf("unknown account class '%s'", tpe))
	}
}
