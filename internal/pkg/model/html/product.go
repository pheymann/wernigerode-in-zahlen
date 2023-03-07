package html

import (
	"html/template"

	"wernigode-in-zahlen.de/internal/pkg/model"
)

type Product struct {
	Meta               model.Metadata
	FpaBalanceSections []BalanceSection
	FpbBalanceSections []BalanceSection
	Copy               ProductCopy
	CSS                ProductCSS
}

type ProductCopy struct {
	BackLink string

	IntroCashflowTotal string
	IntroDescription   string

	CashflowTotal string

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
}

type ProductCSS struct {
	TotalCashflowClass string
}

type BalanceSection struct {
	ID string

	HasIncomeAndExpenses bool
	HasIncome            bool
	IncomeCashflowTotal  float64
	Income               ChartJSDataset

	HasExpenses           bool
	ExpensesCashflowTotal float64
	Expenses              ChartJSDataset

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
