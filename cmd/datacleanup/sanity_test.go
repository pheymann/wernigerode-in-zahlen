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

	checkBalance(t, "City", "Administration", financialPlanCity.AdministrationBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -710_700.00,
			model.BudgetYear2023: -1_779_500.00,
			model.BudgetYear2024: -1_710_600.00,
			model.BudgetYear2025: -1_273_200.00,
			model.BudgetYear2026: -1_031_400.00,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 70_385_700.00,
			model.BudgetYear2023: 72_990_900.00,
			model.BudgetYear2024: 74_673_000.00,
			model.BudgetYear2025: 76_804_200.00,
			model.BudgetYear2026: 77_902_600.00,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -71_096_400.00,
			model.BudgetYear2023: -74_770_400.00,
			model.BudgetYear2024: -76_383_600.00,
			model.BudgetYear2025: -78_077_400.00,
			model.BudgetYear2026: -78_934_000.00,
		},
	})

	checkBalance(t, "City", "Investments", financialPlanCity.InvestmentsBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: 0.00,
			model.BudgetYear2023: 0.00,
			model.BudgetYear2024: -3_419_000.00,
			model.BudgetYear2025: -2_134_100.00,
			model.BudgetYear2026: -1_636_000.00,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 12_955_700.00,
			model.BudgetYear2023: 10_651_900.00,
			model.BudgetYear2024: 8_571_800.00,
			model.BudgetYear2025: 5_557_600.00,
			model.BudgetYear2026: 5_557_500.00,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -12_955_700.00,
			model.BudgetYear2023: -10_651_900.00,
			model.BudgetYear2024: -11_990_800.00,
			model.BudgetYear2025: -7_691_700.00,
			model.BudgetYear2026: -7_193_500.00,
		},
	})
}

func checkBalance(t *testing.T, department string, balance string, actual model.Cashflow, expected model.Cashflow) {
	checkCashflowType(t, expected.Income, actual.Income, "Income", department, balance)
	checkCashflowType(t, expected.Expenses, actual.Expenses, "Expenses", department, balance)
	checkCashflowType(t, expected.Total, actual.Total, "Total", department, balance)
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
			t.Errorf("%s.%s.%s: year %s with difference %.2f: expected=%f, got=%f", department, balance, cashflow, year, value-actual[year], value, actual[year])
		}
	}
}
