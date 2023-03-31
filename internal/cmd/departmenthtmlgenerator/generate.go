package departmenthtmlgenerator

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/google/uuid"
	htmlDepartmentEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/department"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	html "wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Generate(
	financialPlan model.FinancialPlan,
	productData []html.ProductData,
	compressed *model.CompressedDepartment,
	year model.BudgetYear,
	debugRootPath string,
) string {
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

	var depProductData = []html.ProductTableData{}
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

	sanityCheck(financialPlan, compressed, year)

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

func sanityCheck(financialPlan model.FinancialPlan, compressed *model.CompressedDepartment, year model.BudgetYear) {
	var cashflowTotal = 0.0
	var cashflowAdministration = 0.0
	var cashflowInvestmens = 0.0

	for _, balance := range financialPlan.Balances {
		budget := balance.Budgets[year]
		cashflowTotal += budget

		switch balance.Class {
		case model.AccountClassAdministration:
			cashflowAdministration = budget

		case model.AccountClassInvestments:
			cashflowInvestmens = budget
		}
	}

	if shared.IsUnequal(cashflowTotal, compressed.CashflowTotal) {
		fmt.Printf(
			"[WARNING] Compressed and financial plan cashflow divert. Expected %f, got %f\n",
			cashflowTotal,
			compressed.CashflowTotal,
		)
	}
	if shared.IsUnequal(cashflowAdministration, compressed.CashflowAdministration) {
		fmt.Printf(
			"[WARNING] Compressed and financial plan administration cashflow divert. Expected %f, got %f\n",
			cashflowAdministration,
			compressed.CashflowAdministration,
		)
	}
	if shared.IsUnequal(cashflowInvestmens, compressed.CashflowInvestments) {
		fmt.Printf(
			"[WARNING] Compressed and financial plan investments cashflow divert. Expected %f, got %f\n",
			cashflowInvestmens,
			compressed.CashflowInvestments,
		)
	}
}

func populateChartData(
	year model.BudgetYear,
	product html.ProductData,
	compressed *model.CompressedDepartment,

	expensesProductLinks *[]string,
	chartExpensesDataPerProduct *html.ChartJSDataset,

	incomeProductLinks *[]string,
	chartIncomeDataPerProduct *html.ChartJSDataset,
) html.ProductTableData {
	data := html.ProductTableData{}

	var cashflowTotal = 0.0

	for _, balance := range product.FinancialPlan.Balances {
		budget := balance.Budgets[year]
		cashflowTotal += budget

		switch balance.Class {
		case model.AccountClassAdministration:
			compressed.CashflowAdministration += budget
			data.CashflowAdministration = budget

		case model.AccountClassInvestments:
			compressed.CashflowInvestments += budget
			data.CashflowInvestments = budget
		}
	}

	if shared.IsUnequal(cashflowTotal, product.CashflowTotal) {
		panic(fmt.Sprintf("[WARNING] Product and financial plan cashflow divert. Expected %f, got %f\n", cashflowTotal, product.CashflowTotal))
	}

	compressed.CashflowTotal += cashflowTotal
	data.CashflowTotal = cashflowTotal

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
