package department

import (
	"fmt"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	htmlEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	department model.FinancialPlanDepartment,
	year model.BudgetYear,
	productTable []html.ProductTableData,

	incomeProductLinks []string,
	chartIncomeDataPerProduct html.ChartJSDataset,

	expensesProductLinks []string,
	chartExpensesDataPerProduct html.ChartJSDataset,
) html.Department {
	p := message.NewPrinter(language.German)

	hasIncome := shared.IsUnequal(department.Cashflow.Income[year], 0)
	hasExpenses := shared.IsUnequal(department.Cashflow.Expenses[year], 0)
	return html.Department{
		HasIncomeAndExpenses: hasIncome && hasExpenses,

		HasIncome:          hasIncome,
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		HasExpenses:          hasExpenses,
		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,

		Copy: html.DepartmentCopy{
			Year:               year,
			Department:         department.Name,
			IntroCashflowTotal: fmt.Sprintf("In %s planen wir", year),
			IntroDescription:   encodeIntroDescription(department.Cashflow.Total[year], len(productTable)),

			CashflowTotal:          htmlEncoder.EncodeBudget(department.Cashflow.Total[year], p),
			CashflowAdministration: htmlEncoder.EncodeBudget(department.AdministrationBalance.Total[year], p),
			CashflowInvestments:    htmlEncoder.EncodeBudget(department.InvestmentsBalance.Total[year], p),
			IncomeCashflowTotal:    "Einnahmen: " + htmlEncoder.EncodeBudget(department.Cashflow.Income[year], p),
			ExpensesCashflowTotal:  "Ausgaben: " + htmlEncoder.EncodeBudget(department.Cashflow.Expenses[year], p),

			Products: shared.MapSlice(productTable, func(productData html.ProductTableData) html.DepartmentProductCopy {
				return encodeDepartmentProductData(productData, p)
			}),

			BackLink: "Zurück zur Übersicht",

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2023.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			<a href="https://www.wernigerode.de/buergerinformationssystem/vo020.asp?VOLFDNR=3344">hier findet</a>.`,
		},
		CSS: html.DepartmentCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(department.Cashflow.Total[year]),
		},
	}
}

func encodeDepartmentProductData(data html.ProductTableData, p *message.Printer) html.DepartmentProductCopy {
	return html.DepartmentProductCopy{
		Name:                   data.Name,
		CashflowTotal:          htmlEncoder.EncodeBudget(data.CashflowTotal, p),
		CashflowAdministration: htmlEncoder.EncodeBudget(data.CashflowAdministration, p),
		CashflowInvestments:    htmlEncoder.EncodeBudget(data.CashflowInvestments, p),
		Link:                   data.Link,
	}
}

func encodeIntroDescription(cashflowTotal float64, numberOfProducts int) template.HTML {
	var earnOrExpese = "über diesen Fachbereich einzunehmen"
	if cashflowTotal < 0 {
		earnOrExpese = "für diesen Fachbereich auszugeben"
	}

	return template.HTML(fmt.Sprintf(
		"%s. Dabei geht das Geld an die folgenden <b>%d Produkte</b>. Klicke auf eines in den Diagrammen, um mehr zu erfahren.",
		earnOrExpese,
		numberOfProducts,
	))
}
