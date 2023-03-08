package html

import "html/template"

type Department struct {
	IncomeProductLinks []string
	Income             ChartJSDataset

	ExpensesProductLinks []string
	Expenses             ChartJSDataset

	Copy DepartmentCopy
	CSS  DepartmentCSS
}

type DepartmentCopy struct {
	Department         string
	IntroCashflowTotal string
	IntroDescription   template.HTML

	CashflowTotal         string
	IncomeCashflowTotal   string
	ExpensesCashflowTotal string

	Products []DepartmentProductCopy

	BackLink string

	DataDisclosure template.HTML
}

type DepartmentProductCopy struct {
	Name      string
	Link      string
	CashflowA string
	CashflowB string
}

type DepartmentCSS struct {
	TotalCashflowClass string
}
