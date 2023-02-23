package model

type FinancialPlanA struct {
	RequiredMetadata RequiredMetadata
	Balances         []AccountBalance
}

type RequiredMetadata struct {
	Department  string
	Accountable string
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

func (fpa *FinancialPlanA) UpdateLastAccountBalance(f func(AccountBalance) AccountBalance) {
	lastBalanceIndex := len(fpa.Balances) - 1

	if lastBalanceIndex < 0 {
		fpa.Balances = append(fpa.Balances, f(AccountBalance{}))
	} else {
		fpa.Balances[lastBalanceIndex] = f(fpa.Balances[lastBalanceIndex])
	}
}

func (fpa *FinancialPlanA) UpdateLastAccount(f func(Account) Account) {
	lastBalanceIndex := len(fpa.Balances) - 1
	lastAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts) - 1

	if lastAccountIndex < 0 {
		fpa.Balances[lastBalanceIndex].Accounts = append(fpa.Balances[lastBalanceIndex].Accounts, f(Account{}))
	} else {
		lastAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts) - 1

		fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex] = f(fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex])
	}
}

func (fpa *FinancialPlanA) UpdateLastSubAccount(f func(SubAccount) SubAccount) {
	lastBalanceIndex := len(fpa.Balances) - 1
	lastAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts) - 1
	lastSubAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs) - 1

	if lastSubAccountIndex < 0 {
		fpa.Balances[lastBalanceIndex].Accounts[lastBalanceIndex].Subs = append(
			fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs,
			f(SubAccount{}),
		)
	} else {
		fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs[lastSubAccountIndex] = f(
			fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs[lastSubAccountIndex],
		)
	}
}

func (fpa *FinancialPlanA) UpdateLastUnitAccount(f func(UnitAccount) UnitAccount) {
	lastBalanceIndex := len(fpa.Balances) - 1
	lastAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts) - 1
	lastSubAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs) - 1
	lastUnitAccountIndex := len(fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs[lastSubAccountIndex].Units) - 1

	if lastUnitAccountIndex < 0 {
		fpa.Balances[lastBalanceIndex].Accounts[lastBalanceIndex].Subs[lastSubAccountIndex].Units = append(
			fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs[lastSubAccountIndex].Units,
			f(UnitAccount{}),
		)
	} else {
		fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs[lastSubAccountIndex].Units[lastUnitAccountIndex] = f(
			fpa.Balances[lastBalanceIndex].Accounts[lastAccountIndex].Subs[lastSubAccountIndex].Units[lastUnitAccountIndex],
		)
	}
}
