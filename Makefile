.PHONY: clean-up-all
clean-up-all:
	make clean-up department=1
	make clean-up department=2
	make clean-up department=3
	make clean-up department=4
	make clean-up department=5

.PHONY: clean-up
clean-up:
	./scripts/cleanup_products.sh $(department)
	go run cmd/cleaner/main.go --dir=assets/data/raw/$(department) --type=department
	./scripts/merge_financialplans.sh $(department)

.PHONY: generate-html-all
generate-html-all:
	make generate-html department=1 name="Budget des Bürgermeisters"
	make generate-html department=2 name="Budget Finanzen"
	make generate-html department=3 name="Budget Betriebsbereiche"
	make generate-html department=4 name="Budget Bürgerservice"
	make generate-html department=5 name="Budget Stadtentwicklung"

	go run cmd/overviewhtmlgenerator/main.go --departments="1,2,3,4,5"

.PHONY: generate-html
generate-html:
	./scripts/generate_html_products.sh $(department)
	go run cmd/departmenthtmlgenerator/main.go --department=$(department) --name="$(name)"
