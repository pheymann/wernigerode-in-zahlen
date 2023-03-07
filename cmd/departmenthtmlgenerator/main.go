package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	htmlEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
	html "wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

var (
	productDirRegex = regexp.MustCompile(`^assets/data/processed/\d+/\d+/\d+/\d+/\d+$`)
)

func main() {
	department := flag.String("department", "", "department to generate a HTML file from")
	departmentName := flag.String("name", "", "department name")

	flag.Parse()

	if *department == "" {
		panic("department is required")
	}

	if *departmentName == "" {
		panic("department name is required")
	}

	var productData = []ProductData{}
	errWalk := filepath.Walk("assets/data/processed/"+*department, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() && productDirRegex.MatchString(path) {
			fmt.Printf("Read %s\n", path)

			financialPlanAFile, err := os.Open(path + "/financial_plan_a.json")
			if err != nil {
				panic(err)
			}
			defer financialPlanAFile.Close()

			var financialPlanBJSONOpt = shared.None[string]()
			financialPlanBFile, err := os.Open(path + "/financial_plan_a.json")
			if err == nil {
				defer financialPlanAFile.Close()

				financialPlanBJSONOpt = shared.Some(io.ReadCompleteFile(financialPlanBFile))
			}

			metadataFile, err := os.Open(path + "/metadata.json")
			if err != nil {
				panic(err)
			}
			defer metadataFile.Close()

			financialPlanA := fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanAFile))
			financialPlanBOpt := shared.Map(financialPlanBJSONOpt, func(financialPlanBJSON string) model.FinancialPlan {
				return fpDecoder.DecodeFromJSON(financialPlanBJSON)
			})
			metadata := metaDecoder.DecodeFromJSON(io.ReadCompleteFile(metadataFile))

			productData = append(productData, ProductData{
				FinancialPlanA:    financialPlanA,
				FinancialPlanBOpt: financialPlanBOpt,
				Metadata:          metadata,
			})
			return nil
		}

		return nil
	})

	if errWalk != nil {
		panic(errWalk)
	}

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

	var cashflowTotal = 0.0
	for _, product := range productData {
		cashflow, incomeTotal, expensesTotal := populateChartData(
			product,
			&expensesProductLinks,
			&chartExpensesDataPerProduct,
			&incomeProductLinks,
			&chartIncomeDataPerProduct,
		)

		cashflowTotal += cashflow
		incomeTotalCashFlow += incomeTotal
		expensesTotalCashFlow += expensesTotal
	}

	year := model.BudgetYear2023
	p := message.NewPrinter(language.German)

	departmentHTML := Department{
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,

		Copy: DepartmentCopy{
			Department:         *departmentName,
			IntroCashflowTotal: fmt.Sprintf("In %s planen wir", year),
			IntroDescription:   encodeIntroDescription(cashflowTotal),

			CashflowTotal:         htmlEncoder.EncodeBudget(cashflowTotal, p),
			IncomeCashflowTotal:   "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal: "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			BackLink: "Zurück zur Übersicht",

			DataDisclosure: `Die Daten auf dieser Webseite beruhen auf dem Haushaltsplan der Statdt Wernigerode aus dem Jahr 2022.
			Da dieser Plan sehr umfangreich ist, muss ich die Daten automatisiert auslesen. Dieser Prozess ist nicht fehlerfrei
			und somit kann ich keine Garantie für die Richtigkeit geben. Schaut zur Kontrolle immer auf das Original, dass ihr
			hier findet: <a href="https://www.wernigerode.de/B%C3%BCrgerservice/Stadtrat/Haushaltsplan/">https://www.wernigerode.de/Bürgerservice/Stadtrat/Haushaltsplan/</a>
			<br><br>
			Die Budgets auf dieser Webseite ergeben sich aus dem Teilfinanzplan A und B und weichen damit vom Haushaltsplan ab, der
			nur Teilfinanzplan A Daten enthält.`,
		},
		CSS: DepartmentCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(cashflowTotal),
		},
	}

	departmentTmpl := template.Must(template.ParseFiles("assets/html/templates/department.template.html"))

	var htmlBytes bytes.Buffer
	if err := departmentTmpl.Execute(&htmlBytes, departmentHTML); err != nil {
		panic(err)
	}

	target := model.TargetFile{
		Path: "assets/html/" + *department + "/",
		Name: "department",
		Tpe:  "html",
	}

	io.WriteFile(target, htmlBytes.String())

}

func populateChartData(
	product ProductData,
	expensesProductLinks *[]string,
	chartExpensesDataPerProduct *html.ChartJSDataset,
	incomeProductLinks *[]string,
	chartIncomeDataPerProduct *html.ChartJSDataset,
) (float64, float64, float64) {
	var productTotalCashflow = 0.0
	var incomeTotalCashFlow = 0.0
	var expensesTotalCashFlow = 0.0

	for _, balance := range product.FinancialPlanA.Balances {
		productTotalCashflow += balance.Budgets[model.BudgetYear2023]
	}

	if product.FinancialPlanBOpt.IsSome {
		for _, balance := range product.FinancialPlanBOpt.Value.Balances {
			productTotalCashflow += balance.Budgets[model.BudgetYear2023]
		}
	}

	if productTotalCashflow < 0 {
		expensesTotalCashFlow += productTotalCashflow
		*expensesProductLinks = append(*expensesProductLinks, fmt.Sprintf(
			"/html/%s/%s/%s/%s/%s/product.html",
			product.Metadata.Department.ID,
			product.Metadata.ProductClass.ID,
			product.Metadata.ProductDomain.ID,
			product.Metadata.ProductGroup.ID,
			product.Metadata.Product.ID,
		))

		chartExpensesDataPerProduct.Labels = append(chartExpensesDataPerProduct.Labels, product.Metadata.Description)
		chartExpensesDataPerProduct.Data = append(chartExpensesDataPerProduct.Data, productTotalCashflow)
	} else {
		incomeTotalCashFlow += productTotalCashflow
		*incomeProductLinks = append(*incomeProductLinks, fmt.Sprintf(
			"/html/%s/%s/%s/%s/%s/product.html",
			product.Metadata.Department.ID,
			product.Metadata.ProductClass.ID,
			product.Metadata.ProductDomain.ID,
			product.Metadata.ProductGroup.ID,
			product.Metadata.Product.ID,
		))

		chartIncomeDataPerProduct.Labels = append(chartIncomeDataPerProduct.Labels, product.Metadata.Description)
		chartIncomeDataPerProduct.Data = append(chartIncomeDataPerProduct.Data, productTotalCashflow)
	}

	return productTotalCashflow, incomeTotalCashFlow, expensesTotalCashFlow
}

func encodeIntroDescription(cashflowTotal float64) string {
	if cashflowTotal < 0 {
		return "für diesen Fachbereich auszugeben."
	}
	return "über diesen Fachbereich einzunehmen."
}

type ProductData struct {
	FinancialPlanA    model.FinancialPlan
	FinancialPlanBOpt shared.Option[model.FinancialPlan]
	Metadata          model.Metadata
}

type Department struct {
	IncomeProductLinks []string
	Income             html.ChartJSDataset

	ExpensesProductLinks []string
	Expenses             html.ChartJSDataset

	Copy DepartmentCopy
	CSS  DepartmentCSS
}

type DepartmentCopy struct {
	Department         string
	IntroCashflowTotal string
	IntroDescription   string

	CashflowTotal         string
	IncomeCashflowTotal   string
	ExpensesCashflowTotal string

	BackLink string

	DataDisclosure template.HTML
}

type DepartmentCSS struct {
	TotalCashflowClass string
}
