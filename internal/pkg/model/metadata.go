package model

import "fmt"

type Metadata struct {
	Department    Department
	ProductClass  ProductClass
	ProductDomain ProductDomain
	ProductGroup  ProductGroup
	Product       Product
	SubProduct    *SubProduct
	Description   string
	Mission       string
	Target        string
	Objectives    string
	Services      string
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

type SubProduct struct {
	ID   string
	Name string
}

func (metadata Metadata) GetCanonicalProductID() ID {
	return ID(fmt.Sprintf(
		"%s.%s.%s.%s",
		metadata.ProductClass.ID,
		metadata.ProductDomain.ID,
		metadata.ProductGroup.ID,
		metadata.Product.ID,
	))
}

func (metadata *Metadata) Validate() {
	if metadata.Department.ID == "" {
		panic(fmt.Sprintf("metadata.Department.ID is empty. Got: %+v", metadata))
	}
	if metadata.ProductClass.ID == "" {
		panic(fmt.Sprintf("metadata.ProductClass.ID is empty. Got: %+v", metadata))
	}
	if metadata.ProductDomain.ID == "" {
		panic(fmt.Sprintf("metadata.ProductDomain.ID is empty. Got: %+v", metadata))
	}
	if metadata.ProductGroup.ID == "" {
		panic(fmt.Sprintf("metadata.ProductGroup.ID is empty. Got: %+v", metadata))
	}
	if metadata.Product.ID == "" {
		panic(fmt.Sprintf("metadata.Product.ID is empty. Got: %+v", metadata))
	}
	if metadata.Description == "" {
		panic(fmt.Sprintf("metadata.Description is empty. Got: %+v", metadata))
	}
	if metadata.Mission == "" {
		panic(fmt.Sprintf("metadata.Mission is empty. Got: %+v", metadata))
	}
	if metadata.Target == "" {
		panic(fmt.Sprintf("metadata.Target is empty. Got: %+v", metadata))
	}
	if metadata.Objectives == "" {
		fmt.Printf("WARNING: metadata.Objectives is empty. Got: %+v\n", metadata)
	}
	if metadata.Services == "" {
		panic(fmt.Sprintf("metadata.Services is empty. Got: %+v", metadata))
	}
}
