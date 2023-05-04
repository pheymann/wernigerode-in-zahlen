# Wernigerode in Zahlen
This is the project for [Wernigerode in Zahlen](https://wernigerode-in-zahlen.de). It is a tool to generate a website from my city's yearly financial plans.

## Project structure
* `assets/data/raw`: Raw financial data (CVS export provided by the city).
* `assets/data/processed`: Cleaned up financial data (JSON).
* `docs`: Generated website.

The remaining projects follows standard Go structures.

## How Tos
### Generate Website

```shell
make
```

### Only clean up raw CVS data

```shell
make clean-up
```

### Only generate HTML

```shell
make generate-html
```
