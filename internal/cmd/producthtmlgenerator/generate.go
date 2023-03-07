package htmlgenerator

import (
	"bytes"
	"html/template"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	fpaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
	htmlProductEncoder "wernigode-in-zahlen.de/internal/pkg/encoder/html/product"
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/model/html"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

func GenerateHTMLForProduct(financialPlanAJSON string, financialPlanBJSONOpt shared.Option[string], metadataJSON string, year model.BudgetYear) string {
	fpa := fpaDecoder.DecodeFromJSON(financialPlanAJSON)
	metadata := metaDecoder.DecodeFromJSON(metadataJSON)
	p := message.NewPrinter(language.German)

	fpaBalanceData, fpaCashflow := readBalanceDataAndCashflow(fpa, year)

	fpbBalanceDataOpt := shared.None[[]html.BalanceData]()
	fpbCashflowOpt := shared.None[float64]()
	financialPlanBJSONOpt.ForEach(func(financialPlanBJSON string) {
		fpb := fpaDecoder.DecodeFromJSON(financialPlanBJSON)

		fpbBalanceData, fpbCashflow := readBalanceDataAndCashflow(fpb, year)

		fpbBalanceDataOpt.ToSome(fpbBalanceData)
		fpbCashflowOpt.ToSome(fpbCashflow)
	})

	productTmpl := template.Must(template.ParseFiles("assets/html/templates/product.template.html"))

	var htmlBytes bytes.Buffer
	if err := productTmpl.Execute(
		&htmlBytes,
		htmlProductEncoder.Encode(metadata, fpaBalanceData, fpaCashflow, fpbBalanceDataOpt, fpbCashflowOpt, year, p),
	); err != nil {
		panic(err)
	}

	return htmlBytes.String()
}

func readBalanceDataAndCashflow(fp model.FinancialPlan, year model.BudgetYear) ([]html.BalanceData, float64) {
	var cashflowTotal float64
	var balanceData = []html.BalanceData{}
	for _, balance := range fp.Balances {
		cashflowTotal += balance.Budgets[year]

		balanceData = append(balanceData, html.BalanceData{Balance: balance})
		balanceIndex := len(balanceData) - 1

		for _, account := range balance.Accounts {
			accountClass := html.ClassifyAccount(account)

			if isUnequal(account.Budgets[year], 0) {
				for _, sub := range account.Subs {
					if len(sub.Units) > 0 {
						for _, unit := range sub.Units {
							if isUnequal(unit.Budgets[year], 0) {
								dataPoint := html.DataPoint{
									Label:  unit.Desc,
									Budget: unit.Budgets[year],
								}

								balanceData[balanceIndex].AddDataPoint(dataPoint, accountClass)
							}
						}
					} else {
						if isUnequal(sub.Budgets[year], 0) {
							dataPoint := html.DataPoint{
								Label:  sub.Desc,
								Budget: sub.Budgets[year],
							}

							balanceData[balanceIndex].AddDataPoint(dataPoint, accountClass)
						}
					}
				}
			}
		}

		if len(balanceData[balanceIndex].Expenses) == 0 && len(balanceData[balanceIndex].Income) == 0 {
			balanceData = balanceData[:len(balanceData)-1]
		}
	}

	return balanceData, cashflowTotal
}

func isUnequal(a float64, b float64) bool {
	return a < b-0.001 || a > b+0.001
}
