set +x
set -e

department=$1
products=$(find assets/data/processed/${department}/* -type f -name '*metadata*' -exec dirname {} \;)

for product in $products; do
  echo "Merging financial plans for ${product}"
  go run cmd/financialplanmerger/main.go --dir ${product}
done
