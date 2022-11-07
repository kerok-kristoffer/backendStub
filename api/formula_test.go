package api

import (
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetFormulaApi(t *testing.T) {
	user, err := randomUser()
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// TODO kerok - implement tests for creating, updating and deleting Formulas!
}

func TestCreateFormulaApi(t *testing.T) {

	user, err := randomUser()
	require.NoError(t, err)
	require.NotEmpty(t, user)

	ingredient := db.CreateIngredientParams{
		Name:   f.Food().Vegetable(),
		Hash:   f.Hash().MD5(),
		UserID: user.ID,
	}

	require.NotEmpty(t, ingredient)

	var formulaIngredients []FormulaIngredient
	formulaIngredient := FormulaIngredient{
		Id:             0,
		IngredientId:   2,
		Name:           "",
		Inci:           "",
		Percentage:     1,
		WeightInGrams:  1,
		WeightInOunces: 0.0353,
		Cost:           100,
		Description:    "",
	}
	formulaIngredients = append(formulaIngredients, formulaIngredient)

	var phases []Phase
	phaseA := Phase{
		ID:                 0,
		Name:               "TestCreateFormulaApiPhaseA",
		FormulaIngredients: formulaIngredients,
		Description:        "",
	}
	phases = append(phases, phaseA)

	request := addFormulaRequest{
		Phases:        phases,
		Name:          "TestCreateFormulaApi",
		TotalWeight:   100,
		TotalWeightOz: 3.53,
		Description:   "New description",
	}

	println(request.Name)

}
