package model

import "fmt"

type CompressedDepartment struct {
	ID               string
	DepartmentName   string
	CashflowTotal    float64
	NumberOfProducts int
}

func (department CompressedDepartment) GetDepartmentLink() string {
	return fmt.Sprintf("/%s/department.html", department.ID)
}
