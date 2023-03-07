package product

import (
	"fmt"
	"html/template"

	"github.com/google/uuid"
	"golang.org/x/text/message"
	encodeHtml "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	metadata model.Metadata,
	fpaBalanceData []html.BalanceData,
	fpaCashflowTotal float64,
	fpbBalanceDataOpt shared.Option[[]html.BalanceData],
	fpbCashflowTotalOpt shared.Option[float64],
	year model.BudgetYear,
	p *message.Printer,
) html.Product {
	return html.Product{
		Meta:               metadata,
		FpaBalanceSections: balanceDataToSections(fpaBalanceData, year, p),
		FpbBalanceSections: shared.Map(fpbBalanceDataOpt, func(fpbBalanceData []html.BalanceData) []html.BalanceSection {
			return balanceDataToSections(fpbBalanceData, year, p)
		}).GetOrElse([]html.BalanceSection{}),
		Copy: html.ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",

			IntroCashflowTotal: fmt.Sprintf("In %s haben wir", year),
			IntroDescription:   encodeIntroDescription(fpaCashflowTotal+fpbCashflowTotalOpt.GetOrElse(0), metadata),

			CashflowTotal: encodeHtml.EncodeBudget(fpaCashflowTotal, p),

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
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(fpaCashflowTotal),
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

func encodeBalanceSectionHeader(balance model.AccountBalance, year model.BudgetYear, p *message.Printer) template.HTML {
	return template.HTML(fmt.Sprintf(
		`%s <span class="%s">%s</span> %s`,
		encodeAccountClass(balance.Class, balance.Desc),
		encodeHtml.EncodeCSSCashflowClass(balance.Budgets[year]),
		encodeHtml.EncodeBudget(balance.Budgets[year], p),
		encodeBalance(balance.Budgets[year]),
	))
}

func encodeAccountClass(class model.AccountClass, oneOffDesc string) string {
	switch class {
	case model.AccountClassAdministration:
		return "Die Verwaltung hat dabei"
	case model.AccountClassInvestments:
		return "Die Investitionen haben dabei"
	case model.AccountClassOneOff:
		return fmt.Sprintf("Das Budget \"%s\" hat dabei", oneOffDesc)
	default:
		panic(fmt.Sprintf("unknown account class '%s'", class))
	}
}

func encodeBalance(cashflowTotal float64) string {
	if cashflowTotal >= 0 {
		return "eingenommen"
	}
	return "gekostet"
}
