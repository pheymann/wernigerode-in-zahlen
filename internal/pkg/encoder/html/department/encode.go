package department

import (
	"fmt"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	htmlEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	compressed model.CompressedDepartment,
	year model.BudgetYear,
	productData []html.DepartmentProductData,

	incomeTotalCashFlow float64,
	incomeProductLinks []string,
	chartIncomeDataPerProduct html.ChartJSDataset,

	expensesTotalCashFlow float64,
	expensesProductLinks []string,
	chartExpensesDataPerProduct html.ChartJSDataset,
) html.Department {
	p := message.NewPrinter(language.German)

	return html.Department{
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,

		Copy: html.DepartmentCopy{
			Department:         compressed.DepartmentName,
			IntroCashflowTotal: fmt.Sprintf("In %s planen wir", year),
			IntroDescription:   encodeIntroDescription(compressed.CashflowTotal, compressed.NumberOfProducts),

			CashflowTotal:          htmlEncoder.EncodeBudget(compressed.CashflowTotal, p),
			CashflowFinancialPlanA: htmlEncoder.EncodeBudget(compressed.CashflowFinancialPlanA, p),
			CashflowFinancialPlanB: htmlEncoder.EncodeBudget(compressed.CashflowFinancialPlanB, p),
			IncomeCashflowTotal:    "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal:  "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			Products: shared.MapSlice(productData, func(productData html.DepartmentProductData) html.DepartmentProductCopy {
				return encodeDepartmentProductData(productData, p)
			}),

			BackLink: "Zurück zur Übersicht",

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Statdt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>
			<br><br>
			Die Budgets auf dieser Webseite ergeben sich aus dem Teilfinanzplan A und B und weichen damit vom Haushaltsplan ab, der
			nur Teilfinanzplan A Daten enthält.`,
		},
		CSS: html.DepartmentCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(compressed.CashflowTotal),
		},
	}
}

func encodeDepartmentProductData(data html.DepartmentProductData, p *message.Printer) html.DepartmentProductCopy {
	return html.DepartmentProductCopy{
		Name:      data.Name,
		CashflowA: htmlEncoder.EncodeBudget(data.CashflowFinancialPlanA, p),
		CashflowB: htmlEncoder.EncodeBudget(data.CashflowFinancialPlanB, p),
		Link:      data.Link,
	}
}

func encodeIntroDescription(cashflowTotal float64, numberOfProducts int) template.HTML {
	var earnOrExpese = "einzunehmen"
	if cashflowTotal < 0 {
		earnOrExpese = "auszugeben"
	}

	return template.HTML(fmt.Sprintf(
		"%s. Klick auf eines der <b>%d Produkte</b> in den Diagrammen um mehr zu erfahren.",
		earnOrExpese,
		numberOfProducts,
	))
}
