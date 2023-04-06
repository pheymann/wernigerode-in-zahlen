package main

import (
	"os"
	"testing"

	decodeFp "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	"wernigerode-in-zahlen.de/internal/pkg/io"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Test_SanityCheck(t *testing.T) {
	financialPlanCityFile, err := os.Open("/Users/paul/Projects/wernigerode-in-zahlen/assets/data/processed/financial_data.json")
	if err != nil {
		t.Fatal(err)
	}
	defer financialPlanCityFile.Close()

	financialPlanCity := decodeFp.DecodeFromJSON2(io.ReadCompleteFile(financialPlanCityFile))

	department1 := financialPlanCity.Departments["1"]
	checkBalance(t, department1.AdministrationBalance, department1.AdministrationBalance, "1", "Administration")
}

func checkBalance(t *testing.T, expected model.Cashflow, actual model.Cashflow, department string, balance string) {
	checkCashflowType(t, expected.Total, actual.Total, "Total", department, balance)
	checkCashflowType(t, expected.Income, actual.Income, "Income", department, balance)
	checkCashflowType(t, expected.Expenses, actual.Expenses, "Expenses", department, balance)
}

func checkCashflowType(
	t *testing.T,
	expected map[model.BudgetYear]float64,
	actual map[model.BudgetYear]float64,
	cashflow string,
	department string,
	balance string,
) {
	for year, value := range expected {
		if shared.IsUnequal(value, actual[year]) {
			t.Errorf("%s.%s.%s: For year %s: expected=%f, got=%f", department, balance, cashflow, year, value, actual[year])
		}
	}
}
