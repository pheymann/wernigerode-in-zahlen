package html

import "golang.org/x/text/message"

func EncodeBudget(budget float64, p *message.Printer) string {
	return p.Sprintf("%.2f€", budget)
}

func EncodeCSSCashflowClass(budget float64) string {
	if budget < 0 {
		return "total-cashflow-expenses"
	}
	return "total-cashflow-income"
}
