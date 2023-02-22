package model

type FinancialPlanA struct {
	Balances []AccountBalance
}

type AccountClass = string

const (
	AccountClassAdministration AccountClass = "admininstration"
	AccountClassInvestments    AccountClass = "balance-investments"
)

type AccountBalance struct {
	Id         string
	Class      AccountClass
	Desc       string
	Budget2020 float64
	Budget2021 float64
	Budget2022 float64
	Budget2023 float64
	Budget2024 float64
	Budget2025 float64
	Accounts   []Account
}

type Account struct {
	Id         string
	Desc       string
	Budget2020 float64
	Budget2021 float64
	Budget2022 float64
	Budget2023 float64
	Budget2024 float64
	Budget2025 float64
	Subs       []SubAccount
}

type SubAccount struct {
	Id         string
	Desc       string
	Budget2020 float64
	Budget2021 float64
	Budget2022 float64
	Budget2023 float64
	Budget2024 float64
	Budget2025 float64
	Units      []UnitAccount
}

type UnitAccount struct {
	Id         string
	Desc       string
	Budget2020 float64
	Budget2021 float64
	Budget2022 float64
	Budget2023 float64
	Budget2024 float64
	Budget2025 float64
}
