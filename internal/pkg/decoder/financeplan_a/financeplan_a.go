package financeplan_a

import "wernigode-in-zahlen.de/internal/pkg/model"

func Decode(financePlans []model.FinancePlanACostCenter) model.FinancePlanA {
	groups, perGroupUnits := separateCostCenterTypes(financePlans)

	return model.FinancePlanA{
		Groups: groups,
		Units:  perGroupUnits,
	}
}

func separateCostCenterTypes(financePlans []model.FinancePlanACostCenter) ([]model.FinancePlanACostCenter, map[string][]model.FinancePlanACostCenter) {
	var groups []model.FinancePlanACostCenter
	perGroupUnits := make(map[string][]model.FinancePlanACostCenter)

	currentCostCenterGroupID := ""
	currentGroupUnits := []model.FinancePlanACostCenter{}

	for _, financePlan := range financePlans {
		if financePlan.Tpe == model.CostCenterGroup {
			perGroupUnits[currentCostCenterGroupID] = currentGroupUnits
			groups = append(groups, financePlan)

			currentCostCenterGroupID = financePlan.Id
			currentGroupUnits = []model.FinancePlanACostCenter{}
		} else {
			currentGroupUnits = append(currentGroupUnits, financePlan)
		}
	}

	return groups, perGroupUnits
}
