package html

import "html/template"

type Department struct {
	HasIncomeAndExpenses bool

	HasIncome          bool
	IncomeProductLinks []string
	Income             ChartJSDataset

	HasExpenses          bool
	ExpensesProductLinks []string
	Expenses             ChartJSDataset

	Copy DepartmentCopy
	CSS  DepartmentCSS
}

type DepartmentCopy struct {
	Year               string
	Department         string
	IntroCashflowTotal string
	IntroDescription   template.HTML

	CashflowTotal          string
	CashflowAdministration string
	CashflowInvestments    string
	IncomeCashflowTotal    string
	ExpensesCashflowTotal  string

	Products []DepartmentProductCopy

	BackLink string

	DataDisclosure template.HTML
}

type DepartmentProductCopy struct {
	Name                   string
	Link                   string
	CashflowTotal          string
	CashflowAdministration string
	CashflowInvestments    string
}

type DepartmentCSS struct {
	TotalCashflowClass string
}
