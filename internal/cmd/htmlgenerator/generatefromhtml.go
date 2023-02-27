package htmlgenerator

import (
	"bufio"
	"html/template"
	"os"

	fpaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func GenerateHTMLForProduct(financialPlanAFile *os.File, metadataFile *os.File) {
	metadata := metaDecoder.DecodeFromJSON(readCompleteFile(metadataFile))
	fpa := fpaDecoder.DecodeFromJSON(readCompleteFile(financialPlanAFile))

	outFile, err := os.Create("test.html")
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	var balanceData = []BalanceChartData{}
	for _, balance := range fpa.Balances {
		for _, account := range balance.Accounts {
			// there are active accounts and not just placeholder
			if account.Budget2022 < -0.001 || account.Budget2022 > 0.001 {
				balanceData = append(balanceData, BalanceChartData{Balance: balance})
				balanceIndex := len(balanceData) - 1

				for _, sub := range account.Subs {
					if len(sub.Units) > 0 {
						for _, unit := range sub.Units {
							dataPoint := ChartDataPoint{
								Budget: unit.Budget2022,
							}
							balanceData[balanceIndex].DataPoints = append(balanceData[balanceIndex].DataPoints, dataPoint)
						}
					} else {
						dataPoint := ChartDataPoint{
							Budget: sub.Budget2022,
						}
						balanceData[balanceIndex].DataPoints = append(balanceData[balanceIndex].DataPoints, dataPoint)
					}
				}
			}
		}
	}

	productHtml := ProductHTML{
		Meta:      metadata,
		ChartData: balanceData,
		Copy: ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",

			CashflowTotal:    "1000EUR",
			CashflowIncome:   "500EUR",
			CashflowExpenses: "500EUR",

			MetaDepartment:    "Fachbereich",
			MetaProductClass:  "Produktklasse",
			MetaProductDomain: "Produktbereich",
			MetaProductGroup:  "Produktgruppe",
			MetaProduct:       "Produkt",
			MetaAccountable:   "Verantwortlich",
			MetaResponsible:   "-----",
		},
		CSS: ProductCSS{
			TotalCashflowClass: "total-cashflow-income",
		},
	}

	productTmpl := template.Must(template.ParseFiles("assets/html/templates/product.template.html"))
	productTmpl.Execute(outFile, productHtml)
}

func readCompleteFile(file *os.File) string {
	scanner := bufio.NewScanner(file)

	var content = ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	return content
}

type ProductHTML struct {
	Meta           model.Metadata
	ChartData      []BalanceChartData
	FinancialPlanA model.FinancialPlanA
	Copy           ProductCopy
	CSS            ProductCSS
}

type BalanceChartData struct {
	Balance    model.AccountBalance
	DataPoints []ChartDataPoint
}

type ChartDataPoint struct {
	Budget float64
}

type ProductCopy struct {
	BackLink string

	CashflowTotal    string
	CashflowIncome   string
	CashflowExpenses string

	MetaDepartment    string
	MetaProductClass  string
	MetaProductDomain string
	MetaProductGroup  string
	MetaProduct       string
	MetaAccountable   string
	MetaResponsible   string
}

type ProductCSS struct {
	TotalCashflowClass string
}
