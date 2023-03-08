package html

import (
	"wernigode-in-zahlen.de/internal/pkg/model"
	"wernigode-in-zahlen.de/internal/pkg/shared"
)

type ProductData struct {
	FinancialPlanA    model.FinancialPlan
	FinancialPlanBOpt shared.Option[model.FinancialPlan]
	Metadata          model.Metadata
}
