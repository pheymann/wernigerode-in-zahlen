package departmenthtmlgenerator

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/google/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	htmlEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/model"
	html "wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func GenerateDepartmentHTML(productData []html.ProductData, departmentName string, debugRootPath string) string {
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

	p := message.NewPrinter(language.German)

	var cashflowTotal = 0.0
	var productCopies = []html.DepartmentProductCopy{}
	for _, product := range productData {
		productCopy := html.DepartmentProductCopy{
			Name: product.Metadata.Product.Name,
		}

		cashflow, incomeTotal, expensesTotal := populateChartData(
			product,
			&expensesProductLinks,
			&chartExpensesDataPerProduct,
			&incomeProductLinks,
			&chartIncomeDataPerProduct,
			&productCopy,
			p,
		)

		cashflowTotal += cashflow
		incomeTotalCashFlow += incomeTotal
		expensesTotalCashFlow += expensesTotal
		productCopies = append(productCopies, productCopy)
	}
	sort.Slice(productCopies, func(i, j int) bool {
		return productCopies[i].Name < productCopies[j].Name
	})

	year := model.BudgetYear2023

	departmentHTML := html.Department{
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,

		Copy: html.DepartmentCopy{
			Department:         departmentName,
			IntroCashflowTotal: fmt.Sprintf("In %s planen wir", year),
			IntroDescription:   encodeIntroDescription(cashflowTotal, len(productData)),

			CashflowTotal:         htmlEncoder.EncodeBudget(cashflowTotal, p),
			IncomeCashflowTotal:   "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal: "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			Products: productCopies,

			BackLink: "Zurück zur Übersicht",

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Statdt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>
			<br><br>
			Die Budgets auf dieser Webseite ergeben sich aus dem Teilfinanzplan A und B und weichen damit vom Haushaltsplan ab, der
			nur Teilfinanzplan A Daten enthält.`,
		},
		CSS: html.DepartmentCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(cashflowTotal),
		},
	}

	departmentTmpl := template.Must(template.ParseFiles(debugRootPath + "assets/html/templates/department.template.html"))

	var htmlBytes bytes.Buffer
	if err := departmentTmpl.Execute(&htmlBytes, departmentHTML); err != nil {
		panic(err)
	}

	return htmlBytes.String()
}

func populateChartData(
	product html.ProductData,
	expensesProductLinks *[]string,
	chartExpensesDataPerProduct *html.ChartJSDataset,
	incomeProductLinks *[]string,
	chartIncomeDataPerProduct *html.ChartJSDataset,
	productCopy *html.DepartmentProductCopy,
	p *message.Printer,
) (float64, float64, float64) {
	var productTotalCashflow = 0.0
	var incomeTotalCashFlow = 0.0
	var expensesTotalCashFlow = 0.0

	var financialPlanACashflow = 0.0
	for _, balance := range product.FinancialPlanA.Balances {
		productTotalCashflow += balance.Budgets[model.BudgetYear2023]
		financialPlanACashflow += balance.Budgets[model.BudgetYear2023]
	}
	productCopy.CashflowA = htmlEncoder.EncodeBudget(financialPlanACashflow, p)

	if product.FinancialPlanBOpt.IsSome {
		var financialPlanBCashflow = 0.0
		for _, balance := range product.FinancialPlanBOpt.Value.Balances {
			productTotalCashflow += balance.Budgets[model.BudgetYear2023]
			financialPlanBCashflow += balance.Budgets[model.BudgetYear2023]
		}

		if shared.IsUnequal(financialPlanBCashflow, 0) {
			productCopy.CashflowB = htmlEncoder.EncodeBudget(financialPlanBCashflow, p)
		}
	}

	if productTotalCashflow < 0 {
		expensesTotalCashFlow += productTotalCashflow
		productLink := fmt.Sprintf(
			"/html/%s/%s/%s/%s/%s/product.html",
			product.Metadata.Department.ID,
			product.Metadata.ProductClass.ID,
			product.Metadata.ProductDomain.ID,
			product.Metadata.ProductGroup.ID,
			product.Metadata.Product.ID,
		)

		productCopy.Link = productLink
		*expensesProductLinks = append(*expensesProductLinks, productLink)
		chartExpensesDataPerProduct.Labels = append(chartExpensesDataPerProduct.Labels, product.Metadata.Description)
		chartExpensesDataPerProduct.Data = append(chartExpensesDataPerProduct.Data, productTotalCashflow)
	} else {
		incomeTotalCashFlow += productTotalCashflow
		productLink := fmt.Sprintf(
			"/html/%s/%s/%s/%s/%s/product.html",
			product.Metadata.Department.ID,
			product.Metadata.ProductClass.ID,
			product.Metadata.ProductDomain.ID,
			product.Metadata.ProductGroup.ID,
			product.Metadata.Product.ID,
		)

		productCopy.Link = productLink
		*incomeProductLinks = append(*incomeProductLinks, productLink)
		chartIncomeDataPerProduct.Labels = append(chartIncomeDataPerProduct.Labels, product.Metadata.Description)
		chartIncomeDataPerProduct.Data = append(chartIncomeDataPerProduct.Data, productTotalCashflow)
	}

	return productTotalCashflow, incomeTotalCashFlow, expensesTotalCashFlow
}

func encodeIntroDescription(cashflowTotal float64, numberOfProducts int) template.HTML {
	var earnOrExpese = "einzunehmen"
	if cashflowTotal < 0 {
		earnOrExpese = "auszugeben"
	}

	return template.HTML(fmt.Sprintf(
		"%s. Klick auf eines der <b>%d Produkte</b> in den Diagrammen um mehr zu erfahren.",
		earnOrExpese,
		numberOfProducts,
	))
}
