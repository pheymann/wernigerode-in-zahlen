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
)

func Generate(productData []html.ProductData, compressed *model.CompressedDepartment, year model.BudgetYear, debugRootPath string) string {
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
		data := populateChartData(
			year,
			product,
			compressed,
			&expensesProductLinks,
			&chartExpensesDataPerProduct,
			&incomeProductLinks,
			&chartIncomeDataPerProduct,
		)

		if data.CashflowTotal < 0 {
			expensesTotalCashFlow += data.CashflowTotal
		} else {
			incomeTotalCashFlow += data.CashflowTotal
		}

		data.Name = product.Metadata.Product.Name
		depProductData = append(depProductData, data)
	}
	sort.Slice(depProductData, func(i, j int) bool {
		return depProductData[i].Name < depProductData[j].Name
	})

	compressed.NumberOfProducts = len(productData)

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
	year model.BudgetYear,
	product html.ProductData,
	compressed *model.CompressedDepartment,

	expensesProductLinks *[]string,
	chartExpensesDataPerProduct *html.ChartJSDataset,

	incomeProductLinks *[]string,
	chartIncomeDataPerProduct *html.ChartJSDataset,
) html.DepartmentProductData {
	data := html.DepartmentProductData{}

	var cashflowTotal = 0.0
	var cashflowB = 0.0

	for _, balance := range product.FinancialPlan.Balances {
		cashflowTotal += balance.Budgets[year]

		for _, account := range balance.Accounts {
			for _, sub := range account.Subs {
				for _, unit := range sub.Units {
					if unit.AboveValueLimit != nil {
						cashflowB += unit.Budgets[year]
					}
				}
			}
		}
	}
	compressed.CashflowTotal += cashflowTotal
	compressed.CashflowB += cashflowB
	data.CashflowTotal = cashflowTotal
	data.CashflowB = cashflowB

	productLink := fmt.Sprintf(
		"/%s/%s/%s/%s/%s/product.html",
		product.Metadata.Department.ID,
		product.Metadata.ProductClass.ID,
		product.Metadata.ProductDomain.ID,
		product.Metadata.ProductGroup.ID,
		product.Metadata.Product.ID,
	)

	if cashflowTotal < 0 {
		data.Link = productLink
		*expensesProductLinks = append(*expensesProductLinks, productLink)
		chartExpensesDataPerProduct.Labels = append(chartExpensesDataPerProduct.Labels, product.Metadata.Product.Name)
		chartExpensesDataPerProduct.Data = append(chartExpensesDataPerProduct.Data, cashflowTotal)
	} else {
		data.Link = productLink
		*incomeProductLinks = append(*incomeProductLinks, productLink)
		chartIncomeDataPerProduct.Labels = append(chartIncomeDataPerProduct.Labels, product.Metadata.Product.Name)
		chartIncomeDataPerProduct.Data = append(chartIncomeDataPerProduct.Data, cashflowTotal)
	}

	return data
}
