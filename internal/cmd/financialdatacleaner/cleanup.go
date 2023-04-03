package financialdatacleaner

import (
	"encoding/csv"
	"log"
	"os"

	fd "wernigerode-in-zahlen.de/internal/pkg/decoder/financialdata"
	fp "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Cleanup(financialDataFile *os.File, productToDepartment map[model.ID]model.ID) map[string]string {
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
	for productID, financialPlan := range productFinancialPlans {
		departmentID := productToDepartment[productID]

		if departmentFinancialPlans[departmentID].DepartmentID == "" {
			departmentFinancialPlans[departmentID] = model.FinancialPlanDepartment{
				DepartmentID: departmentID,
			}
		}

		departmentFinancialPlans[departmentID] = append(departmentFinancialPlans[departmentID], financialPlan)
	}

	return map[string]string{}
}
