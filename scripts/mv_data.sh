set +x

dir=$1
financial_plan_a=$2
metadata=$3
financial_plan_b=$4

if [ -z "${dir}" ]; then
  echo "set a directory"
  exit 1
fi

if [ -z "${financial_plan_a}" ]; then
  echo "set a financial plan a"
  exit 1
fi

if [ -z "${metadata}" ]; then
  echo "No metadata"
else
  mkdir -p assets/data/raw/${dir}; mv ~/Downloads/tabula-wernigerode_haushaltsplan_2022\ /tabula-wernigerode_haushaltsplan_2022\ -${metadata}.csv $_/metadata.csv
fi

if [ -z "${financial_plan_b}" ]; then
  echo "No financial plan b"
else
  mkdir -p assets/data/raw/${dir}; mv ~/Downloads/tabula-wernigerode_haushaltsplan_2022\ /tabula-wernigerode_haushaltsplan_2022\ -${metadata}.csv $_/financial_plan_b.csv
fi

mkdir -p assets/data/raw/${dir}; mv ~/Downloads/tabula-wernigerode_haushaltsplan_2022\ /tabula-wernigerode_haushaltsplan_2022\ -${financial_plan_a}.csv $_/financial_plan_a.csv
