# Wernigerode in Zahlen
This is the project for [Wernigerode in Zahlen](https://wernigerode-in-zahlen.de). It is a tool to generate a website from my towns yearly financial plans.

## Project structure
* `assets/data/raw`: Raw financial data (CVS) based on the town's financial plan (PDF).
* `assets/data/processed`: Clean up financial data (JSON).
* `docs`: Generated websites.

The remaining projects follows standard Go structures.

## How Tos
### Generate Website

```shell
make
```

### Only clean up raw CVS data

```shell
make clean-up-all
```

### Only generate HTML

```shell
make generate-html-all
```

### Clean up a specific department

Available departments: 1, 2, 3, 4, 5.

```shell
make clean-up department=<ID>
```

### Generate HTML for a specific department

Available departments: 1, 2, 3, 4, 5.

```shell
make generate-html department=<ID> name="<DEPARTMENT NAME>"
```

## Issues in the financial report
For department 3 the summary on page 94 states total expenditures of -5,835,200.00€, but all products for this department combined only add up to -5,804,200.00€. There seems to a be an error in the PDF. An additional 31,000.00€ show up that are not accounted for in any of the products.

The list of balances for all products in that department (manually extrated to verify generated data):
* -2,788,700.00€
* -88,000.00€
* 88,900.00€
* -10,000.00€
* -312,300.00€
* -7,500.00€
* 48,000.00€
* 0.00€
* -51,600.00€
* -60,000.00€
* 249,000.00€
* 0.00€
* 144,900.00€
* -10,000.00€
* -1,606,900.00€
* -19,500.00€
* 34,100.00€
* 0.00€
* -218,900.00€
* -5,800.00€
* -593,600.00€
* -85,000.00€
* -113,300.00€
* -2,500.00€
* -374,000.00€
* -21,500.00€
