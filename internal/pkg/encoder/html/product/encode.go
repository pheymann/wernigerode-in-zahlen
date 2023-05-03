package product

import (
	"fmt"
	"html/template"
	"sort"

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

func balanceToSection(balance model.AccountBalance, year model.BudgetYear, p *message.Printer) html.BalanceSection {
	accountsSplit := splitAccountsByType(balance.Accounts, year)
	accountsIncome := accountsSplit.First
	accountsExpenses := accountsSplit.Second

	hasIncome := len(accountsIncome) > 0
	hasExpenses := len(accountsExpenses) > 0

	chartIDUniq := balance.Type

	section := html.BalanceSection{
		ID: "balance_" + chartIDUniq,

		HasIncomeAndExpenses: hasIncome && hasExpenses,

		HasIncome:           hasIncome,
		IncomeCashflowTotal: balance.Cashflow.Income[year],
		Income:              dataPointsToChartJSDataset(accountsIncome, year, chartIDUniq+"_income"),

		HasExpenses:           hasExpenses,
		ExpensesCashflowTotal: balance.Cashflow.Expenses[year],
		Expenses:              dataPointsToChartJSDataset(accountsExpenses, year, chartIDUniq+"_expenses"),

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

func splitAccountsByType(accounts []model.Account, year model.BudgetYear) shared.Pair[[]model.Account, []model.Account] {
	var income = make([]model.Account, 0)
	var expenses = make([]model.Account, 0)

	for _, account := range accounts {
		if !shared.IsUnequal(account.Budget[year], 0.0) {
			continue
		}

		if account.Type == model.AccountTypeIncome {
			income = append(income, account)
		} else if account.Type == model.AccountTypeExpense {
			expenses = append(expenses, account)
		} else {
			panic("Unknown account type: " + account.Type)
		}
	}

	return shared.NewPair(income, expenses)
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

func dataPointsToChartJSDataset(accounts []model.Account, year model.BudgetYear, chartIDUniq string) html.ChartJSDataset {
	var labels = []string{}
	var data = []float64{}

	for _, account := range accounts {
		if shared.IsUnequal(account.Budget[year], 0) {
			labels = append(labels, account.Description)
			data = append(data, account.Budget[year])
		}
	}

	return html.ChartJSDataset{
		ID:           "chartjs_" + chartIDUniq,
		Labels:       labels,
		DatasetLabel: "Budget",
		Data:         data,
	}
}

func encodeBalanceSectionHeader(balance model.AccountBalance, year model.BudgetYear, p *message.Printer) template.HTML {
	return template.HTML(fmt.Sprintf(
		`- %s <span class="%s">%s</span>`,
		encodeAccountClass(balance.Type, balance.Cashflow.Total[year]),
		encodeHtml.EncodeCSSCashflowClass(balance.Cashflow.Total[year]),
		encodeHtml.EncodeBudget(balance.Cashflow.Total[year], p),
	))
}

func encodeAccountClass(tpe model.AccountBalanceType, cashflowTotal float64) string {
	switch tpe {
	case model.AccountBalanceTypeAdministration:
		if cashflowTotal >= 0 {
			return "Die Verwaltung erwirtschaftet"
		}
		return "Die Verwaltung kostet"

	case model.AccountBalanceTypeInvestments:
		if cashflowTotal >= 0 {
			return "Investitionen erwirtschaften"
		}
		return "Investitionen kosten"

	default:
		panic(fmt.Sprintf("unknown account class '%s'", tpe))
	}
}
