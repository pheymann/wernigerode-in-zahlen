all: clean-up generate-html

.PHONY: clean-up
clean-up:
	go run cmd/datacleanup/main.go

.PHONY: generate-html
generate-html:
	go run cmd/htmlgenerator/main.go
