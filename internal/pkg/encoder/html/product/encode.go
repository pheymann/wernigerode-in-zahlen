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

			IntroCashflowTotal: fmt.Sprintf("Das Produkt - %s - wird in %s", metadata.Description, year),
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

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Statdt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>
			<br><br>
			Das Gesamtbudget auf dieser Webseite ergibt sich aus dem Teilfinanzplan A und B.`,
		},
		CSS: html.ProductCSS{
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(fpaCashflowTotal),
		},
	}
}

func encodeIntroDescription(cashflowTotal float64, meta model.Metadata) string {
	if cashflowTotal >= 0 {
		return "einbringen."
	}
	return "kosten."
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
		`%s <span class="%s">%s</span>`,
		encodeAccountClass(balance.Class, balance.Budgets[year], balance.Desc),
		encodeHtml.EncodeCSSCashflowClass(balance.Budgets[year]),
		encodeHtml.EncodeBudget(balance.Budgets[year], p),
	))
}

func encodeAccountClass(class model.AccountClass, cashflowTotal float64, oneOffDesc string) string {
	switch class {
	case model.AccountClassAdministration:
		if cashflowTotal >= 0 {
			return "Die Verwaltung erwirtschaftet"
		}
		return "Die Verwaltung kostet"

	case model.AccountClassInvestments:
		if cashflowTotal >= 0 {
			return "Investitionen unterhalb der Wertgrenze erwirtschaften"
		}
		return "Investitionen unterhalb der Wertgrenze kosten"

	case model.AccountClassOneOff:
		if cashflowTotal >= 0 {
			return fmt.Sprintf("Die Investition \"%s\" erwirtschaftet", oneOffDesc)
		}
		return fmt.Sprintf("Die Investition \"%s\" kostet", oneOffDesc)

	default:
		panic(fmt.Sprintf("unknown account class '%s'", class))
	}
}
