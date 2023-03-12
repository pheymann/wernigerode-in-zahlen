package html

type ProductTableData struct {
	Name          string
	CashflowTotal float64
	CashflowB     float64
	Link          string
}

type AccountTableData struct {
	Name          string
	CashflowTotal float64
	AboveLimit    string
}
