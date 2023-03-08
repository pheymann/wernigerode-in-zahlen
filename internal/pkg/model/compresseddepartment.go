package model

type CompressedDepartment struct {
	ID                     string
	DepartmentName         string
	CashflowTotal          float64
	CashflowFinancialPlanA float64
	CashflowFinancialPlanB float64
	NumberOfProducts       int
}
