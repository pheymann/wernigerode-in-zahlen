set +x
set -e

department=$1
products=$(find assets/data/raw/${department}/* -type f -name '*metadata*' -exec dirname {} \;)

for product in $products; do
  echo "Cleaning up ${product}"
  go run cmd/cleaner/main.go --dir ${product} --metadata
done
