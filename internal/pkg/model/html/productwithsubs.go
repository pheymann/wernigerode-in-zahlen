package html

import (
	"html/template"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

type ProductWithSubs struct {
	Meta model.Metadata

	HasIncomeAndExpenses bool

	HasIncome          bool
	IncomeProductLinks []string
	Income             ChartJSDataset

	HasExpenses          bool
	ExpensesProductLinks []string
	Expenses             ChartJSDataset

	Copy ProductWithSubsCopy
	CSS  ProductWithSubsCSS
}

type ProductWithSubsCopy struct {
	Year     string
	BackLink string

	IntroCashflowTotal string
	IntroDescription   template.HTML

	CashflowTotal          string
	CashflowAdministration string
	CashflowInvestments    string
	IncomeCashflowTotal    string
	ExpensesCashflowTotal  string

	MetaDepartment    string
	MetaProductClass  string
	MetaProductDomain string
	MetaProductGroup  string
	MetaProduct       string
	MetaAccountable   string
	MetaResponsible   string
	MetaMission       string
	MetaTargets       string
	MetaServices      string
	MetaGrouping      string

	SubProducts []SubProductCopy

	DataDisclosure template.HTML
}

type SubProductCopy struct {
	Name                   string
	Link                   string
	CashflowTotal          string
	CashflowAdministration string
	CashflowInvestments    string
}

type ProductWithSubsCSS struct {
	TotalCashflowClass string
}
