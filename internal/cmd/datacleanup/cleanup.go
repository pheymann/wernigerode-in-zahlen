package datacleanup

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	fd "wernigerode-in-zahlen.de/internal/pkg/decoder/financialdata"
	fp "wernigerode-in-zahlen.de/internal/pkg/decoder/financialplan"
	decodeMeta "wernigerode-in-zahlen.de/internal/pkg/decoder/metadata"
	"wernigerode-in-zahlen.de/internal/pkg/encoder"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Cleanup(financialDataFile *os.File, metadataFiles []*os.File) string {
	metadata := cleanupAllMetadata(metadataFiles)
	productToMetadata := createProductToMetadataMapping(metadata)
	cityFinancialPlan := cleanupFinancialPlans(financialDataFile, productToMetadata)

	return encoder.EncodeToJSON(cityFinancialPlan)
}

func cleanupAllMetadata(metadataFiles []*os.File) []model.Metadata {
	var metadata = []model.Metadata{}

	for _, metadataFile := range metadataFiles {
		metadata = append(metadata, cleanupMetadata(metadataFile))
	}

	return metadata
}

func cleanupMetadata(metadataFile *os.File) model.Metadata {
	metadataScanner := bufio.NewScanner(metadataFile)
	metadataLines := []string{}

	for metadataScanner.Scan() {
		metadataLines = append(metadataLines, metadataScanner.Text())
	}

	metadataDecoder := decodeMeta.NewMetadataDecoder()

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		metadataDecoder.Debug()
	// 		fmt.Printf("\n%+v\n", r)
	// 		os.Exit(1)
	// 	}
	// }()

	metadata := metadataDecoder.DecodeFromCSV(metadataLines)

	return metadata
}

func createProductToMetadataMapping(metadata []model.Metadata) map[model.ID]model.Metadata {
	productToMetadata := make(map[model.ID]model.Metadata)

	for _, m := range metadata {
		var key = fmt.Sprintf("%s.%s.%s.%s", m.ProductClass.ID, m.ProductDomain.ID, m.ProductGroup.ID, m.Product.ID)

		if m.SubProduct != nil {
			key = fmt.Sprintf("%s.%s", key, m.SubProduct.ID)
		}
		productToMetadata[key] = m
	}

	return productToMetadata
}

func cleanupFinancialPlans(financialDataFile *os.File, productToMetadata map[model.ID]model.Metadata) model.FinancialPlanCity {
	csvReader := csv.NewReader(financialDataFile)
	csvReader.Comma = ';'
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse financial data CSV", err)
	}

	productAccounts := fd.DecodeAccounts(rows)

	productFinancialPlans := make(map[string]model.FinancialPlanProduct)
	for productID, accounts := range productAccounts {
		plan := fp.DecodeFromAccounts(accounts)
		plan.Metadata = findMetadata(productID, productToMetadata)
		productFinancialPlans[productID] = plan
	}

	departmentFinancialPlans := make(map[string]model.FinancialPlanDepartment)
	for productID, productFinancialPlan := range productFinancialPlans {
		departmentID := productFinancialPlan.Metadata.Department.ID

		if departmentFinancialPlans[departmentID].DepartmentID == "" {
			departmentFinancialPlans[departmentID] = model.FinancialPlanDepartment{
				DepartmentID: departmentID,
				Products:     make(map[model.ID]model.FinancialPlanProduct),
			}
		}

		departmentFinancialPlans[departmentID].Products[productID] = productFinancialPlan
		for budgetYear, budget := range productFinancialPlan.AdministrationBalance.Budget {
			departmentFinancialPlans[departmentID].AdministrationBalance[budgetYear] += budget
		}
		for budgetYear, budget := range productFinancialPlan.InvestmentsBalance.Budget {
			departmentFinancialPlans[departmentID].InvestmentsBalance[budgetYear] += budget
		}
	}

	cityFinancialPlan := model.FinancialPlanCity{
		AdministrationBalance: make(map[string]float64),
		InvestmentsBalance:    make(map[string]float64),
		Departments:           departmentFinancialPlans,
	}
	for _, departmentFinancialPlan := range departmentFinancialPlans {
		for budgetYear, budget := range departmentFinancialPlan.AdministrationBalance {
			cityFinancialPlan.AdministrationBalance[budgetYear] += budget
		}
		for budgetYear, budget := range departmentFinancialPlan.InvestmentsBalance {
			cityFinancialPlan.InvestmentsBalance[budgetYear] += budget
		}
	}

	return cityFinancialPlan
}

func findMetadata(productID model.ID, productToMetadata map[model.ID]model.Metadata) model.Metadata {
	if metadata, ok := productToMetadata[productID]; ok {
		return metadata
	}
	// panic(fmt.Sprintf("No metadata found for product %s", productID))
	fmt.Printf("WARN >> No metadata found for product %s\n", productID)
	return model.Metadata{}
}
