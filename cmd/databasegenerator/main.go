package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"wernigode-in-zahlen.de/internal/pkg/decoder"
	fpDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan"
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

var (
	productDirRegex = regexp.MustCompile(
		`assets/data/processed/(?P<department>\d+)/(?P<class>\d+)/(?P<domain>\d+)/(?P<group>\d+)/(?P<product>\d+)$`,
	)
)

func main() {
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	var database = map[string]map[string]map[string]map[string]map[string][]model.FinancialPlan{}
	errWalk := filepath.Walk(*debugRootPath+"assets/data/processed/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			matches := productDirRegex.FindStringSubmatch(path)

			if len(matches) != 0 {
				fmt.Printf("Read %s\n", path)

				financialPlanFile, err := os.Open(path + "/merged_financial_plan.json")
				if err != nil {
					panic(err)
				}
				defer financialPlanFile.Close()

				department := decoder.DecodeString(productDirRegex, "department", matches)
				class := decoder.DecodeString(productDirRegex, "class", matches)
				domain := decoder.DecodeString(productDirRegex, "domain", matches)
				group := decoder.DecodeString(productDirRegex, "group", matches)
				product := decoder.DecodeString(productDirRegex, "product", matches)

				if _, ok := database[department]; !ok {
					database[department] = map[string]map[string]map[string]map[string][]model.FinancialPlan{}
				}
				if _, ok := database[department][class]; !ok {
					database[department][class] = map[string]map[string]map[string][]model.FinancialPlan{}
				}
				if _, ok := database[department][class][domain]; !ok {
					database[department][class][domain] = map[string]map[string][]model.FinancialPlan{}
				}
				if _, ok := database[department][class][domain][group]; !ok {
					database[department][class][domain][group] = map[string][]model.FinancialPlan{}
				}
				if _, ok := database[department][class][domain][group][product]; !ok {
					database[department][class][domain][group][product] = []model.FinancialPlan{}
				}

				database[department][class][domain][group][product] = append(
					database[department][class][domain][group][product],
					fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanFile)),
				)
			}

			return nil
		}

		return nil
	})

	if errWalk != nil {
		panic(errWalk)
	}

	jsonDatabase, err := json.MarshalIndent(database, "", "  ")
	if err != nil {
		panic(err)
	}

	io.WriteFile(
		model.TargetFile{
			Path: "assets/data",
			Name: "database",
			Tpe:  "json",
		},
		string(jsonDatabase),
	)

	var csvDatabase bytes.Buffer
	w := csv.NewWriter(&csvDatabase)

	for _, record := range databaseToCSV(database) {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	io.WriteFile(
		model.TargetFile{
			Path: "assets/data",
			Name: "database",
			Tpe:  "csv",
		},
		csvDatabase.String(),
	)
}

func databaseToCSV(database map[string]map[string]map[string]map[string]map[string][]model.FinancialPlan) [][]string {
	var csvRows = [][]string{
		{
			"department id",
			"product class id",
			"product domain id",
			"product group id",
			"product id",
			"account id",
			"account description",
			"sub account id",
			"sub account description",
			"above value limit category",
			"above value limit sub category",
			"budget 2020",
			"budget 2021",
			"budget 2022",
			"budget 2023",
			"budget 2024",
			"budget 2025",
		},
	}

	for department, classes := range database {
		for class, domains := range classes {
			for domain, groups := range domains {
				for group, products := range groups {
					for product, financialPlans := range products {
						for _, plan := range financialPlans {
							for _, balance := range plan.Balances {
								for _, account := range balance.Accounts {
									for _, sub := range account.Subs {
										if len(sub.Units) > 0 {
											for _, unit := range sub.Units {
												var row = []string{
													department,
													class,
													domain,
													group,
													product,
													sub.Id,
													sub.Desc,
													unit.Id,
													unit.Desc,
												}

												if unit.AboveValueLimit != nil {
													row = append(row, unit.AboveValueLimit.Category)
													row = append(row, unit.AboveValueLimit.SubCategory)
												} else {
													row = append(row, "")
													row = append(row, "")
												}

												for _, budget := range unit.Budgets {
													row = append(row, fmt.Sprintf("%f", budget))
												}

												csvRows = append(csvRows, row)
											}
										} else {
											var row = []string{
												department,
												class,
												domain,
												group,
												product,
												sub.Id,
												sub.Desc,
											}

											row = append(row, "")
											row = append(row, "")

											for _, budget := range sub.Budgets {
												row = append(row, fmt.Sprintf("%f", budget))
											}

											csvRows = append(csvRows, row)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return csvRows
}
