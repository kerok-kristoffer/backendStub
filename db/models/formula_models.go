package models

type UpdateFullFormulaIngredientParams struct {
	FormulaIngredientId         int64  `json:"id" binding:"required"`
	IngredientId                int64  `json:"ingredient_id" binding:"required"`
	FormulaIngredientPercentage int32  `json:"percentage" binding:"required"`
	FormulaIngredientName       string `json:"name" binding:"required"`
	FormulaIngredientCost       int32  `json:"cost" binding:"required"`
}

type UpdateFullFormulaPhaseParams struct {
	PhaseName        string                              `json:"name" binding:"required"`
	PhaseDescription string                              `json:"description" binding:"required"`
	PhaseId          int64                               `json:"id" binding:"required"`
	Ingredients      []UpdateFullFormulaIngredientParams `json:"ingredients" binding:"required"`
}

type UpdateFullFormulaParams struct {
	FormulaId          int64
	FormulaName        string
	Weight             float32
	WeightOz           float32
	FormulaDescription string
	Phases             []UpdateFullFormulaPhaseParams
	UserId             int64
}
