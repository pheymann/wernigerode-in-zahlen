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
