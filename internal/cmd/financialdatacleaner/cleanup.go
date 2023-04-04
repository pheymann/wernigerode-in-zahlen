package financialdatacleaner

import (
	"encoding/csv"
	"log"
	"os"

	fd "wernigerode-in-zahlen.de/internal/pkg/decoder/financialdata"
	fp "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	"wernigerode-in-zahlen.de/internal/pkg/encoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Cleanup(financialDataFile *os.File, productToDepartment map[model.ID]model.ID) string {
	csvReader := csv.NewReader(financialDataFile)
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse financial data CSV", err)
	}

	productAccounts := fd.DecodeAccounts(rows)

	productFinancialPlans := make(map[string]model.FinancialPlanProduct)
	for productID, accounts := range productAccounts {
		productFinancialPlans[productID] = fp.DecodeFromAccounts(accounts)
	}

	departmentFinancialPlans := make(map[string]model.FinancialPlanDepartment)
	for productID, productFinancialPlan := range productFinancialPlans {
		departmentID := productToDepartment[productID]

		if departmentFinancialPlans[departmentID].DepartmentID == "" {
			departmentFinancialPlans[departmentID] = model.FinancialPlanDepartment{
				DepartmentID: departmentID,
				Products:     make(map[model.ID]model.FinancialPlanProduct),
			}
		}

		productFinancialPlan.DepartmentID = departmentID
		departmentFinancialPlans[departmentID].Products[productID] = productFinancialPlan
		for budgetYear, budget := range productFinancialPlan.AdministrationBalance.Budget {
			departmentFinancialPlans[departmentID].AdministrationBalance[budgetYear] += budget
		}
		for budgetYear, budget := range productFinancialPlan.InvestmentsBalance.Budget {
			departmentFinancialPlans[departmentID].InvestmentsBalance[budgetYear] += budget
		}
	}

	cityFinancialPlan := model.FinancialPlanCity{
		AdministrationBalance: make(map[string]float64),
		InvestmentsBalance:    make(map[string]float64),
		Departments:           departmentFinancialPlans,
	}
	for _, departmentFinancialPlan := range departmentFinancialPlans {
		for budgetYear, budget := range departmentFinancialPlan.AdministrationBalance {
			cityFinancialPlan.AdministrationBalance[budgetYear] += budget
		}
		for budgetYear, budget := range departmentFinancialPlan.InvestmentsBalance {
			cityFinancialPlan.InvestmentsBalance[budgetYear] += budget
		}
	}

	return encoder.EncodeToJSON(cityFinancialPlan)
}
