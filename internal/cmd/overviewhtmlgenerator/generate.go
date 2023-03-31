package overviewhtmlgenerator

import (
	"bytes"
	"html/template"
	"sort"

	"github.com/google/uuid"
	htmlOverviewtEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html/overview"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
)

func Generate(departments []model.CompressedDepartment, debugRootPath string) string {
	var cashflowTotal = 0.0

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

	year := model.BudgetYear2023

	overviewTmpl := template.Must(template.ParseFiles(debugRootPath + "assets/html/templates/overview.template.html"))

	var htmlBytes bytes.Buffer
	if err := overviewTmpl.Execute(
		&htmlBytes,
		htmlOverviewtEncoder.Encode(
			departments,
			year,

			cashflowTotal,

			incomeTotalCashFlow,
			incomeDepartmentLinks,
			chartIncomeDataPerDepartment,

			expensesTotalCashFlow,
			expensesDepartmentLinks,
			chartExpensesDataPerDepartment,
		),
	); err != nil {
		panic(err)
	}

	return htmlBytes.String()
}
