package htmlgenerator

import (
	"bytes"
	"html/template"
	"sort"

	htmlDepartmentEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/department"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func GenerateDepartments(plan model.FinancialPlanCity, budgetYear model.BudgetYear, departmentTmpl *template.Template) []shared.Pair[model.TargetFile, string] {
	var result = []shared.Pair[model.TargetFile, string]{}

	for _, department := range plan.Departments {
		result = append(result, generateDepartment(department, budgetYear, departmentTmpl))
	}

	return result
}

func generateDepartment(department model.FinancialPlanDepartment, budgetYear model.BudgetYear, departmentTmpl *template.Template) shared.Pair[model.TargetFile, string] {
	var incomeProductLinks = []string{}
	chartIncomeDataPerProduct := html.ChartJSDataset{
		ID:           "chartjs_products_income",
		DatasetLabel: "Einnahmen",
	}

	var expensesProductLinks = []string{}
	chartExpensesDataPerProduct := html.ChartJSDataset{
		ID:           "chartjs_products_expenses",
		DatasetLabel: "Ausgaben",
	}

	populateChartData(department, budgetYear, &expensesProductLinks, &chartExpensesDataPerProduct, &incomeProductLinks, &chartIncomeDataPerProduct)

	productTable := generateProductTable(department, budgetYear)
	sort.Slice(productTable, func(i, j int) bool {
		return productTable[i].Name < productTable[j].Name
	})

	var htmlBytes bytes.Buffer
	if err := departmentTmpl.Execute(
		&htmlBytes,
		htmlDepartmentEncoder.Encode(
			department,
			budgetYear,
			productTable,

			incomeProductLinks,
			chartIncomeDataPerProduct,

			expensesProductLinks,
			chartExpensesDataPerProduct,
		),
	); err != nil {
		panic(err)
	}

	content := htmlBytes.String()
	file := model.TargetFile{
		Path: "docs/" + department.ID + "/",
		Name: "department",
		Tpe:  "html",
	}

	return shared.NewPair(file, content)
}

func generateProductTable(
	department model.FinancialPlanDepartment,
	budgetYear model.BudgetYear,
) []html.ProductTableData {
	table := []html.ProductTableData{}

	for _, product := range department.Products {
		data := html.ProductTableData{
			Name:                   product.Metadata.Product.Name,
			CashflowTotal:          product.Cashflow.Total[budgetYear],
			CashflowAdministration: product.AdministrationBalance.Cashflow.Total[budgetYear],
			CashflowInvestments:    product.InvestmentsBalance.Cashflow.Total[budgetYear],
			Link:                   product.CreateLink(),
		}

		table = append(table, data)
	}
	return table
}

func populateChartData(
	department model.FinancialPlanDepartment,
	budgetYear model.BudgetYear,

	expensesProductLinks *[]string,
	chartExpensesDataPerProduct *html.ChartJSDataset,

	incomeProductLinks *[]string,
	chartIncomeDataPerProduct *html.ChartJSDataset,
) {
	for _, product := range department.Products {
		if product.Cashflow.Total[budgetYear] < 0 {
			*expensesProductLinks = append(*expensesProductLinks, product.CreateLink())
			chartExpensesDataPerProduct.Labels = append(chartExpensesDataPerProduct.Labels, product.Metadata.Product.Name)
			chartExpensesDataPerProduct.Data = append(chartExpensesDataPerProduct.Data, product.Cashflow.Total[budgetYear])
		} else {
			*incomeProductLinks = append(*incomeProductLinks, product.CreateLink())
			chartIncomeDataPerProduct.Labels = append(chartIncomeDataPerProduct.Labels, product.Metadata.Product.Name)
			chartIncomeDataPerProduct.Data = append(chartIncomeDataPerProduct.Data, product.Cashflow.Total[budgetYear])
		}
	}
}
