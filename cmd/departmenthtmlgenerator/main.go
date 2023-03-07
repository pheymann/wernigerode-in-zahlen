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
	fpaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	htmlEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
	html "wernigode-in-zahlen.de/internal/pkg/model/html"
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

	financialPlanDepartmentFile, err := os.Open("assets/data/processed/" + *department + "/financial_plan_a.json")
	if err != nil {
		panic(err)
	}

	defer financialPlanDepartmentFile.Close()

	finacialPlanDepartment := fpaDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanDepartmentFile))

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

			metadataFile, err := os.Open(path + "/metadata.json")
			if err != nil {
				panic(err)
			}
			defer metadataFile.Close()

			financialPlanA := fpaDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanAFile))
			metadata := metaDecoder.DecodeFromJSON(io.ReadCompleteFile(metadataFile))

			productData = append(productData, ProductData{
				FinancialPlanA: financialPlanA,
				Metadata:       metadata,
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

	for _, product := range productData {
		var productTotalCashflow = 0.0
		for _, balance := range product.FinancialPlanA.Balances {
			productTotalCashflow += balance.Budgets[model.BudgetYear2023]
		}

		if productTotalCashflow < 0 {
			expensesTotalCashFlow += productTotalCashflow
			expensesProductLinks = append(expensesProductLinks, fmt.Sprintf(
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
			incomeProductLinks = append(incomeProductLinks, fmt.Sprintf(
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
	}

	year := model.BudgetYear2023
	p := message.NewPrinter(language.German)

	var cashflowTotal = 0.0
	for _, balance := range finacialPlanDepartment.Balances {
		cashflowTotal += balance.Budgets[year]
	}

	departmentHTML := Department{
		IncomeProductLinks: incomeProductLinks,
		Income:             chartIncomeDataPerProduct,

		ExpensesProductLinks: expensesProductLinks,
		Expenses:             chartExpensesDataPerProduct,

		Copy: DepartmentCopy{
			Department:         *departmentName,
			IntroCashflowTotal: fmt.Sprintf("In %s haben wir", year),
			IntroDescription:   encodeIntroDescription(cashflowTotal),

			CashflowTotal:         htmlEncoder.EncodeBudget(cashflowTotal, p),
			IncomeCashflowTotal:   "Einnahmen: " + htmlEncoder.EncodeBudget(incomeTotalCashFlow, p),
			ExpensesCashflowTotal: "Ausgaben: " + htmlEncoder.EncodeBudget(expensesTotalCashFlow, p),

			BackLink: "Zurück zur Übersicht",
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

func encodeIntroDescription(cashflowTotal float64) string {
	if cashflowTotal < 0 {
		return "für diesen Fachbereich ausgegeben ausgegeben."
	}
	return "über diesen Fachbereich eingenommen."
}

type ProductData struct {
	FinancialPlanA model.FinancialPlan
	Metadata       model.Metadata
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
}

type DepartmentCSS struct {
	TotalCashflowClass string
}
