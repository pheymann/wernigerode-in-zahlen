package html

import (
	"strings"

	"wernigerode-in-zahlen.de/internal/pkg/model"
)

type CashflowClass = string

const (
	CashflowClassIncome   CashflowClass = "income"
	CashflowClassExpenses CashflowClass = "expenses"
)

func ClassifyAccount(account model.Account) string {
	if strings.Contains(account.Desc, "Einzahlungen") {
		return CashflowClassIncome
	}
	return CashflowClassExpenses
}
