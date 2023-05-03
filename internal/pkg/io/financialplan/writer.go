package financialplan_a

import (
	"wernigerode-in-zahlen.de/internal/pkg/io"
	"wernigerode-in-zahlen.de/internal/pkg/model"
)

func Write(FinancialPlan string, target model.TargetFile) {
	target.Tpe = "json"

	io.WriteFile(target, FinancialPlan)
}
