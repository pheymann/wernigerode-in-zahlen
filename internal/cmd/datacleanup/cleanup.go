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
	"wernigerode-in-zahlen.de/internal/pkg/shared"
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

	defer func() {
		if r := recover(); r != nil {
			metadataDecoder.Debug()
			fmt.Printf("\n%+v\n", r)
			os.Exit(1)
		}
	}()

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

	subProductFinancialPlans := make(map[string][]model.FinancialPlanProduct)
	productFinancialPlans := make(map[string]model.FinancialPlanProduct)
	for productID, accounts := range productAccounts {
		plan := fp.DecodeFromAccounts(accounts)
		metaOpt := findMetadata(productID, productToMetadata)

		if metaOpt.IsSome {
			plan.Metadata = metaOpt.Value

			if plan.IsSubProduct() {
				subs := subProductFinancialPlans[plan.Metadata.GetCanonicalProductID()]
				subProductFinancialPlans[plan.Metadata.GetCanonicalProductID()] = append(subs, plan)
			} else {
				productFinancialPlans[productID] = plan
			}
		}
	}

	for productID, subPlans := range subProductFinancialPlans {
		metaOpt := findMetadata(productID, productToMetadata)

		if metaOpt.IsSome {
			plan := *model.NewFinancialPlanProduct()

			plan.ID = productID
			plan.Metadata = metaOpt.Value
			plan.SubProducts = subPlans

			for _, subPlan := range subPlans {
				plan.AdministrationBalance.Cashflow = plan.AdministrationBalance.Cashflow.AddCashflow(subPlan.AdministrationBalance.Cashflow)
				plan.InvestmentsBalance.Cashflow = plan.InvestmentsBalance.Cashflow.AddCashflow(subPlan.InvestmentsBalance.Cashflow)
				plan.Cashflow = plan.Cashflow.AddCashflow(subPlan.Cashflow)
			}

			productFinancialPlans[productID] = plan
		}
	}

	departmentFinancialPlans := make(map[string]model.FinancialPlanDepartment)
	for productID, productFinancialPlan := range productFinancialPlans {
		departmentID := productFinancialPlan.Metadata.Department.ID

		if departmentFinancialPlans[departmentID].ID == "" {
			departmentFinancialPlans[departmentID] = model.FinancialPlanDepartment{
				ID:                    departmentID,
				Name:                  departmentNames[departmentID],
				Products:              make(map[model.ID]model.FinancialPlanProduct),
				AdministrationBalance: model.NewCashFlow(),
				InvestmentsBalance:    model.NewCashFlow(),
				Cashflow:              model.NewCashFlow(),
			}
		}

		department := departmentFinancialPlans[departmentID]
		department.Products[productID] = productFinancialPlan
		department.AdministrationBalance = department.AdministrationBalance.AddCashflow(productFinancialPlan.AdministrationBalance.Cashflow)
		department.InvestmentsBalance = department.InvestmentsBalance.AddCashflow(productFinancialPlan.InvestmentsBalance.Cashflow)
		department.Cashflow = department.Cashflow.AddCashflow(productFinancialPlan.Cashflow)
		departmentFinancialPlans[departmentID] = department
	}

	cityFinancialPlan := model.FinancialPlanCity{
		AdministrationBalance: model.NewCashFlow(),
		InvestmentsBalance:    model.NewCashFlow(),
		Cashflow:              model.NewCashFlow(),
		Departments:           departmentFinancialPlans,
	}
	for _, departmentFinancialPlan := range departmentFinancialPlans {
		cityFinancialPlan.AdministrationBalance = cityFinancialPlan.AdministrationBalance.AddCashflow(departmentFinancialPlan.AdministrationBalance)
		cityFinancialPlan.InvestmentsBalance = cityFinancialPlan.InvestmentsBalance.AddCashflow(departmentFinancialPlan.InvestmentsBalance)
		cityFinancialPlan.Cashflow = cityFinancialPlan.Cashflow.AddCashflow(departmentFinancialPlan.Cashflow)
	}

	return cityFinancialPlan
}

var (
	departmentNames = map[string]string{
		"1": "Budget Oberbürgermeister",
		"2": "Budget Finanzen",
		"3": "Budget Betriebsbereiche",
		"4": "Budget Bürgerservice",
		"5": "Budget Stadtentwicklung",
	}
)

func findMetadata(productID model.ID, productToMetadata map[model.ID]model.Metadata) shared.Option[model.Metadata] {
	if metadata, ok := productToMetadata[productID]; ok {
		return shared.Some(metadata)
	}
	// panic(fmt.Sprintf("No metadata found for product %s", productID))
	fmt.Printf("WARN >> No metadata found for product %s\n", productID)
	return shared.None[model.Metadata]()
}
