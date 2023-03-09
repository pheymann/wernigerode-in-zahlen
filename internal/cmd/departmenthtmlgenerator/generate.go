package departmenthtmlgenerator

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/google/uuid"
	htmlDepartmentEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html/department"
	"wernigode-in-zahlen.de/internal/pkg/model"
	html "wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func GenerateDepartmentHTML(productData []html.ProductData, compressed *model.CompressedDepartment, debugRootPath string) string {
	var incomeTotalCashFlow = 0.0
	var incomeProductLinks = []string{}
	chartIncomeDataPerProduct := html.ChartJSDataset{
		ID:           "chartjs-" + uuid.New().String(),
		DatasetLabel: "Einnahmen",
	}

	var expensesTotalCashFlow = 0.0
	var expensesProductLinks = []string{}
	chartExpensesDataPerProduct := html.ChartJSDataset{
		ID:           "chartjs-" + uuid.New().String(),
		DatasetLabel: "Ausgaben",
	}

	var depProductData = []html.DepartmentProductData{}
	for _, product := range productData {
		incomeTotal, expensesTotal, data := populateChartData(
			product,
			compressed,
			&expensesProductLinks,
			&chartExpensesDataPerProduct,
			&incomeProductLinks,
			&chartIncomeDataPerProduct,
		)

		data.Name = product.Metadata.Product.Name
		incomeTotalCashFlow += incomeTotal
		expensesTotalCashFlow += expensesTotal
		depProductData = append(depProductData, data)
	}
	sort.Slice(depProductData, func(i, j int) bool {
		return depProductData[i].Name < depProductData[j].Name
	})

	compressed.NumberOfProducts = len(productData)

	year := model.BudgetYear2023

	departmentTmpl := template.Must(template.ParseFiles(debugRootPath + "assets/html/templates/department.template.html"))

	var htmlBytes bytes.Buffer
	if err := departmentTmpl.Execute(
		&htmlBytes,
		htmlDepartmentEncoder.Encode(
			*compressed,
			year,
			depProductData,

			incomeTotalCashFlow,
			incomeProductLinks,
			chartIncomeDataPerProduct,

			expensesTotalCashFlow,
			expensesProductLinks,
			chartExpensesDataPerProduct,
		),
	); err != nil {
		panic(err)
	}

	return htmlBytes.String()
}

func populateChartData(
	product html.ProductData,
	compressed *model.CompressedDepartment,

	expensesProductLinks *[]string,
	chartExpensesDataPerProduct *html.ChartJSDataset,

	incomeProductLinks *[]string,
	chartIncomeDataPerProduct *html.ChartJSDataset,
) (float64, float64, html.DepartmentProductData) {
	data := html.DepartmentProductData{}

	var productTotalCashflow = 0.0
	var incomeTotalCashFlow = 0.0
	var expensesTotalCashFlow = 0.0

	var cashflowFinancialPlanA = 0.0
	for _, balance := range product.FinancialPlanA.Balances {
		productTotalCashflow += balance.Budgets[model.BudgetYear2023]
		cashflowFinancialPlanA += balance.Budgets[model.BudgetYear2023]
	}
	compressed.CashflowFinancialPlanA += cashflowFinancialPlanA
	data.CashflowFinancialPlanA = cashflowFinancialPlanA

	if product.FinancialPlanBOpt.IsSome {
		var cashflowFinancialPlanB = 0.0
		for _, balance := range product.FinancialPlanBOpt.Value.Balances {
			productTotalCashflow += balance.Budgets[model.BudgetYear2023]
			cashflowFinancialPlanB += balance.Budgets[model.BudgetYear2023]
		}

		if shared.IsUnequal(cashflowFinancialPlanB, 0) {
			compressed.CashflowFinancialPlanB += cashflowFinancialPlanB
			data.CashflowFinancialPlanB = cashflowFinancialPlanB
		}
	}

	if productTotalCashflow < 0 {
		expensesTotalCashFlow += productTotalCashflow
		productLink := fmt.Sprintf(
			"/%s/%s/%s/%s/%s/product.html",
			product.Metadata.Department.ID,
			product.Metadata.ProductClass.ID,
			product.Metadata.ProductDomain.ID,
			product.Metadata.ProductGroup.ID,
			product.Metadata.Product.ID,
		)

		data.Link = productLink
		*expensesProductLinks = append(*expensesProductLinks, productLink)
		chartExpensesDataPerProduct.Labels = append(chartExpensesDataPerProduct.Labels, product.Metadata.Product.Name)
		chartExpensesDataPerProduct.Data = append(chartExpensesDataPerProduct.Data, productTotalCashflow)
	} else {
		incomeTotalCashFlow += productTotalCashflow
		productLink := fmt.Sprintf(
			"/%s/%s/%s/%s/%s/product.html",
			product.Metadata.Department.ID,
			product.Metadata.ProductClass.ID,
			product.Metadata.ProductDomain.ID,
			product.Metadata.ProductGroup.ID,
			product.Metadata.Product.ID,
		)

		data.Link = productLink
		*incomeProductLinks = append(*incomeProductLinks, productLink)
		chartIncomeDataPerProduct.Labels = append(chartIncomeDataPerProduct.Labels, product.Metadata.Product.Name)
		chartIncomeDataPerProduct.Data = append(chartIncomeDataPerProduct.Data, productTotalCashflow)
	}

	compressed.CashflowTotal += productTotalCashflow

	return incomeTotalCashFlow, expensesTotalCashFlow, data
}
