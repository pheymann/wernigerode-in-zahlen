package productwithsubs

import (
	"fmt"
	"html/template"

	"golang.org/x/text/message"
	encodeHtml "wernigerode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	plan model.FinancialPlanProduct,
	year model.BudgetYear,

	incomeProductLinks []string,
	chartIncomeDataPerProduct html.ChartJSDataset,

	expensesProductLinks []string,
	chartExpensesDataPerProduct html.ChartJSDataset,

	p *message.Printer,
) html.ProductWithSubs {
	hasIncome := shared.IsUnequal(plan.Cashflow.Income[year], 0)
	hasExpenses := shared.IsUnequal(plan.Cashflow.Expenses[year], 0)

	return html.ProductWithSubs{
		Meta:                 plan.Metadata,
		HasIncomeAndExpenses: hasIncome && hasExpenses,

		HasIncome:          hasIncome,
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		HasExpenses:          hasExpenses,
		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,
		Copy: html.ProductWithSubsCopy{
			BackLink: "Zurück zur Bereichsübersicht",
			Year:     year,

			IntroCashflowTotal: fmt.Sprintf("Das Produkt - %s - wird in %s", plan.Metadata.Description, year),
			IntroDescription:   encodeIntroDescription(plan.Cashflow.Total[year], len(plan.SubProducts), plan.Metadata),

			CashflowTotal:         encodeHtml.EncodeBudget(plan.Cashflow.Total[year], p),
			IncomeCashflowTotal:   "Einnahmen: " + encodeHtml.EncodeBudget(plan.Cashflow.Income[year], p),
			ExpensesCashflowTotal: "Ausgaben: " + encodeHtml.EncodeBudget(plan.Cashflow.Expenses[year], p),

			SubProducts: encodeSubProducts(plan.SubProducts, p),

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
		CSS: html.ProductWithSubsCSS{
			TotalCashflowClass: encodeHtml.EncodeCSSCashflowClass(plan.Cashflow.Total[year]),
		},
	}
}

func encodeIntroDescription(cashflowTotal float64, numberOfSubProducts int, meta model.Metadata) template.HTML {
	var expenseEarnCopy = "kosten"
	if cashflowTotal >= 0 {
		expenseEarnCopy = "einbringen"
	}

	return template.HTML(fmt.Sprintf(
		`%s. Dabei geht das Geld an die folgenden <b>%d unter Produkte</b>. Klicke auf eines in den Diagrammen, um mehr zu erfahren.`,
		expenseEarnCopy,
		numberOfSubProducts,
	))
}

func encodeSubProducts(subProducts []model.FinancialPlanProduct, p *message.Printer) []html.SubProductCopy {
	var result []html.SubProductCopy

	for _, subProduct := range subProducts {
		result = append(result, html.SubProductCopy{
			Name:                   subProduct.Metadata.Description,
			CashflowTotal:          encodeHtml.EncodeBudget(subProduct.Cashflow.Total[model.BudgetYear2022], p),
			CashflowAdministration: encodeHtml.EncodeBudget(subProduct.AdministrationBalance.Cashflow.Total[model.BudgetYear2022], p),
			CashflowInvestments:    encodeHtml.EncodeBudget(subProduct.InvestmentsBalance.Cashflow.Total[model.BudgetYear2022], p),
			Link:                   subProduct.CreateLink(),
		})
	}

	return result
}
