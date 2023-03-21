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

	CashflowTotal         string
	CashflowB             string
	IncomeCashflowTotal   string
	ExpensesCashflowTotal string

	AdditionalInfo template.HTML
	Departments    []OverviewDepartmentCopy

	BackLink string

	DataDisclosure template.HTML
}

type OverviewDepartmentCopy struct {
	Name          string
	Link          string
	CashflowTotal string
	CashflowB     string
}

type OverviewCSS struct {
	TotalCashflowClass string
}
