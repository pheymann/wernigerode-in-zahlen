package html

type ProductTableData struct {
	Name                   string
	CashflowTotal          float64
	CashflowAdministration float64
	CashflowInvestments    float64
	Link                   string
}

type AccountTableData struct {
	Name          string
	CashflowTotal float64
}
