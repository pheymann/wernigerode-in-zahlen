package htmlgenerator

import (
	"bytes"
	"html/template"
	"sort"

	htmlOverviewtEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/overview"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
)

func GenerateOverview(plan model.FinancialPlanCity, budgetYear model.BudgetYear, overviewTmpl *template.Template) (model.TargetFile, string) {
	var incomeDepartmentLinks = []string{}
	chartIncomeDataPerDepartment := html.ChartJSDataset{
		ID:           "chartjs_departments_income",
		DatasetLabel: "Einnahmen",
	}

	var expensesDepartmentLinks = []string{}
	chartExpensesDataPerDepartment := html.ChartJSDataset{
		ID:           "chartjs_departments_expenses",
		DatasetLabel: "Ausgaben",
	}

	sortedDepartments := []model.FinancialPlanDepartment{}
	for _, department := range plan.Departments {
		sortedDepartments = append(sortedDepartments, department)
	}
	sort.Slice(sortedDepartments, func(i, j int) bool {
		return sortedDepartments[i].Name < sortedDepartments[j].Name
	})

	for _, department := range sortedDepartments {
		departmentLink := department.CreateLink()

		if department.Cashflow.Total[budgetYear] > 0 {
			incomeDepartmentLinks = append(incomeDepartmentLinks, departmentLink)
			chartIncomeDataPerDepartment.Labels = append(chartIncomeDataPerDepartment.Labels, department.Name)
			chartIncomeDataPerDepartment.Data = append(chartIncomeDataPerDepartment.Data, department.Cashflow.Total[budgetYear])
		} else {
			expensesDepartmentLinks = append(expensesDepartmentLinks, departmentLink)
			chartExpensesDataPerDepartment.Labels = append(chartExpensesDataPerDepartment.Labels, department.Name)
			chartExpensesDataPerDepartment.Data = append(chartExpensesDataPerDepartment.Data, department.Cashflow.Total[budgetYear])
		}
	}

	var htmlBytes bytes.Buffer
	if err := overviewTmpl.Execute(
		&htmlBytes,
		htmlOverviewtEncoder.Encode(
			plan,
			budgetYear,
			incomeDepartmentLinks,
			chartIncomeDataPerDepartment,
			expensesDepartmentLinks,
			chartExpensesDataPerDepartment,
		),
	); err != nil {
		panic(err)
	}

	content := htmlBytes.String()
	file := model.TargetFile{
		Path: "docs/",
		Name: "index",
		Tpe:  "html",
	}

	return file, content
}
