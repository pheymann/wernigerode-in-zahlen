package htmlgenerator

import (
	"bufio"
	"fmt"
	"html/template"
	"os"

	fpaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/financialplan_a"
	metaDecoder "wernigode-in-zahlen.de/internal/pkg/decoder/metadata"
)

func GenerateHTMLForProduct(financialPlanAFile *os.File, metadataFile *os.File) {
	metadata := metaDecoder.DecodeFromJSON(readCompleteFile(metadataFile))
	fmt.Printf("%+v", fpaDecoder.DecodeFromJSON(readCompleteFile(financialPlanAFile)))

	outFile, err := os.Create("test.html")
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	productTmpl := template.Must(template.ParseFiles("assets/html/templates/product.template.html"))
	productTmpl.Execute(outFile, metadata)
}

func readCompleteFile(file *os.File) string {
	scanner := bufio.NewScanner(file)

	var content = ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	return content
}
