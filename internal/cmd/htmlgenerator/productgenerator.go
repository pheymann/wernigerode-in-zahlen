package htmlgenerator

import (
	"bytes"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	htmlProductEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/product"
	htmlProductWithSubsEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/productwithsubs"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func GenerateProducts(
	plan model.FinancialPlanCity,
	budgetYear model.BudgetYear,
	productTmpl *template.Template,
	productWithSubsTempl *template.Template,
) []shared.Pair[model.TargetFile, string] {
	var pairs = []shared.Pair[model.TargetFile, string]{}

	for _, department := range plan.Departments {
		for _, product := range department.Products {
			if product.SubProducts == nil || len(product.SubProducts) == 0 {
				pairs = append(pairs, generateProduct(product, budgetYear, productTmpl))
			} else {
				pairs = append(pairs, generateProductWithSubs(product, budgetYear, productWithSubsTempl))

				for _, subProduct := range product.SubProducts {
					pairs = append(pairs, generateProduct(subProduct, budgetYear, productTmpl))
				}
			}
		}
	}

	return pairs
}

func generateProduct(
	plan model.FinancialPlanProduct,
	budgetYear model.BudgetYear,
	productTmpl *template.Template,
) shared.Pair[model.TargetFile, string] {
	p := message.NewPrinter(language.German)

	accountTable := generateAccountTable(plan, budgetYear)

	var htmlBytes bytes.Buffer
	if err := productTmpl.Execute(
		&htmlBytes,
		htmlProductEncoder.Encode(plan, accountTable, budgetYear, p),
	); err != nil {
		panic(err)
	}

	content := htmlBytes.String()
	file := model.TargetFile{
		Path: "docs/" + plan.GetPath(),
		Name: "product",
		Tpe:  "html",
	}

	return shared.NewPair(file, content)
}

func generateAccountTable(plan model.FinancialPlanProduct, budgetYear model.BudgetYear) []html.AccountTableData {
	var tableData = []html.AccountTableData{}

	for _, account := range plan.AdministrationBalance.Accounts {
		if shared.IsUnequal(account.Budget[budgetYear], 0) {
			tableData = append(tableData, html.AccountTableData{
				Name:          account.Description,
				CashflowTotal: account.Budget[budgetYear],
			})
		}
	}

	return tableData
}

func generateProductWithSubs(
	plan model.FinancialPlanProduct,
	budgetYear model.BudgetYear,
	productWithSubs *template.Template,
) shared.Pair[model.TargetFile, string] {
	p := message.NewPrinter(language.German)

	var incomeSubProductLinks = []string{}
	chartIncomeDataPerSubProduct := html.ChartJSDataset{
		ID:           "chartjs_sub_products_income",
		DatasetLabel: "Einnahmen",
	}

	var expensesSubProductLinks = []string{}
	chartExpensesDataPerSubProduct := html.ChartJSDataset{
		ID:           "chartjs_sub_products_expenses",
		DatasetLabel: "Ausgaben",
	}

	populateSubProductChartData(
		plan,
		budgetYear,
		&expensesSubProductLinks,
		&chartExpensesDataPerSubProduct,
		&incomeSubProductLinks,
		&chartIncomeDataPerSubProduct,
	)

	var htmlBytes bytes.Buffer
	if err := productWithSubs.Execute(
		&htmlBytes,
		htmlProductWithSubsEncoder.Encode(
			plan,
			budgetYear,
			incomeSubProductLinks,
			chartIncomeDataPerSubProduct,
			expensesSubProductLinks,
			chartExpensesDataPerSubProduct,
			p,
		),
	); err != nil {
		panic(err)
	}

	content := htmlBytes.String()
	file := model.TargetFile{
		Path: "docs/" + plan.GetPath(),
		Name: "product",
		Tpe:  "html",
	}

	return shared.NewPair(file, content)
}

func populateSubProductChartData(
	plan model.FinancialPlanProduct,
	budgetYear model.BudgetYear,

	expensesSubProductLinks *[]string,
	chartExpensesDataPerSubProduct *html.ChartJSDataset,

	incomeSubProductLinks *[]string,
	chartIncomeDataPerSubProduct *html.ChartJSDataset,
) {
	for _, subProduct := range plan.SubProducts {
		if subProduct.Cashflow.Total[budgetYear] < 0 {
			*expensesSubProductLinks = append(*expensesSubProductLinks, subProduct.CreateLink())
			chartExpensesDataPerSubProduct.Labels = append(chartExpensesDataPerSubProduct.Labels, subProduct.Metadata.Product.Name)
			chartExpensesDataPerSubProduct.Data = append(chartExpensesDataPerSubProduct.Data, subProduct.Cashflow.Total[budgetYear])
		} else {
			*incomeSubProductLinks = append(*incomeSubProductLinks, subProduct.CreateLink())
			chartIncomeDataPerSubProduct.Labels = append(chartIncomeDataPerSubProduct.Labels, subProduct.Metadata.Product.Name)
			chartIncomeDataPerSubProduct.Data = append(chartIncomeDataPerSubProduct.Data, subProduct.Cashflow.Total[budgetYear])
		}
	}
}
