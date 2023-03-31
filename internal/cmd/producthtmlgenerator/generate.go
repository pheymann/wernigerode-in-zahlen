package htmlgenerator

import (
	"bytes"
	"fmt"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	fpDecoder "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	metaDecoder "wernigerode-in-zahlen.de/internal/pkg/decoder/metadata"
	htmlProductEncoder "wernigerode-in-zahlen.de/internal/pkg/encoder/html/product"
	"wernigerode-in-zahlen.de/internal/pkg/model"
	"wernigerode-in-zahlen.de/internal/pkg/model/html"
	"wernigerode-in-zahlen.de/internal/pkg/shared"
)

func Generate(financialPlanJSON string, metadataJSON string, subProductData []html.ProductData, year model.BudgetYear, debugRootPath string) (string, float64) {
	p := message.NewPrinter(language.German)

	fp := fpDecoder.DecodeFromJSON(financialPlanJSON)
	fpBalanceData, tableData, fpCashflow := readBalanceDataAndCashflow(fp, year)
	metadata := metaDecoder.DecodeFromJSON(metadataJSON)

	sanityCheck(fpCashflow, tableData)

	productTmpl := template.Must(template.ParseFiles(debugRootPath + "assets/html/templates/product.template.html"))

	var htmlBytes bytes.Buffer
	if err := productTmpl.Execute(
		&htmlBytes,
		htmlProductEncoder.Encode(metadata, fpBalanceData, fpCashflow, tableData, subProductData, year, p),
	); err != nil {
		panic(err)
	}

	return htmlBytes.String(), fpCashflow
}

func sanityCheck(cashflowTotal float64, tableData []html.AccountTableData) {
	var tableDataTotal float64
	for _, data := range tableData {
		tableDataTotal += data.CashflowTotal
	}

	if shared.IsUnequal(tableDataTotal, cashflowTotal) {
		panic(fmt.Sprintf("Table data total does not match cashflow total. Expected %f, got %f", cashflowTotal, tableDataTotal))
	}
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
			var accountCashflow = 0.0

			if shared.IsUnequal(account.Budgets[year], 0) {
				for _, sub := range account.Subs {
					var subAccountCashflow = 0.0

					if len(sub.Units) > 0 {
						for _, unit := range sub.Units {
							if shared.IsUnequal(unit.Budgets[year], 0) {
								dataPoint := html.DataPoint{
									Label:  unit.Desc,
									Budget: unit.Budgets[year],
								}

								balanceData[balanceIndex].AddDataPoint(dataPoint)

								tableData = append(tableData, html.AccountTableData{
									Name:          unit.Desc,
									CashflowTotal: unit.Budgets[year],
								})

								subAccountCashflow += unit.Budgets[year]
								accountCashflow += unit.Budgets[year]
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

							subAccountCashflow += sub.Budgets[year]
							accountCashflow += sub.Budgets[year]
						}
					}

					if shared.IsUnequal(subAccountCashflow, sub.Budgets[year]) {
						panic(fmt.Sprintf("Sub-Account %s cashflow does not match budget. Expected %f, got %f", sub.Id, sub.Budgets[year], subAccountCashflow))
					}
				}
			}

			if shared.IsUnequal(accountCashflow, account.Budgets[year]) {
				panic(fmt.Sprintf("Account %s cashflow does not match budget. Expected %f, got %f", account.Id, account.Budgets[year], accountCashflow))
			}
		}

		if len(balanceData[balanceIndex].Expenses) == 0 && len(balanceData[balanceIndex].Income) == 0 {
			balanceData = balanceData[:len(balanceData)-1]
		}
	}

	return balanceData, tableData, cashflowTotal
}
