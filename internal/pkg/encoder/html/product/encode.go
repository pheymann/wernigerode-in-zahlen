package product

import (
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/text/message"
	encodeHtml "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	metadata model.Metadata,
	fpBalanceData []html.BalanceData,
	fpCashflowTotal float64,
	tableData []html.AccountTableData,
	subProductData []html.ProductData,
	year model.BudgetYear,
	p *message.Printer,
) html.Product {
	var sections = balanceDataToSections(fpBalanceData, year, p)
	subProductSection := subProductsToSection(subProductData, year, p)
	if subProductSection != nil {
		sections = append(sections, *subProductSection)
	}

	return html.Product{
		Meta:            metadata,
		BalanceSections: sections,
		Copy: html.ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",
			Year:     year,

			IntroCashflowTotal: fmt.Sprintf("Das Produkt - %s - wird in %s", metadata.Description, year),
			IntroDescription:   encodeIntroDescription(fpCashflowTotal, metadata),

			CashflowTotal: encodeHtml.EncodeBudget(fpCashflowTotal, p),

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
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(fpCashflowTotal),
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

func balanceDataToSections(data []html.BalanceData, year model.BudgetYear, p *message.Printer) []html.BalanceSection {
	var sections = []html.BalanceSection{}
	for _, balance := range data {
		var incomeCashflowTotal float64
		for _, income := range balance.Income {
			incomeCashflowTotal += income.Budget
		}
		var expensesCashflowTotal float64
		for _, expense := range balance.Expenses {
			expensesCashflowTotal += expense.Budget
		}

		sections = append(sections, html.BalanceSection{
			ID: strings.ReplaceAll("balance-"+uuid.New().String(), "-", ""),

			HasIncomeAndExpenses: len(balance.Income) > 0 && len(balance.Expenses) > 0,

			HasIncome:            len(balance.Income) > 0,
			HasMoreThanOneIncome: len(balance.Income) > 1,
			IncomeCashflowTotal:  incomeCashflowTotal,
			Income:               dataPointsToChartJSDataset(balance.Income),

			HasExpenses:           len(balance.Expenses) > 0,
			HasMoreThanOneExpense: len(balance.Expenses) > 1,
			ExpensesCashflowTotal: expensesCashflowTotal,
			Expenses:              dataPointsToChartJSDataset(balance.Expenses),

			Copy: html.BalanceSectionCopy{
				Header:                encodeBalanceSectionHeader(balance.Balance, year, p),
				IncomeCashflowTotal:   "Einnahmen: " + encodeHtml.EncodeBudget(incomeCashflowTotal, p),
				ExpensesCashflowTotal: "Ausgaben: " + encodeHtml.EncodeBudget(expensesCashflowTotal, p),
			},
			CSS: html.BalanceSectionCSS{
				CashflowTotalClass: encodeHtml.EncodeCSSCashflowClass(balance.Balance.Budgets[year]),
			},
		})
	}

	return sections
}

func subProductsToSection(subProductData []html.ProductData, year model.BudgetYear, p *message.Printer) *html.BalanceSection {
	if len(subProductData) == 0 {
		return nil
	}

	var incomeCashflowTotal = 0.0
	var expensesCashflowTotal = 0.0

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

	for _, subProduct := range subProductData {
		var cashflow = 0.0
		for _, balance := range subProduct.FinancialPlan.Balances {
			cashflow += balance.Budgets[year]
		}

		link := fmt.Sprintf("%s/product.html", subProduct.Metadata.SubProduct.ID)

		if cashflow < 0 {
			section.Expenses.Labels = append(section.Expenses.Labels, subProduct.Metadata.SubProduct.Name)
			section.Expenses.Data = append(section.Expenses.Data, cashflow)
			section.ExpensesSubProductLinks = append(section.ExpensesSubProductLinks, link)

			expensesCashflowTotal += cashflow
		} else {
			section.Income.Labels = append(section.Income.Labels, subProduct.Metadata.SubProduct.Name)
			section.Income.Data = append(section.Income.Data, cashflow)
			section.IncomeSubProductLinks = append(section.IncomeSubProductLinks, link)

			incomeCashflowTotal += cashflow
		}
	}

	if shared.IsUnequal(expensesCashflowTotal, 0.0) {
		section.HasExpenses = true
		section.HasMoreThanOneExpense = len(section.ExpensesSubProductLinks) > 1
		section.ExpensesCashflowTotal = expensesCashflowTotal
		section.HasExpensesSubProductLinks = true
	}
	if shared.IsUnequal(incomeCashflowTotal, 0.0) {
		section.HasIncome = true
		section.HasMoreThanOneIncome = len(section.IncomeSubProductLinks) > 1
		section.IncomeCashflowTotal = incomeCashflowTotal
		section.HasIncomeSubProductLinks = true
	}

	if section.HasIncome && section.HasExpenses {
		section.HasIncomeAndExpenses = true
	}

	cashflowTotal := incomeCashflowTotal + expensesCashflowTotal
	section.Copy = html.BalanceSectionCopy{
		Header:                encodeSubProductBalanceSectionHeader(cashflowTotal, p),
		IncomeCashflowTotal:   "Einnahmen: " + encodeHtml.EncodeBudget(incomeCashflowTotal, p),
		ExpensesCashflowTotal: "Ausgaben: " + encodeHtml.EncodeBudget(expensesCashflowTotal, p),
	}

	section.CSS = html.BalanceSectionCSS{
		CashflowTotalClass: encodeHtml.EncodeCSSCashflowClass(cashflowTotal),
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

func dataPointsToChartJSDataset(dataPoints []html.DataPoint) html.ChartJSDataset {
	var labels = []string{}
	var data = []float64{}

	for _, dataPoint := range dataPoints {
		labels = append(labels, dataPoint.Label)
		data = append(data, dataPoint.Budget)
	}

	return html.ChartJSDataset{
		ID:           strings.ReplaceAll("chartjs-"+uuid.New().String(), "-", "_"),
		Labels:       labels,
		DatasetLabel: "Budget",
		Data:         data,
	}
}

func encodeBalanceSectionHeader(balance model.AccountBalance, year model.BudgetYear, p *message.Printer) template.HTML {
	return template.HTML(fmt.Sprintf(
		`- %s <span class="%s">%s</span>`,
		encodeAccountClass(balance.Class, balance.Budgets[year]),
		encodeHtml.EncodeCSSCashflowClass(balance.Budgets[year]),
		encodeHtml.EncodeBudget(balance.Budgets[year], p),
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

func encodeAccountClass(class model.AccountClass, cashflowTotal float64) string {
	switch class {
	case model.AccountClassAdministration:
		if cashflowTotal >= 0 {
			return "Die Verwaltung erwirtschaftet"
		}
		return "Die Verwaltung kostet"

	case model.AccountClassInvestments:
		if cashflowTotal >= 0 {
			return "Investitionen erwirtschaften"
		}
		return "Investitionen kosten"

	case model.AccountClassOneOff:
		if cashflowTotal >= 0 {
			return "Investitionen oberhalb der Wertgrenze erwirtschaften"
		}
		return "Investitionen oberhalb der Wertgrenze kosten"

	default:
		panic(fmt.Sprintf("unknown account class '%s'", class))
	}
}
