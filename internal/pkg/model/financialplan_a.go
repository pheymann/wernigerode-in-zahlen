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

func (fpa *FinancialPlanA) AddAccountBalance(balance AccountBalance) {
	fpa.Balances = append(fpa.Balances, balance)
}

func (fpa *FinancialPlanA) RemoveLastAccountBalance() {
	fpa.Balances = fpa.Balances[:len(fpa.Balances)-1]
}

func (fpa *FinancialPlanA) UpdateLastAccountBalance(f func(AccountBalance) AccountBalance) {
	lastBalanceIndex := len(fpa.Balances) - 1

	if lastBalanceIndex < 0 {
		fpa.AddAccountBalance(f(AccountBalance{}))
	} else {
		fpa.Balances[lastBalanceIndex] = f(fpa.Balances[lastBalanceIndex])
	}
}

func (fpa *FinancialPlanA) AddAccount(account Account) {
	fpa.UpdateLastAccountBalance(func(balance AccountBalance) AccountBalance {
		balance.Accounts = append(balance.Accounts, account)

		return balance
	})
}

func (fpa *FinancialPlanA) RemoveLastAccount() {
	fpa.UpdateLastAccountBalance(func(balance AccountBalance) AccountBalance {
		balance.Accounts = balance.Accounts[:len(balance.Accounts)-1]

		return balance
	})
}

func (fpa *FinancialPlanA) UpdateLastAccount(f func(Account) Account) {
	fpa.UpdateLastAccountBalance(func(balance AccountBalance) AccountBalance {
		lastAccountIndex := len(balance.Accounts) - 1

		if lastAccountIndex < 0 {
			balance.Accounts = append(balance.Accounts, f(Account{}))
		} else {
			balance.Accounts[lastAccountIndex] = f(balance.Accounts[lastAccountIndex])
		}

		return balance
	})
}

func (fpa *FinancialPlanA) AddSubAccount(subAccount SubAccount) {
	fpa.UpdateLastAccount(func(account Account) Account {
		account.Subs = append(account.Subs, subAccount)

		return account
	})
}

func (fpa *FinancialPlanA) UpdateLastSubAccount(f func(SubAccount) SubAccount) {
	fpa.UpdateLastAccount(func(account Account) Account {
		lastSubAccountIndex := len(account.Subs) - 1

		if lastSubAccountIndex < 0 {
			account.Subs = append(account.Subs, f(SubAccount{}))
		} else {
			account.Subs[lastSubAccountIndex] = f(account.Subs[lastSubAccountIndex])
		}

		return account
	})
}

func (fpa *FinancialPlanA) AddUnitAccount(unitAccount UnitAccount) {
	fpa.UpdateLastSubAccount(func(subAccount SubAccount) SubAccount {
		subAccount.Units = append(subAccount.Units, unitAccount)

		return subAccount
	})
}

func (fpa *FinancialPlanA) UpdateLastUnitAccount(f func(UnitAccount) UnitAccount) {
	fpa.UpdateLastSubAccount(func(subAccount SubAccount) SubAccount {
		lastUnitAccountIndex := len(subAccount.Units) - 1

		if lastUnitAccountIndex < 0 {
			subAccount.Units = append(subAccount.Units, f(UnitAccount{}))
		} else {
			subAccount.Units[lastUnitAccountIndex] = f(subAccount.Units[lastUnitAccountIndex])
		}

		return subAccount
	})
}

func (sub SubAccount) HasUnits() bool {
	return len(sub.Units) > 0
}
