package html

import "wernigode-in-zahlen.de/internal/pkg/model"

type Product struct {
	Meta            model.Metadata
	BalanceSections []BalanceSection
	Copy            ProductCopy
	CSS             ProductCSS
}

type ProductCopy struct {
	BackLink string

	IntroCashflowTotal string
	IntroDescription   string

	CashflowTotal    string
	CashflowIncome   string
	CashflowExpenses string

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
	Header                string
	IncomeCashflowTotal   string
	ExpensesCashflowTotal string
}

type BalanceSectionCSS struct {
	CashflowTotalClass string
}
