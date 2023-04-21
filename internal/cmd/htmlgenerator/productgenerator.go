package htmlgenerator

import (
	"bytes"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	htmlProductEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/product"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func GenerateProducts(plan model.FinancialPlanCity, budgetYear model.BudgetYear, productTmpl *template.Template) []shared.Pair[model.TargetFile, string] {
	var pairs = []shared.Pair[model.TargetFile, string]{}

	for _, department := range plan.Departments {
		for _, product := range department.Products {
			pairs = append(pairs, generateProduct(product, budgetYear, productTmpl))
		}
	}

	return pairs
}

func generateProduct(plan model.FinancialPlanProduct, budgetYear model.BudgetYear, productTmpl *template.Template) shared.Pair[model.TargetFile, string] {
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
