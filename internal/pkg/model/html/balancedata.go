package html

import "wernigode-in-zahlen.de/internal/pkg/model"

type BalanceData struct {
	Balance  model.AccountBalance
	Income   []DataPoint
	Expenses []DataPoint
}

type DataPoint struct {
	Label  string
	Budget float64
}

func (b *BalanceData) AddDataPoint(dataPoint DataPoint, class CashflowClass) {
	if class == CashflowClassIncome {
		b.Income = append(b.Income, dataPoint)
	} else {
		b.Expenses = append(b.Expenses, dataPoint)
	}
}
