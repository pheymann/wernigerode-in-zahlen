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
	compressed model.CompressedDepartment,
	year model.BudgetYear,
	productData []html.ProductTableData,

	incomeTotalCashFlow float64,
	incomeProductLinks []string,
	chartIncomeDataPerProduct html.ChartJSDataset,

	expensesTotalCashFlow float64,
	expensesProductLinks []string,
	chartExpensesDataPerProduct html.ChartJSDataset,
) html.Department {
	p := message.NewPrinter(language.German)

	hasIncome := shared.IsUnequal(incomeTotalCashFlow, 0)
	hasExpenses := shared.IsUnequal(expensesTotalCashFlow, 0)
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
			Department:         compressed.DepartmentName,
			IntroCashflowTotal: fmt.Sprintf("In %s planen wir", year),
			IntroDescription:   encodeIntroDescription(compressed.CashflowTotal, compressed.NumberOfProducts),

			CashflowTotal:          htmlEncoder.EncodeBudget(compressed.CashflowTotal, p),
			CashflowAdministration: htmlEncoder.EncodeBudget(compressed.CashflowAdministration, p),
			CashflowInvestments:    htmlEncoder.EncodeBudget(compressed.CashflowInvestments, p),
			IncomeCashflowTotal:    "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal:  "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			Products: shared.MapSlice(productData, func(productData html.ProductTableData) html.DepartmentProductCopy {
				return encodeDepartmentProductData(productData, p)
			}),

			BackLink: "Zurück zur Übersicht",

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>.`,
		},
		CSS: html.DepartmentCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(compressed.CashflowTotal),
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
