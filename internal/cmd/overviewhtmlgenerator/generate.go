package overviewhtmlgenerator

import (
	"bytes"
	"sort"

	"github.com/google/uuid"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
)

func Generate(departments []model.CompressedDepartment, debugRootPath string) string {
	var cashflowTotal = 0.0
	var cashflowAdministration = 0.0
	var cashflowInvestments = 0.0

	var incomeTotalCashFlow = 0.0
	var incomeDepartmentLinks = []string{}
	chartIncomeDataPerDepartment := html.ChartJSDataset{
		ID:           "chartjs-" + uuid.New().String(),
		DatasetLabel: "Einnahmen",
	}

	var expensesTotalCashFlow = 0.0
	var expensesDepartmentLinks = []string{}
	chartExpensesDataPerDepartment := html.ChartJSDataset{
		ID:           "chartjs-" + uuid.New().String(),
		DatasetLabel: "Ausgaben",
	}

	for _, department := range departments {
		cashflowTotal += department.CashflowTotal
		cashflowAdministration += department.CashflowAdministration
		cashflowInvestments += department.CashflowInvestments

		if department.CashflowTotal > 0 {
			incomeTotalCashFlow += department.CashflowTotal
			incomeDepartmentLinks = append(incomeDepartmentLinks, department.GetDepartmentLink())
			chartIncomeDataPerDepartment.Labels = append(chartIncomeDataPerDepartment.Labels, department.DepartmentName)
			chartIncomeDataPerDepartment.Data = append(chartIncomeDataPerDepartment.Data, department.CashflowTotal)
		} else {
			expensesTotalCashFlow += department.CashflowTotal
			expensesDepartmentLinks = append(expensesDepartmentLinks, department.GetDepartmentLink())
			chartExpensesDataPerDepartment.Labels = append(chartExpensesDataPerDepartment.Labels, department.DepartmentName)
			chartExpensesDataPerDepartment.Data = append(chartExpensesDataPerDepartment.Data, department.CashflowTotal)
		}
	}
	sort.Slice(departments, func(i, j int) bool {
		return departments[i].DepartmentName < departments[j].DepartmentName
	})

	var htmlBytes bytes.Buffer

	return htmlBytes.String()
}
