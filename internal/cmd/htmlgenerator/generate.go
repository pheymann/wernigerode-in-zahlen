package htmlgenerator

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/google/uuid"

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

	var cashflowTotal float64
	var balanceData = []BalanceData{}
	for _, balance := range fpa.Balances {
		cashflowTotal += balance.Budget2022

		balanceData = append(balanceData, BalanceData{Balance: balance})
		balanceIndex := len(balanceData) - 1

		for _, account := range balance.Accounts {
			accountClass := classifyAccount(account)

			if isUnequal(account.Budget2022, 0) {
				for _, sub := range account.Subs {
					if len(sub.Units) > 0 {
						for _, unit := range sub.Units {
							if isUnequal(unit.Budget2022, 0) {
								dataPoint := DataPoint{
									Label:  unit.Desc,
									Budget: unit.Budget2022,
								}

								balanceData[balanceIndex].addDataPoint(dataPoint, accountClass)
							}
						}
					} else {
						if isUnequal(sub.Budget2022, 0) {
							dataPoint := DataPoint{
								Label:  sub.Desc,
								Budget: sub.Budget2022,
							}

							balanceData[balanceIndex].addDataPoint(dataPoint, accountClass)
						}
					}
				}
			}
		}

		if len(balanceData[balanceIndex].Expenses) == 0 || len(balanceData[balanceIndex].Income) == 0 {
			balanceData = balanceData[:len(balanceData)-1]
		}
	}

	productHtml := ProductHTML{
		Meta:            metadata,
		BalanceSections: balanceDataToSections(balanceData),
		Copy: ProductCopy{
			BackLink: "Zurück zur Bereichsübersicht",

			CashflowTotal: printBudget(cashflowTotal),

			MetaDepartment:    "Fachbereich",
			MetaProductClass:  "Produktklasse",
			MetaProductDomain: "Produktbereich",
			MetaProductGroup:  "Produktgruppe",
			MetaProduct:       "Produkt",
			MetaAccountable:   "Verantwortlich",
			MetaResponsible:   "Zuständig",
		},
		CSS: ProductCSS{
			TotalCashflowClass: cssCashflowClass(cashflowTotal),
		},
	}

	productTmpl := template.Must(template.ParseFiles("assets/html/templates/product.template.html"))
	if err := productTmpl.Execute(outFile, productHtml); err != nil {
		panic(err)
	}
}

type AccountClass = string

const (
	AccountClassIncome   AccountClass = "income"
	AccountClassExpenses AccountClass = "expenses"
)

func classifyAccount(account model.Account) string {
	if strings.Contains(account.Desc, "Einzahlungen") {
		return AccountClassIncome
	}
	return AccountClassExpenses
}

func isUnequal(a float64, b float64) bool {
	return a < b-0.001 || a > b+0.001
}

func printBudget(budget float64) string {
	return fmt.Sprintf("%.2f EUR", budget)
}

func cssCashflowClass(budget float64) string {
	if budget < 0 {
		return "total-cashflow-expenses"
	}
	return "total-cashflow-income"
}

func readCompleteFile(file *os.File) string {
	scanner := bufio.NewScanner(file)

	var content = ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	return content
}

type BalanceData struct {
	Balance  model.AccountBalance
	Income   []DataPoint
	Expenses []DataPoint
}

type DataPoint struct {
	Label  string
	Budget float64
}

func (b *BalanceData) addDataPoint(dataPoint DataPoint, class AccountClass) {
	if class == AccountClassIncome {
		b.Income = append(b.Income, dataPoint)
	} else {
		b.Expenses = append(b.Expenses, dataPoint)
	}
}

func balanceDataToSections(data []BalanceData) []BalanceSection {
	var sections = []BalanceSection{}
	for _, balance := range data {
		var incomeCashflowTotal float64
		for _, income := range balance.Income {
			incomeCashflowTotal += income.Budget
		}
		var expensesCashflowTotal float64
		for _, expense := range balance.Expenses {
			expensesCashflowTotal += expense.Budget
		}

		sections = append(sections, BalanceSection{
			ID:                    "balance-" + uuid.New().String(),
			Label:                 balance.Balance.Class,
			CashflowTotal:         printBudget(balance.Balance.Budget2022),
			CSSCashflowTotal:      cssCashflowClass(balance.Balance.Budget2022),
			HasIncomeAndExpenses:  len(balance.Income) > 0 && len(balance.Expenses) > 0,
			HasIncome:             len(balance.Income) > 0,
			IncomeCashflowTotal:   incomeCashflowTotal,
			Income:                dataPointsToChartJSDataset(balance.Income),
			HasExpenses:           len(balance.Expenses) > 0,
			ExpensesCashflowTotal: expensesCashflowTotal,
			Expenses:              dataPointsToChartJSDataset(balance.Expenses),
		})
	}

	return sections
}

func dataPointsToChartJSDataset(dataPoints []DataPoint) ChartJSDataset {
	var labels = []string{}
	var data = []float64{}

	for _, dataPoint := range dataPoints {
		labels = append(labels, dataPoint.Label)
		data = append(data, dataPoint.Budget)
	}

	return ChartJSDataset{
		ID:           "chartjs-" + uuid.New().String(),
		Labels:       labels,
		DatasetLabel: "Budget",
		Data:         data,
	}
}

type ProductHTML struct {
	Meta            model.Metadata
	BalanceSections []BalanceSection
	Copy            ProductCopy
	CSS             ProductCSS
}

type BalanceSection struct {
	ID                    string
	Label                 string
	CashflowTotal         string
	CSSCashflowTotal      string
	HasIncomeAndExpenses  bool
	HasIncome             bool
	IncomeCashflowTotal   float64
	Income                ChartJSDataset
	HasExpenses           bool
	ExpensesCashflowTotal float64
	Expenses              ChartJSDataset
}

type ChartJSDataset struct {
	ID           string
	Labels       []string
	DatasetLabel string
	Data         []float64
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
