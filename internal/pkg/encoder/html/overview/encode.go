package overview

import (
	"fmt"
	"html/template"
	"sort"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	htmlEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Encode(
	plan model.FinancialPlanCity,
	year model.BudgetYear,

	incomeDepartmentLinks []string,
	chartIncomeDataPerProduct html.ChartJSDataset,

	expensesDepartmentLinks []string,
	chartExpensesDataPerProduct html.ChartJSDataset,
) html.Overview {
	p := message.NewPrinter(language.German)

	departmentsTable := shared.ReduceMap(shared.MapMap(plan.Departments, func(department model.FinancialPlanDepartment) html.OverviewDepartmentCopy {
		return encodeDepartment(department, year, p)
	}))

	sort.Slice(departmentsTable, func(i, j int) bool {
		return departmentsTable[i].Name < departmentsTable[j].Name
	})

	return html.Overview{
		HasIncome:             plan.Cashflow.Income[year] > 0,
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
			IntroDescription: encodeIntroDescription(plan.Cashflow.Total[year], len(plan.Departments)),

			CashflowTotal:          htmlEncoder.EncodeBudget(plan.Cashflow.Total[year], p),
			CashflowAdministration: htmlEncoder.EncodeBudget(plan.AdministrationBalance.Total[year], p),
			CashflowInvestments:    htmlEncoder.EncodeBudget(plan.InvestmentsBalance.Total[year], p),
			IncomeCashflowTotal:    "Einnahmen: " + htmlEncoder.EncodeBudget(plan.Cashflow.Income[year], p),
			ExpensesCashflowTotal:  "Ausgaben: " + htmlEncoder.EncodeBudget(plan.Cashflow.Expenses[year], p),

			AdditionalInfo: `Aktuell bildet diese Webseite die Finanzdaten aus den Teilfinanzplänen A ab. Zusätzliche finanzielle Mittel zum Beispiel aus
			dem Finanzmittelüberschuss sind nicht enthalten. Die Gesamtausgaben würden sich dann auf <strong>-3.284.100,00€</strong> reduzieren (siehe Haushaltsplan).
			Zudem besteht eine Differenz wie Konten zusammengerechnet werden. Der Haushaltsplan summiert alle Einnahmen und Ausgaben für laufende Verwaltungstätigkeiten und
			Investitionen separat auf. Diese Webseite dagegen summiert Einnahmen und Ausgaben basierend auf Produkten und Fachbereichen. Die finale Differenz stimmt jedoch
			am Ende wieder überein.`,
			Departments: departmentsTable,
			AdditionalInfoAfterTable: `Du willst dir die Daten selber mal anschauen? Kein Problem. <a href="https://github.com/pheymann/wernigerode-in-zahlen/tree/main/assets">Hier</a> findest du eine Zusammenfassung der Daten. Die CSV Datei
			kann einfach in ein Tabellenprogramm importiert werden. Falls du dich mit Datenanalyse auskennst, habe ich auch noch JSON Dateien bereitgestellt.`,

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>.`,
		},
		CSS: html.OverviewCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(plan.Cashflow.Total[year]),
		},
	}
}

func encodeDepartment(department model.FinancialPlanDepartment, year model.BudgetYear, p *message.Printer) html.OverviewDepartmentCopy {
	return html.OverviewDepartmentCopy{
		Name:                   department.Name,
		CashflowTotal:          htmlEncoder.EncodeBudget(department.Cashflow.Total[year], p),
		CashflowAdministration: htmlEncoder.EncodeBudget(department.AdministrationBalance.Total[year], p),
		CashflowInvestments:    htmlEncoder.EncodeBudget(department.InvestmentsBalance.Total[year], p),
		Link:                   department.CreateLink(),
	}
}

func encodeIntroDescription(cashflowTotal float64, numberOfProducts int) template.HTML {
	var earnOrExpese = "einzunehmen"
	if cashflowTotal < 0 {
		earnOrExpese = "auszugeben"
	}

	return template.HTML(fmt.Sprintf(
		`%s. Die Gelder teilen sich auf <b>%d Fachbereiche</b> auf und setzen sich aus den laufenden Verwaltungstätigkeiten
		und gesonderten Investitionen, wie zum Beispiel Baumaßnahmen, zusammen. Klicke auf einen in den Diagrammen um mehr zu erfahren.`,
		earnOrExpese,
		numberOfProducts,
	))
}
