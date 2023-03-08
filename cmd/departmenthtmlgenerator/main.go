package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"

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
	productDirRegex = regexp.MustCompile(`assets/data/processed/\d+/\d+/\d+/\d+/\d+$`)
)

func main() {
	department := flag.String("department", "", "department to generate a HTML file from")
	departmentName := flag.String("name", "", "department name")
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	if *department == "" {
		panic("department is required")
	}

	if *departmentName == "" {
		panic("department name is required")
	}

	var productData = []ProductData{}
	errWalk := filepath.Walk(*debugRootPath+"assets/data/processed/"+*department, func(path string, info os.FileInfo, err error) error {
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
			financialPlanBFile, err := os.Open(path + "/financial_plan_b.json")
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

	p := message.NewPrinter(language.German)

	var cashflowTotal = 0.0
	var productCopies = []DepartmentProductCopy{}
	for _, product := range productData {
		productCopy := DepartmentProductCopy{
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

	departmentHTML := Department{
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,

		Copy: DepartmentCopy{
			Department:         *departmentName,
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
		CSS: DepartmentCSS{
			TotalCashflowClass: htmlEncoder.EncodeCSSCashflowClass(cashflowTotal),
		},
	}

	departmentTmpl := template.Must(template.ParseFiles(*debugRootPath + "assets/html/templates/department.template.html"))

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
	productCopy *DepartmentProductCopy,
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
	IntroDescription   template.HTML

	CashflowTotal         string
	IncomeCashflowTotal   string
	ExpensesCashflowTotal string

	Products []DepartmentProductCopy

	BackLink string

	DataDisclosure template.HTML
}

type DepartmentProductCopy struct {
	Name      string
	Link      string
	CashflowA string
	CashflowB string
}

type DepartmentCSS struct {
	TotalCashflowClass string
}
