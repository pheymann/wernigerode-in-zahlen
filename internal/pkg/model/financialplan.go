package model

import (
	"fmt"
)

type ID = string

type FinancialPlanCity struct {
	AdministrationBalance Cashflow
	InvestmentsBalance    Cashflow
	Cashflow              Cashflow
	Departments           map[ID]FinancialPlanDepartment
}

type FinancialPlanDepartment struct {
	ID                    ID
	Name                  string
	AdministrationBalance Cashflow
	InvestmentsBalance    Cashflow
	Cashflow              Cashflow
	Products              map[ID]FinancialPlanProduct
}

func (department FinancialPlanDepartment) CreateLink() string {
	return fmt.Sprintf("/%s/department.html", department.ID)
}

type FinancialPlanProduct struct {
	ID                    ID
	AdministrationBalance AccountBalance
	InvestmentsBalance    AccountBalance
	Cashflow              Cashflow
	Metadata              Metadata
	SubProducts           []FinancialPlanProduct
}

func NewFinancialPlanProduct() *FinancialPlanProduct {
	return &FinancialPlanProduct{
		AdministrationBalance: AccountBalance{
			Type:     AccountBalanceTypeAdministration,
			Cashflow: NewCashFlow(),
			Accounts: make([]Account, 0),
		},
		InvestmentsBalance: AccountBalance{
			Type:     AccountBalanceTypeInvestments,
			Cashflow: NewCashFlow(),
			Accounts: make([]Account, 0),
		},
		Cashflow:    NewCashFlow(),
		SubProducts: make([]FinancialPlanProduct, 0),
	}
}

func (product FinancialPlanProduct) CreateLink() string {
	return product.GetPath() + "product.html"
}

func (product FinancialPlanProduct) GetPath() string {
	var productPath = fmt.Sprintf(
		"/%s/%s/%s/%s/%s/",
		product.Metadata.Department.ID,
		product.Metadata.ProductClass.ID,
		product.Metadata.ProductDomain.ID,
		product.Metadata.ProductGroup.ID,
		product.Metadata.Product.ID,
	)

	if product.IsSubProduct() {
		productPath = fmt.Sprintf("%s%s/", productPath, product.Metadata.SubProduct.ID)
	}

	return productPath
}

func (product FinancialPlanProduct) IsSubProduct() bool {
	return product.Metadata.SubProduct != nil
}

type AccountBalance struct {
	Type     AccountBalanceType
	Cashflow Cashflow
	Accounts []Account
}

type AccountBalanceType = string

const (
	AccountBalanceTypeAdministration AccountBalanceType = "administration"
	AccountBalanceTypeInvestments    AccountBalanceType = "investments"
)

type Account struct {
	ID          string
	ProductID   string
	Description string
	Type        AccountType
	Budget      map[BudgetYear]float64
}

type AccountType = string

const (
	AccountTypeExpense AccountType = "expense"
	AccountTypeIncome  AccountType = "income"
)

type Cashflow struct {
	Total    map[BudgetYear]float64
	Income   map[BudgetYear]float64
	Expenses map[BudgetYear]float64
}

func (cf Cashflow) AddCashflow(other Cashflow) Cashflow {
	for year, total := range other.Total {
		cf.Total[year] += total
	}
	for year, income := range other.Income {
		cf.Income[year] += income
	}
	for year, expenses := range other.Expenses {
		cf.Expenses[year] += expenses
	}
	return cf
}

func NewCashFlow() Cashflow {
	return Cashflow{
		Total:    make(map[BudgetYear]float64),
		Income:   make(map[BudgetYear]float64),
		Expenses: make(map[BudgetYear]float64),
	}
}

type BudgetYear = string

const (
	BudgetYear2020 BudgetYear = "2020"
	BudgetYear2021 BudgetYear = "2021"
	BudgetYear2022 BudgetYear = "2022"
	BudgetYear2023 BudgetYear = "2023"
	BudgetYear2024 BudgetYear = "2024"
	BudgetYear2025 BudgetYear = "2025"
	BudgetYear2026 BudgetYear = "2026"
)
