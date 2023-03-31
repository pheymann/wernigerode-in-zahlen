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
			Year:     year,
			Headline: "Wernigerode in Zahlen",
			IntroCashflowTotal: template.HTML(fmt.Sprintf(`Als Teil unserer Stadt und Gemeinde habe ich mich gefragt, wo wir eigentlich unsere
			Gelder investieren. Nach einigem Suchen habe ich den Wernigeröder <a href="https://www.wernigerode.de/B%%C3%%BCrgerservice/Stadtrat/Haushaltsplan/">Haushaltsplan</a>
			gefunden. Ein mehrere hundert Seiten langes Dokument, das genau aufführt, wo Gelder gebraucht werden. Leider ist es alles andere als leicht zu lesen und so
			ist die Idee für diese Webseite bei mir entstanden. Es soll eine Darstellung des Finanzhaushalts Wernigerodes sein, die gut zu lesen und verstehen ist.
			<br><br>
			Alles startet mit der Gesamtübersicht. In %s planen wir`, year)),
			IntroDescription: encodeIntroDescription(cashflowTotal, len(departments)),

			CashflowTotal:         htmlEncoder.EncodeBudget(cashflowTotal, p),
			IncomeCashflowTotal:   "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal: "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			AdditionalInfo: `Aktuell bildet diese Webseite die Finanzdaten aus den Teilfinanzplänen A und B ab. Zusätzliche finanzielle Mittel zum Beispiel aus
			dem Finanzmittelüberschuss sind nicht enthalten. Die Gesamtausgaben würden sich dann auf <strong>-3.284.100,00€</strong> reduzieren (siehe Haushaltsplan).`,
			Departments: shared.MapSlice(departments, func(department model.CompressedDepartment) html.OverviewDepartmentCopy {
				return encodeCompressedDepartment(department, p)
			}),
			AdditionalInfoAfterTable: `Du willst dir die Daten selber mal anschauen? Kein Problem. <a href="https://github.com/pheymann/wernigerode-in-zahlen/tree/main/assets">Hier</a> findest du eine Zusammenfassung der Daten. Die CSV Datei
			kann einfach in ein Tabellenprogramm importiert werden. Falls du dich mit Datenanalyse auskennst, habe ich auch noch JSON Dateien bereitgestellt.`,

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>.`,
		},
		CSS: html.OverviewCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(cashflowTotal),
		},
	}
}

func encodeCompressedDepartment(department model.CompressedDepartment, p *message.Printer) html.OverviewDepartmentCopy {
	return html.OverviewDepartmentCopy{
		Name:          department.DepartmentName,
		CashflowTotal: htmlEncoder.EncodeBudget(department.CashflowTotal, p),
		Link:          department.GetDepartmentLink(),
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
