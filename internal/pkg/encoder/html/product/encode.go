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
	fpbBalanceSection := shared.Map(fpbBalanceDataOpt, func(fpbBalanceData []html.BalanceData) html.BalanceSection {
		return compressPlanBBalancesToSingleSection(fpbBalanceData, year, p)
	})

	return html.Product{
		Meta:                 metadata,
		FpaBalanceSections:   balanceDataToSections(fpaBalanceData, year, p),
		HasFpbBalanceSection: fpbBalanceSection.IsSome,
		FpbBalanceSection:    fpbBalanceSection.GetOrElse(html.BalanceSection{}),
		Copy: html.ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",

			IntroCashflowTotal: fmt.Sprintf("Das Produkt - %s - wird in %s", metadata.Description, year),
			IntroDescription:   encodeIntroDescription(fpaCashflowTotal+fpbCashflowTotalOpt.GetOrElse(0), metadata),

			CashflowTotal: encodeHtml.EncodeBudget(fpaCashflowTotal+fpbCashflowTotalOpt.GetOrElse(0), p),

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

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>
			<br><br>
			Das Gesamtbudget auf dieser Webseite ergibt sich aus dem Teilfinanzplan A und B.`,
		},
		CSS: html.ProductCSS{
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(fpaCashflowTotal + fpbCashflowTotalOpt.GetOrElse(0)),
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

func compressPlanBBalancesToSingleSection(balanceData []html.BalanceData, year model.BudgetYear, p *message.Printer) html.BalanceSection {
	var income []html.DataPoint
	var expenses []html.DataPoint

	for _, balance := range balanceData {
		income = append(income, balance.Income...)
		expenses = append(expenses, balance.Expenses...)
	}

	var incomeCashflowTotal float64
	for _, income := range income {
		incomeCashflowTotal += income.Budget
	}
	var expensesCashflowTotal float64
	for _, expense := range expenses {
		expensesCashflowTotal += expense.Budget
	}
	cashflowTotal := incomeCashflowTotal - expensesCashflowTotal

	return html.BalanceSection{
		ID:                   "balance-" + uuid.New().String(),
		HasIncomeAndExpenses: len(income) > 0 && len(expenses) > 0,
		HasIncome:            len(income) > 0,
		IncomeCashflowTotal:  incomeCashflowTotal,
		Income:               dataPointsToChartJSDataset(income),

		HasExpenses:           len(expenses) > 0,
		ExpensesCashflowTotal: expensesCashflowTotal,
		Expenses:              dataPointsToChartJSDataset(expenses),

		Copy: html.BalanceSectionCopy{
			Header: template.HTML(fmt.Sprintf(
				`%s <span class="%s">%s</span>`,
				encodeAccountClass(model.AccountClassOneOff, cashflowTotal),
				encodeHtml.EncodeCSSCashflowClass(cashflowTotal),
				encodeHtml.EncodeBudget(cashflowTotal, p),
			)),
			IncomeCashflowTotal:   "Einnahmen: " + encodeHtml.EncodeBudget(incomeCashflowTotal, p),
			ExpensesCashflowTotal: "Ausgaben: " + encodeHtml.EncodeBudget(expensesCashflowTotal, p),
		},
		CSS: html.BalanceSectionCSS{
			CashflowTotalClass: encodeHtml.EncodeCSSCashflowClass(cashflowTotal),
		},
	}
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
		encodeAccountClass(balance.Class, balance.Budgets[year]),
		encodeHtml.EncodeCSSCashflowClass(balance.Budgets[year]),
		encodeHtml.EncodeBudget(balance.Budgets[year], p),
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
			return "Investitionen unterhalb der Wertgrenze erwirtschaften"
		}
		return "Investitionen unterhalb der Wertgrenze kosten"

	case model.AccountClassOneOff:
		if cashflowTotal >= 0 {
			return "Investitionen oberhalb der Wertgrenze erwirtschaften"
		}
		return "Investitionen oberhalb der Wertgrenze kosten"

	default:
		panic(fmt.Sprintf("unknown account class '%s'", class))
	}
}
