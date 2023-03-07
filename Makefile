.PHONY: clean-up
clean-up:
	./scripts/cleanup_products.sh 1
	go run cmd/cleaner/main.go --dir=assets/data/raw/1 --type=department
