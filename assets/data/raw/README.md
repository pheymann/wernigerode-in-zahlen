### Handy Regexes
#### For wernigerode_haushaltsplan_2022.csv

Search:
```regex
(".+")[ ]+([\-\d\.,]+)[ ]+([\-\d\.]+)[ ]+([\-\d\.]+)[ ]+([\-\d\.]+)[ ]+([\-\d\.]+)[ ]+([\-\d\.]+)[ ]+([\-\d\.]+)"
````

Replace:

```regex
$1,"$2",$3,$4,$5,$6,$7,$8
```
