package model

type CostCenterType = string

const (
	CostCenterGroup CostCenterType = "group"
	CostCenterUnit  CostCenterType = "unit"
)

type FinancePlanACostCenter struct {
	Id         string
	Tpe        CostCenterType
	Desc       string
	Budget2020 float64
	Budget2021 float64
	Budget2022 float64
	Budget2023 float64
	Budget2024 float64
	Budget2025 float64
}

type FinancePlanA struct {
	Groups []FinancePlanACostCenter
	Units  map[string][]FinancePlanACostCenter
}
