.PHONY: clean-up-all
clean-up-all:
	./scripts/cleanup_products.sh 1
	go run cmd/cleaner/main.go --dir=assets/data/raw/1 --type=department
	./scripts/merge_financialplans.sh 1

	./scripts/cleanup_products.sh 2
	go run cmd/cleaner/main.go --dir=assets/data/raw/2 --type=department
	./scripts/merge_financialplans.sh 2

	./scripts/cleanup_products.sh 3
	go run cmd/cleaner/main.go --dir=assets/data/raw/3 --type=department
	./scripts/merge_financialplans.sh 3

	./scripts/cleanup_products.sh 4
	go run cmd/cleaner/main.go --dir=assets/data/raw/4 --type=department
	./scripts/merge_financialplans.sh 4

.PHONY: generate-html-all
generate-html-all:
	./scripts/generate_html_products.sh 1
	go run cmd/departmenthtmlgenerator/main.go --department=1 --name="Budget des Bürgermeisters"

	./scripts/generate_html_products.sh 2
	go run cmd/departmenthtmlgenerator/main.go --department=2 --name="Budget Finanzen"

	./scripts/generate_html_products.sh 3
	go run cmd/departmenthtmlgenerator/main.go --department=3 --name="Budget Betriebsbereiche"

	./scripts/generate_html_products.sh 4
	go run cmd/departmenthtmlgenerator/main.go --department=4 --name="Budget Bürgerservice"

	go run cmd/overviewhtmlgenerator/main.go --departments="1,2,3,4"
