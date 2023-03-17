package html

import (
	"wernigode-in-zahlen.de/internal/pkg/model"
)

type ProductData struct {
	FinancialPlan model.FinancialPlan
	Metadata      model.Metadata
	CashflowTotal float64
}
