package model

type Metadata struct {
	Department    Department
	ProductClass  ProductClass
	ProductDomain ProductDomain
	ProductGroup  ProductGroup
	Product       Product
	Description   string
	Mission       string
	Target        string
	Objectives    string
	Services      string
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

func (metadata *Metadata) Validate() bool {
	return metadata.Department.ID != "" &&
		metadata.ProductClass.ID != "" &&
		metadata.ProductDomain.ID != "" &&
		metadata.ProductGroup.ID != "" &&
		metadata.Product.ID != "" &&
		metadata.Description != "" &&
		metadata.Mission != "" &&
		metadata.Target != "" &&
		metadata.Objectives != "" &&
		metadata.Services != ""
}
