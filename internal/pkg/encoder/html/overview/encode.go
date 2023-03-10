package overview

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
	departments []model.CompressedDepartment,
	year model.BudgetYear,

	cashflowTotal float64,
	cashflowFinancialPlanA float64,
	cashflowFinancialPlanB float64,

	incomeTotalCashFlow float64,
	incomeDepartmentLinks []string,
	chartIncomeDataPerProduct html.ChartJSDataset,

	expensesTotalCashFlow float64,
	expensesDepartmentLinks []string,
	chartExpensesDataPerProduct html.ChartJSDataset,
) html.Overview {
	p := message.NewPrinter(language.German)

	return html.Overview{
		HasIncome:             incomeTotalCashFlow > 0,
		IncomeDepartmentLinks: incomeDepartmentLinks,
		Income:                chartIncomeDataPerProduct,

		ExpensesDepartmentLinks: expensesDepartmentLinks,
		Expenses:                chartExpensesDataPerProduct,

		Copy: html.OverviewCopy{
			Headline: "Wernigerode in Zahlen",
			IntroCashflowTotal: template.HTML(fmt.Sprintf(`Als Teil unserer Stadt und Gemeinde habe ich mich gefragt, wo wir eigentlich unser
			Geld investieren und nach einigem Suchen habe ich den <a href="https://www.wernigerode.de/B%%C3%%BCrgerservice/Stadtrat/Haushaltsplan/">Haushaltsplan</a>
			der Stadt gefunden. Der führt genau auf, wo Gelder gebraucht werden, ist aber alles andere als leicht zu lesen. Und so
			ist die Idee bei mir für diese Webseite entstanden. Es soll eine Darstellung des Finanzhaushalts Wernigerodes sein, die gut zu lesen und verstehen ist.
			<br><br>
			Alles startet mit der Gesamtübersicht. In %s planen wir`, year)),
			IntroDescription: encodeIntroDescription(cashflowTotal, len(departments)),

			CashflowTotal:          htmlEncoder.EncodeBudget(cashflowTotal, p),
			CashflowFinancialPlanA: htmlEncoder.EncodeBudget(cashflowFinancialPlanA, p),
			CashflowFinancialPlanB: htmlEncoder.EncodeBudget(cashflowFinancialPlanB, p),
			IncomeCashflowTotal:    "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal:  "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			Departments: shared.MapSlice(departments, func(department model.CompressedDepartment) html.OverviewDepartmentCopy {
				return encodeCompressedDepartment(department, p)
			}),

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>
			<br><br>
			Die Budgets auf dieser Webseite ergeben sich aus dem Teilfinanzplan A und B und weichen damit vom Haushaltsplan ab, der
			nur Teilfinanzplan A Daten enthält.`,
		},
		CSS: html.OverviewCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(cashflowTotal),
		},
	}
}

func encodeCompressedDepartment(department model.CompressedDepartment, p *message.Printer) html.OverviewDepartmentCopy {
	return html.OverviewDepartmentCopy{
		Name:      department.DepartmentName,
		CashflowA: htmlEncoder.EncodeBudget(department.CashflowFinancialPlanA, p),
		CashflowB: htmlEncoder.EncodeBudget(department.CashflowFinancialPlanB, p),
		Link:      department.GetDepartmentLink(),
	}
}

func encodeIntroDescription(cashflowTotal float64, numberOfProducts int) template.HTML {
	var earnOrExpese = "einzunehmen"
	if cashflowTotal < 0 {
		earnOrExpese = "auszugeben"
	}

	return template.HTML(fmt.Sprintf(
		"%s. Die Gelder teilen sich auf <b>%d Fachbereiche</b> auf. Klicke auf einen in den Diagrammen um mehr zu erfahren.",
		earnOrExpese,
		numberOfProducts,
	))
}