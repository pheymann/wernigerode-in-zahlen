package htmlgenerator

import (
	"bytes"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	htmlProductEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html/product"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func Generate(financialPlanJSON string, metadataJSON string, subProductData []html.ProductData, year model.BudgetYear, debugRootPath string) string {
	p := message.NewPrinter(language.German)

	fp := fpDecoder.DecodeFromJSON(financialPlanJSON)
	fpBalanceData, tableData, fpCashflow := readBalanceDataAndCashflow(fp, year)
	metadata := metaDecoder.DecodeFromJSON(metadataJSON)

	productTmpl := template.Must(template.ParseFiles(debugRootPath + "assets/html/templates/product.template.html"))

	var htmlBytes bytes.Buffer
	if err := productTmpl.Execute(
		&htmlBytes,
		htmlProductEncoder.Encode(metadata, fpBalanceData, fpCashflow, tableData, subProductData, year, p),
	); err != nil {
		panic(err)
	}

	return htmlBytes.String()
}

func readBalanceDataAndCashflow(fp model.FinancialPlan, year model.BudgetYear) ([]html.BalanceData, []html.AccountTableData, float64) {
	var cashflowTotal float64
	var balanceData = []html.BalanceData{}
	var tableData = []html.AccountTableData{}

	for _, balance := range fp.Balances {
		cashflowTotal += balance.Budgets[year]

		balanceData = append(balanceData, html.BalanceData{Balance: balance})
		balanceIndex := len(balanceData) - 1

		for _, account := range balance.Accounts {
			if shared.IsUnequal(account.Budgets[year], 0) {
				for _, sub := range account.Subs {
					if len(sub.Units) > 0 {
						for _, unit := range sub.Units {
							if shared.IsUnequal(unit.Budgets[year], 0) {
								dataPoint := html.DataPoint{
									Label:  unit.Desc,
									Budget: unit.Budgets[year],
								}

								balanceData[balanceIndex].AddDataPoint(dataPoint)

								aboveLimitDesc := ""
								if unit.AboveValueLimit != nil {
									aboveLimitDesc = unit.AboveValueLimit.Category
								}

								tableData = append(tableData, html.AccountTableData{
									Name:          unit.Desc,
									CashflowTotal: unit.Budgets[year],
									AboveLimit:    aboveLimitDesc,
								})
							}
						}
					} else {
						if shared.IsUnequal(sub.Budgets[year], 0) {
							dataPoint := html.DataPoint{
								Label:  sub.Desc,
								Budget: sub.Budgets[year],
							}

							balanceData[balanceIndex].AddDataPoint(dataPoint)

							tableData = append(tableData, html.AccountTableData{
								Name:          sub.Desc,
								CashflowTotal: sub.Budgets[year],
							})
						}
					}
				}
			}
		}

		if len(balanceData[balanceIndex].Expenses) == 0 && len(balanceData[balanceIndex].Income) == 0 {
			balanceData = balanceData[:len(balanceData)-1]
		}
	}

	return balanceData, tableData, cashflowTotal
}
