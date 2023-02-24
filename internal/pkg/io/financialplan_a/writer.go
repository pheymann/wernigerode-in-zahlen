package financialplan_a

import (
	"wernigode-in-zahlen.de/internal/pkg/io"
	"wernigode-in-zahlen.de/internal/pkg/model"
)

func Write(FinancialPlanA string, target model.TargetFile) {
	target.Tpe = "json"

	io.WriteFile(target, FinancialPlanA)
}
