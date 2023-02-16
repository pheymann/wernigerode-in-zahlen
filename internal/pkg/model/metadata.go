package model

type Metadata struct {
	Department    Department
	ProductClass  ProductClass
	ProductDomain ProductDomain
	ProductGroup  ProductGroup
	Product       Product
	Description   string
	FileName      string
	FileType      string
}

type Department struct {
	ID          string
	Name        string
	Accountable string
}

type ProductClass struct {
	ID          string
	Name        string
	Accountable string
}

type ProductDomain struct {
	ID          string
	Name        string
	Responsible string
}

type ProductGroup struct {
	ID   string
	Name string
	Desc string
}

type Product struct {
	ID               string
	Name             string
	LegalRequirement string
}
