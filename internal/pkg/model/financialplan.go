package model

type FinancialPlan struct {
	Balances []AccountBalance
}

type FinancialPlan2 struct {
	Balances []AccountBalance2
}

type AccountBalance2 struct {
	Class    AccountClass
	Budgets  map[BudgetYear]float64
	Accounts []Account2
}

type Account2 struct {
	ID          string
	ProductID   string
	Description string
	Budget      map[BudgetYear]float64
}

type AccountClass = string

const (
	AccountClassAdministration AccountClass = "admininstration"
	AccountClassInvestments    AccountClass = "balance-investments"
)

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

type AccountBalance struct {
	Id       string
	Class    AccountClass
	Desc     string
	Budgets  map[BudgetYear]float64
	Accounts []Account
}

type Account struct {
	Id      string
	Desc    string
	Budgets map[BudgetYear]float64
	Subs    []SubAccount
}

type SubAccount struct {
	Id      string
	Desc    string
	Budgets map[BudgetYear]float64
	Units   []UnitAccount
}

type UnitAccount struct {
	Id      string
	Desc    string
	Budgets map[BudgetYear]float64
}

func (fpa *FinancialPlan) AddAccountBalance(balance AccountBalance) {
	fpa.Balances = append(fpa.Balances, balance)
}

func (fpa *FinancialPlan) RemoveLastAccountBalance() {
	fpa.Balances = fpa.Balances[:len(fpa.Balances)-1]
}

func (fpa *FinancialPlan) UpdateLastAccountBalance(f func(AccountBalance) AccountBalance) {
	lastBalanceIndex := len(fpa.Balances) - 1

	if lastBalanceIndex < 0 {
		fpa.AddAccountBalance(f(AccountBalance{}))
	} else {
		fpa.Balances[lastBalanceIndex] = f(fpa.Balances[lastBalanceIndex])
	}
}

func (fpa *FinancialPlan) AddAccount(account Account) {
	fpa.UpdateLastAccountBalance(func(balance AccountBalance) AccountBalance {
		balance.Accounts = append(balance.Accounts, account)

		return balance
	})
}

func (fpa *FinancialPlan) RemoveLastAccount() {
	fpa.UpdateLastAccountBalance(func(balance AccountBalance) AccountBalance {
		balance.Accounts = balance.Accounts[:len(balance.Accounts)-1]

		return balance
	})
}

func (fpa *FinancialPlan) UpdateLastAccount(f func(Account) Account) {
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

func (fpa *FinancialPlan) AddSubAccount(subAccount SubAccount) {
	fpa.UpdateLastAccount(func(account Account) Account {
		account.Subs = append(account.Subs, subAccount)

		return account
	})
}

func (fpa *FinancialPlan) UpdateLastSubAccount(f func(SubAccount) SubAccount) {
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

func (fpa *FinancialPlan) AddUnitAccount(unitAccount UnitAccount) {
	fpa.UpdateLastSubAccount(func(subAccount SubAccount) SubAccount {
		subAccount.Units = append(subAccount.Units, unitAccount)

		return subAccount
	})
}

func (fpa *FinancialPlan) UpdateLastUnitAccount(f func(UnitAccount) UnitAccount) {
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
