package financialplan_a

import (
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Write(FinancialPlan string, target model.TargetFile) {
	target.Tpe = "json"

	io.WriteFile(target, FinancialPlan)
}
