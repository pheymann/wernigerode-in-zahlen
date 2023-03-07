.PHONY: clean-up
clean-up:
	./scripts/cleanup_products.sh 1
	go run cmd/cleaner/main.go --dir=assets/data/raw/1 --type=department

.PHONY: generate-html
generate-html:
	./scripts/generate_html_products.sh 1
	go run cmd/departmenthtmlgenerator/main.go --department=1 --name="Budget des Bürgermeisters"
