package db

type UpdateFullFormulaIngredientParams struct {
	FormulaIngredientId         int64   `json:"id" binding:"required"`
	IngredientId                int64   `json:"ingredient_id" binding:"required"`
	FormulaIngredientPercentage float32 `json:"percentage" binding:"required"`
	FormulaIngredientName       string  `json:"name" binding:"required"`
	FormulaIngredientCost       float32 `json:"cost" binding:"required"`
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

type FullFormulaIngredient struct {
	Id             int64   `json:"id"`
	IngredientId   int64   `json:"ingredient_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Inci           string  `json:"inci" binding:"required"`
	Percentage     float32 `json:"percentage" binding:"required"`
	WeightInGrams  float32 `json:"weightInGrams" binding:"required"`
	WeightInOunces float32 `json:"weightInOunces" binding:"required"`
	Cost           float32 `json:"cost"  binding:"required"`
	Description    string  `json:"description" binding:"required"`
}

type FullFormulaPhase struct {
	Key                int64                   `json:"key" binding:"required"`
	ID                 int64                   `json:"id" binding:"required"`
	Name               string                  `json:"name" binding:"required"`
	FormulaIngredients []FullFormulaIngredient `json:"ingredients" binding:"required"`
	Description        string                  `json:"description" binding:""`
}

func (ffp FullFormulaPhase) GetKey() int64 {
	return ffp.ID
}

type FormulaResponse struct {
	ID            int64              `json:"id" binding:"required"`
	Phases        []FullFormulaPhase `json:"phases"`
	Name          string             `json:"name"`
	TotalWeight   float64            `json:"totalWeight"`
	TotalWeightOz float64            `json:"totalWeightInOunces"`
	Description   string             `json:"description"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
}
