package product

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/text/message"
	encodeHtml "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
)

func Encode(metadata model.Metadata, balanceData []html.BalanceData, cashflowTotal float64, year model.BudgetYear, p *message.Printer) html.Product {
	return html.Product{
		Meta:            metadata,
		BalanceSections: balanceDataToSections(balanceData, year, p),
		Copy: html.ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",

			IntroCashflowTotal: fmt.Sprintf("In %s haben wir", year),
			IntroDescription:   encodeIntroDescription(cashflowTotal, metadata),

			CashflowTotal: encodeHtml.EncodeBudget(cashflowTotal, p),

			MetaDepartment:    "Fachbereich",
			MetaProductClass:  "Produktklasse",
			MetaProductDomain: "Produktbereich",
			MetaProductGroup:  "Produktgruppe",
			MetaProduct:       "Produkt",
			MetaAccountable:   "Verantwortlich",
			MetaResponsible:   "Zuständig",
			MetaMission:       "Aufgabe",
			MetaTargets:       "Ziele",
			MetaServices:      "Dienstleistungen",
			MetaGrouping:      "Gruppierung",
		},
		CSS: html.ProductCSS{
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(cashflowTotal),
		},
	}
}

func encodeIntroDescription(cashflowTotal float64, meta model.Metadata) string {
	if cashflowTotal >= 0 {
		return fmt.Sprintf("eingenommen über: %s.", meta.Description)
	}
	return fmt.Sprintf("ausgegeben für: %s.", meta.Description)
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
			ID: "balance-" + uuid.New().String(),

			HasIncomeAndExpenses: len(balance.Income) > 0 && len(balance.Expenses) > 0,
			HasIncome:            len(balance.Income) > 0,
			IncomeCashflowTotal:  incomeCashflowTotal,
			Income:               dataPointsToChartJSDataset(balance.Income),

			HasExpenses:           len(balance.Expenses) > 0,
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

func dataPointsToChartJSDataset(dataPoints []html.DataPoint) html.ChartJSDataset {
	var labels = []string{}
	var data = []float64{}

	for _, dataPoint := range dataPoints {
		labels = append(labels, dataPoint.Label)
		data = append(data, dataPoint.Budget)
	}

	return html.ChartJSDataset{
		ID:           "chartjs-" + uuid.New().String(),
		Labels:       labels,
		DatasetLabel: "Budget",
		Data:         data,
	}
}

func encodeBalanceSectionHeader(balance model.AccountBalance, year model.BudgetYear, p *message.Printer) string {
	return fmt.Sprintf(
		"%s %s %s",
		encodeAccountClass(balance.Class),
		encodeHtml.EncodeBudget(balance.Budgets[year], p),
		encodeBalance(balance.Budgets[year]),
	)
}

func encodeAccountClass(class model.AccountClass) string {
	switch class {
	case model.AccountClassAdministration:
		return "Die Verwaltung hat dabei"
	case model.AccountClassInvestments:
		return "Die Investitionen haben dabei"
	}
	panic(fmt.Sprintf("unknown account class '%s'", class))
}

func encodeBalance(cashflowTotal float64) string {
	if cashflowTotal >= 0 {
		return "eingenommen"
	}
	return "gekostet"
}
