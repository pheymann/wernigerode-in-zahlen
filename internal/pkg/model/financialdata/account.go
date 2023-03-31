package financialdata

import "wernigerode-in-zahlen.de/internal/pkg/model"

type Account struct {
	ID          string
	ProductID   string
	Description string
	Budget      map[model.BudgetYear]float64
}
