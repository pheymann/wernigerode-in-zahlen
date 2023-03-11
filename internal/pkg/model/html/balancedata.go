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

func (b *BalanceData) AddDataPoint(dataPoint DataPoint) {
	if dataPoint.Budget < 0 {
		b.Expenses = append(b.Expenses, dataPoint)
	} else {
		b.Income = append(b.Income, dataPoint)
	}
}
