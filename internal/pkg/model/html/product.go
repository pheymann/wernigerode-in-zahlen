package html

import (
	"html/template"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

type Product struct {
	Meta            model.Metadata
	BalanceSections []BalanceSection
	Copy            ProductCopy
	CSS             ProductCSS
}

type ProductCopy struct {
	BackLink string
	Year     string

	IntroCashflowTotal string
	IntroDescription   string

	CashflowTotal string

	Accounts []ProductAccountCopy

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

	DataDisclosure template.HTML
}

type ProductAccountCopy struct {
	Name          string
	CashflowTotal string
}

type ProductCSS struct {
	TotalCashflowClass string
}

type BalanceSection struct {
	ID string

	HasIncomeAndExpenses bool

	HasIncome                bool
	IncomeCashflowTotal      float64
	HasIncomeSubProductLinks bool
	IncomeSubProductLinks    []string
	Income                   ChartJSDataset
	IncomeID                 template.JS

	HasExpenses                bool
	ExpensesCashflowTotal      float64
	HasExpensesSubProductLinks bool
	ExpensesSubProductLinks    []string
	Expenses                   ChartJSDataset
	ExpensesID                 template.JS

	Copy BalanceSectionCopy
	CSS  BalanceSectionCSS
}

type BalanceSectionCopy struct {
	Header                template.HTML
	IncomeCashflowTotal   string
	ExpensesCashflowTotal string
}

type BalanceSectionCSS struct {
	CashflowTotalClass string
}
