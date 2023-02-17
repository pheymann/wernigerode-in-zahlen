.PHONY: mv_product_financial_data
mv_product_financial_data:
	mkdir -p assets/data/raw/$(dir); mv ~/Downloads/tabula-wernigerode_haushaltsplan_2022\ /tabula-wernigerode_haushaltsplan_2022\ -$(metadata).csv $_/metadata.csv
	mkdir -p assets/data/raw/$(dir); mv ~/Downloads/tabula-wernigerode_haushaltsplan_2022\ /tabula-wernigerode_haushaltsplan_2022\ -$(financial_plan_a).csv $_/financial_plan_a.csv
