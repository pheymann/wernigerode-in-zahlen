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
	Year               string
	Headline           string
	IntroCashflowTotal template.HTML
	IntroDescription   template.HTML

	CashflowTotal          string
	CashflowAdministration string
	CashflowInvestments    string
	IncomeCashflowTotal    string
	ExpensesCashflowTotal  string

	AdditionalInfo           template.HTML
	Departments              []OverviewDepartmentCopy
	AdditionalInfoAfterTable template.HTML

	BackLink string

	DataDisclosure template.HTML
}

type OverviewDepartmentCopy struct {
	Name                   string
	Link                   string
	CashflowTotal          string
	CashflowAdministration string
	CashflowInvestments    string
}

type OverviewCSS struct {
	TotalCashflowClass string
}
