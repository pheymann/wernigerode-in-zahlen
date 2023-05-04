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

	department1 := financialPlanCity.Departments["1"]
	checkBalance(t, "Department - "+department1.Name, "Administration", department1.AdministrationBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -6_155_000,
			model.BudgetYear2023: -6_220_300,
			model.BudgetYear2024: -6_374_000,
			model.BudgetYear2025: -6_509_300,
			model.BudgetYear2026: -6_343_400,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 3_340_200,
			model.BudgetYear2023: 3_378_700,
			model.BudgetYear2024: 3_421_700,
			model.BudgetYear2025: 3_616_700,
			model.BudgetYear2026: 3_616_700,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -9_495_200,
			model.BudgetYear2023: -9_599_000,
			model.BudgetYear2024: -9_795_700,
			model.BudgetYear2025: -10_126_000,
			model.BudgetYear2026: -9_960_100,
		},
	})

	checkBalance(t, "Department - "+department1.Name, "Investitionen", department1.InvestmentsBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: 718_300,
			model.BudgetYear2023: -407_600,
			model.BudgetYear2024: -165_000,
			model.BudgetYear2025: -160_000,
			model.BudgetYear2026: -160_000,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 3_340_200,
			model.BudgetYear2023: 3_378_700,
			model.BudgetYear2024: 3_421_700,
			model.BudgetYear2025: 3_616_700,
			model.BudgetYear2026: 3_616_700,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -1_051_300,
			model.BudgetYear2023: -1_108_900,
			model.BudgetYear2024: -150_000,
			model.BudgetYear2025: -150_000,
			model.BudgetYear2026: -150_000,
		},
	})

	department2 := financialPlanCity.Departments["2"]
	checkBalance(t, "Department - "+department2.Name, "Administration", department2.AdministrationBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: 33_324_100,
			model.BudgetYear2023: 35_603_300,
			model.BudgetYear2024: 36_528_600,
			model.BudgetYear2025: 37_781_600,
			model.BudgetYear2026: 38_646_400,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 49_764_300,
			model.BudgetYear2023: 52_658_500,
			model.BudgetYear2024: 54_482_100,
			model.BudgetYear2025: 56_264_800,
			model.BudgetYear2026: 57_494_900,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -16_440_200,
			model.BudgetYear2023: -17_055_200,
			model.BudgetYear2024: -17_953_500,
			model.BudgetYear2025: -18_483_200,
			model.BudgetYear2026: -18_848_500,
		},
	})

	checkBalance(t, "Department - "+department2.Name, "Investitionen", department2.InvestmentsBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: 1_819_800,
			model.BudgetYear2023: 1_402_600,
			model.BudgetYear2024: 1_402_600,
			model.BudgetYear2025: 1_402_600,
			model.BudgetYear2026: 1_402_600,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 1_819_800,
			model.BudgetYear2023: 1_402_600,
			model.BudgetYear2024: 1_402_600,
			model.BudgetYear2025: 1_402_600,
			model.BudgetYear2026: 1_402_600,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: 0,
			model.BudgetYear2023: 0,
			model.BudgetYear2024: 0,
			model.BudgetYear2025: 0,
			model.BudgetYear2026: 0,
		},
	})

	// (\d)\.(\d)

	department3 := financialPlanCity.Departments["3"]
	checkBalance(t, "Department - "+department3.Name, "Administration", department3.AdministrationBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -5_510_800,
			model.BudgetYear2023: -5_530_300,
			model.BudgetYear2024: -5_583_100,
			model.BudgetYear2025: -5_731_200,
			model.BudgetYear2026: -5_856_800,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 1_936_400,
			model.BudgetYear2023: 2_053_400,
			model.BudgetYear2024: 2_067_900,
			model.BudgetYear2025: 2_095_600,
			model.BudgetYear2026: 2_105_600,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -7_447_200,
			model.BudgetYear2023: -7_583_700,
			model.BudgetYear2024: -7_651_000,
			model.BudgetYear2025: -7_826_800,
			model.BudgetYear2026: -7_962_400,
		},
	})

	checkBalance(t, "Department - "+department3.Name, "Investitionen", department3.InvestmentsBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -361_000,
			model.BudgetYear2023: -175_600,
			model.BudgetYear2024: -314_200,
			model.BudgetYear2025: -287_300,
			model.BudgetYear2026: -298_300,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 36_000,
			model.BudgetYear2023: 20_000,
			model.BudgetYear2024: 25_000,
			model.BudgetYear2025: 20_000,
			model.BudgetYear2026: 15_000,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -397_000,
			model.BudgetYear2023: -195_600,
			model.BudgetYear2024: -339_200,
			model.BudgetYear2025: -307_300,
			model.BudgetYear2026: -313_300,
		},
	})

	department4 := financialPlanCity.Departments["4"]
	checkBalance(t, "Department - "+department4.Name, "Administration", department4.AdministrationBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -13_783_000,
			model.BudgetYear2023: -16_458_700,
			model.BudgetYear2024: -17_096_500,
			model.BudgetYear2025: -17_588_700,
			model.BudgetYear2026: -18_254_200,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 13_576_000,
			model.BudgetYear2023: 13_469_600,
			model.BudgetYear2024: 13_445_100,
			model.BudgetYear2025: 13_450_700,
			model.BudgetYear2026: 13_420_000,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -27_359_000,
			model.BudgetYear2023: -29_928_300,
			model.BudgetYear2024: -30_541_600,
			model.BudgetYear2025: -31_039_400,
			model.BudgetYear2026: -31_674_200,
		},
	})

	checkBalance(t, "Department - "+department4.Name, "Investitionen", department4.InvestmentsBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -2_641_700,
			model.BudgetYear2023: -2_966_300,
			model.BudgetYear2024: -2_755_200,
			model.BudgetYear2025: -1_701_100,
			model.BudgetYear2026: -1_230_100,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 2_327_200,
			model.BudgetYear2023: 1_072_100,
			model.BudgetYear2024: 550_600,
			model.BudgetYear2025: 129_000,
			model.BudgetYear2026: 237_000,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -4_968_900,
			model.BudgetYear2023: -4_038_400,
			model.BudgetYear2024: -3_305_800,
			model.BudgetYear2025: -1_830_100,
			model.BudgetYear2026: -1_467_100,
		},
	})

	department5 := financialPlanCity.Departments["5"]
	checkBalance(t, "Department - "+department5.Name, "Administration", department5.AdministrationBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: -8_586_000,
			model.BudgetYear2023: -9_173_500,
			model.BudgetYear2024: -9_185_600,
			model.BudgetYear2025: -9_225_600,
			model.BudgetYear2026: -9_223_400,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 1_768_800,
			model.BudgetYear2023: 1_430_700,
			model.BudgetYear2024: 1_256_200,
			model.BudgetYear2025: 1_376_400,
			model.BudgetYear2026: 1_265_400,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -10_354_800,
			model.BudgetYear2023: -10_604_200,
			model.BudgetYear2024: -10_441_800,
			model.BudgetYear2025: -10_602_000,
			model.BudgetYear2026: -10_488_800,
		},
	})

	checkBalance(t, "Department - "+department5.Name, "Investitionen", department5.InvestmentsBalance, model.Cashflow{
		Total: map[model.BudgetYear]float64{
			model.BudgetYear2022: 464_600,
			model.BudgetYear2023: 2_146_900,
			model.BudgetYear2024: -1_587_200,
			model.BudgetYear2025: -1_388_300,
			model.BudgetYear2026: -1_350_200,
		},
		Income: map[model.BudgetYear]float64{
			model.BudgetYear2022: 7_721_400,
			model.BudgetYear2023: 7_048_300,
			model.BudgetYear2024: 6_443_600,
			model.BudgetYear2025: 3_856_000,
			model.BudgetYear2026: 3_752_900,
		},
		Expenses: map[model.BudgetYear]float64{
			model.BudgetYear2022: -7_256_800,
			model.BudgetYear2023: -4_901_400,
			model.BudgetYear2024: -8_030_800,
			model.BudgetYear2025: -5_244_300,
			model.BudgetYear2026: -5_103_100,
		},
	})
}

func checkBalance(t *testing.T, department string, balance string, actual model.Cashflow, expected model.Cashflow) {
	checkTotalCashflow(t, expected, department, balance, "expected")
	checkTotalCashflow(t, actual, department, balance, "actual")
	checkCashflowType(t, expected.Income, actual.Income, "Income", department, balance)
	checkCashflowType(t, expected.Expenses, actual.Expenses, "Expenses", department, balance)
	checkCashflowType(t, expected.Total, actual.Total, "Total", department, balance)
}

func checkTotalCashflow(
	t *testing.T,
	cashflow model.Cashflow,
	department string,
	balance string,
	name string,
) {
	for year, value := range cashflow.Total {
		if shared.IsUnequal(value, cashflow.Income[year]+cashflow.Expenses[year]) {
			t.Errorf("%s.%s.%s: Total cashflow for year %s is not equal to income+expenses: %.2f != %.2f + %.2f", department, balance, name, year, value, cashflow.Income[year], cashflow.Expenses[year])
		}
	}
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
			t.Errorf("%s.%s.%s: year %s with difference %.2f:\n  > expected=%.2f\n  >   actual=%.2f", department, balance, cashflow, year, value-actual[year], value, actual[year])
		}
	}
}
