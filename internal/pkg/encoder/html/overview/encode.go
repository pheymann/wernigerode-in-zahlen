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
			Year:               year,
			Headline:           "Wernigerode in Zahlen",
			IntroCashflowTotal: template.HTML(fmt.Sprintf(`Im Jahr %s werden wir`, year)),
			IntroDescription:   encodeIntroDescription(plan.Cashflow.Total[year], len(plan.Departments)),

			CashflowTotal:          htmlEncoder.EncodeBudget(plan.Cashflow.Total[year], p),
			CashflowAdministration: htmlEncoder.EncodeBudget(plan.AdministrationBalance.Total[year], p),
			CashflowInvestments:    htmlEncoder.EncodeBudget(plan.InvestmentsBalance.Total[year], p),
			IncomeCashflowTotal:    "Einnahmen: " + htmlEncoder.EncodeBudget(plan.Cashflow.Income[year], p),
			ExpensesCashflowTotal:  "Ausgaben: " + htmlEncoder.EncodeBudget(plan.Cashflow.Expenses[year], p),

			AdditionalInfo: `Aktuell bildet diese Webseite die Finanzdaten aus den Teilfinanzplänen A (Verwaltungs- und Investitionstätigkeiten) ab. Zusätzliche Ausgaben und Einnahmen zum Beispiel aus
			dem Finanzierungstätigkeiten sind nicht enthalten. Die Gesamtausgaben belaufen sich für 2023 auf <strong>-3.379.700€</strong> (siehe Haushaltsplan).
			Zudem besteht eine Differenz wie Konten zusammengerechnet werden. Der Haushaltsplan summiert alle Einnahmen und Ausgaben für laufende Verwaltungstätigkeiten und
			Investitionen separat auf. Diese Webseite dagegen summiert Einnahmen und Ausgaben basierend auf Produkten und Fachbereichen. Die finalen Werte stimmen jedoch
			am Ende wieder überein.`,
			Departments: departmentsTable,
			AdditionalInfoAfterTable: `Du willst dir die Daten selber mal anschauen? Kein Problem. <a href="https://github.com/pheymann/wernigerode-in-zahlen/tree/main/assets">Hier</a> findest du eine Zusammenfassung der Daten. Die CSV Datei
			kann einfach in ein Tabellenprogramm importiert werden. Falls du dich mit Datenanalyse auskennst, habe ich auch noch JSON Dateien bereitgestellt.`,

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Stadt Wernigerode aus dem Jahr 2023.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			<a href="https://www.wernigerode.de/buergerinformationssystem/vo020.asp?VOLFDNR=3344">hier findet</a>.`,
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
