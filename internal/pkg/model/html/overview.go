package html

import (
	"html/template"
)

type Overview struct {
	HasIncome             bool
	IncomeDepartmentLinks []string
	Income                ChartJSDataset

	ExpensesDepartmentLinks []string
	Expenses                ChartJSDataset

	Copy OverviewCopy
	CSS  OverviewCSS
}

type OverviewCopy struct {
	Headline           string
	IntroCashflowTotal string
	IntroDescription   template.HTML

	CashflowTotal          string
	CashflowFinancialPlanA string
	CashflowFinancialPlanB string
	IncomeCashflowTotal    string
	ExpensesCashflowTotal  string

	Departments []OverviewDepartmentCopy

	BackLink string

	DataDisclosure template.HTML
}

type OverviewDepartmentCopy struct {
	Name      string
	Link      string
	CashflowA string
	CashflowB string
}

type OverviewCSS struct {
	TotalCashflowClass string
}
