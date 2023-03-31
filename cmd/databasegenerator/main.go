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
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

var (
	productDirRegex = regexp.MustCompile(
		`assets/data/processed/(?P<department>\d+)/(?P<class>\d+)/(?P<domain>\d+)/(?P<group>\d+)/(?P<product>\d+)(/(?P<sub_product>\d+))?$`,
	)
)

type productFinancialPlan struct {
	FinancialPlan model.FinancialPlan
	SubProducts   map[string]model.FinancialPlan
}

func main() {
	debugRootPath := flag.String("root-path", "", "Debug: root path")

	flag.Parse()

	var database = map[string]map[string]map[string]map[string]map[string]productFinancialPlan{}
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

				var subProductOpt = shared.None[string]()
				if len(matches) > productDirRegex.SubexpIndex("sub_product") {
					subProductOpt = shared.Some(decoder.DecodeString(productDirRegex, "sub_product", matches))

					if subProductOpt.Value == "" {
						subProductOpt = shared.None[string]()
					}
				}

				if _, ok := database[department]; !ok {
					database[department] = map[string]map[string]map[string]map[string]productFinancialPlan{}
				}
				if _, ok := database[department][class]; !ok {
					database[department][class] = map[string]map[string]map[string]productFinancialPlan{}
				}
				if _, ok := database[department][class][domain]; !ok {
					database[department][class][domain] = map[string]map[string]productFinancialPlan{}
				}
				if _, ok := database[department][class][domain][group]; !ok {
					database[department][class][domain][group] = map[string]productFinancialPlan{}
				}
				if _, ok := database[department][class][domain][group][product]; !ok {
					database[department][class][domain][group][product] = productFinancialPlan{
						SubProducts: map[string]model.FinancialPlan{},
					}
				}

				if subProductOpt.IsSome {
					if _, ok := database[department][class][domain][group][product].SubProducts[subProductOpt.Value]; !ok {
						database[department][class][domain][group][product].SubProducts[subProductOpt.Value] = model.FinancialPlan{}
					}

					database[department][class][domain][group][product].SubProducts[subProductOpt.Value] = fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanFile))
				} else {
					productFp := database[department][class][domain][group][product]
					productFp.FinancialPlan = fpDecoder.DecodeFromJSON(io.ReadCompleteFile(financialPlanFile))
					database[department][class][domain][group][product] = productFp
				}
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

func databaseToCSV(database map[string]map[string]map[string]map[string]map[string]productFinancialPlan) [][]string {
	var csvRows = [][]string{
		{
			"department id",
			"product class id",
			"product domain id",
			"product group id",
			"product id",
			"sub product id",
			"account id",
			"account description",
			"sub account id",
			"sub account description",
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
					for product, financialPlan := range products {
						if len(financialPlan.SubProducts) > 0 {
							for subProduct, subProductFinancialPlan := range financialPlan.SubProducts {
								csvRows = append(csvRows, financialPlanToCSV(
									department,
									class,
									domain,
									group,
									product,
									shared.Some(subProduct),
									subProductFinancialPlan,
								)...)
							}
						} else {
							csvRows = append(csvRows, financialPlanToCSV(
								department,
								class,
								domain,
								group,
								product,
								shared.None[string](),
								financialPlan.FinancialPlan,
							)...)
						}
					}
				}
			}
		}
	}

	return csvRows
}

func financialPlanToCSV(
	department string,
	class string,
	domain string,
	group string,
	product string,
	subProduct shared.Option[string],
	financialPlan model.FinancialPlan,
) [][]string {
	csvRows := [][]string{}

	for _, balance := range financialPlan.Balances {
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
						}

						if subProduct.IsSome {
							row = append(row, subProduct.Value)
						} else {
							row = append(row, "")
						}

						row = append(row, []string{
							sub.Id,
							sub.Desc,
							unit.Id,
							unit.Desc,
						}...)

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

					if subProduct.IsSome {
						row = append(row, subProduct.Value)
					} else {
						row = append(row, "")
					}

					row = append(row, []string{
						sub.Id,
						sub.Desc,
					}...)

					for _, budget := range sub.Budgets {
						row = append(row, fmt.Sprintf("%f", budget))
					}

					csvRows = append(csvRows, row)
				}
			}
		}
	}
	return csvRows
}
