set +x
set -e

department=$1
products=$(find assets/data/processed/${department}/* -type f -name '*metadata*' -exec dirname {} \;)

for product in $products; do
  echo "Generate HTML for ${product}"
  go run cmd/producthtmlgenerator/main.go --dir ${product}
done
