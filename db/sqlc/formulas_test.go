package db

import (
	"context"
	"database/sql"
	"github.com/kerok-kristoffer/formulating/db/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateFormulaTx(t *testing.T) {
	userAccount := NewUserAccount(testDB)
	user := CreateRandomUser(t)

	ingredient := createRandomIngredient(t, user.ID)
	ingredientToMoveOnUpdate := createRandomIngredient(t, user.ID)
	ingredientToRemoveOnUpdate := createRandomIngredient(t, user.ID)
	ingredientToAddOnUpdate := createRandomIngredient(t, user.ID)
	newPhaseAddedIngredient := createRandomIngredient(t, user.ID)

	formula := CreateRandomFormula(t, user.ID)
	phase := CreateRandomPhase(t, formula.ID)
	phaseToAddOnUpdate := CreateRandomPhase(t, formula.ID)

	formulaIngredient := CreateRandomFormulaIngredient(t, ingredient.ID, phase.ID)
	phaseMoveFormulaIngredient := CreateRandomFormulaIngredient(t, ingredientToMoveOnUpdate.ID, phase.ID)
	CreateRandomFormulaIngredient(t, ingredientToRemoveOnUpdate.ID, phase.ID)

	ingredientParams := ingredientToUpdateParams(formulaIngredient.ID, ingredient)
	addIngredientParams := ingredientToUpdateParams(0, ingredientToAddOnUpdate)
	moveIngredientParams := ingredientToUpdateParams(phaseMoveFormulaIngredient.ID, ingredientToMoveOnUpdate)
	ingredients := []models.UpdateFullFormulaIngredientParams{
		ingredientParams,
		addIngredientParams,
	}

	addPhaseIngredientParams := ingredientToUpdateParams(0, newPhaseAddedIngredient)
	newPhaseIngredients := []models.UpdateFullFormulaIngredientParams{
		addPhaseIngredientParams,
		moveIngredientParams,
	}

	originalPhaseParams := models.UpdateFullFormulaPhaseParams{
		PhaseName:        F.Genre().Name(),
		PhaseDescription: F.Lorem().Sentence(5),
		PhaseId:          phase.ID,
		Ingredients:      ingredients,
	}
	addPhaseParams := models.UpdateFullFormulaPhaseParams{
		PhaseName:        F.Genre().Name(),
		PhaseDescription: F.Lorem().Sentence(7),
		PhaseId:          phaseToAddOnUpdate.ID,
		Ingredients:      newPhaseIngredients,
	}
	phases := []models.UpdateFullFormulaPhaseParams{
		originalPhaseParams,
		addPhaseParams,
	}

	fullFormulaUpdateParams := models.UpdateFullFormulaParams{
		FormulaId:          formula.ID,
		FormulaName:        F.Music().Author(),
		Weight:             50,
		WeightOz:           50 / 28.34,
		FormulaDescription: F.Lorem().Sentence(5),
		Phases:             phases,
		UserId:             user.ID,
	}

	tx, err := userAccount.UpdateFullFormulaTx(context.Background(), fullFormulaUpdateParams)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	compareFormulaProperties(t, fullFormulaUpdateParams, tx)
	for i, phase := range tx.Phases {
		formulaPhaseParams := fullFormulaUpdateParams.Phases[i]
		comparePhaseProperties(t, formulaPhaseParams, phase)
		fullFormulaIngredientParams := formulaPhaseParams.Ingredients
		require.Equal(t, len(fullFormulaIngredientParams), len(phase.Ingredients))
		for j, phaseIngredient := range phase.Ingredients {
			compareIngredientProperties(t, fullFormulaIngredientParams[j], phaseIngredient, phase)
		}
	}
	// TODO kerok: Add more test scenarios

	for _, phaseToRemove := range tx.Phases {
		err = userAccount.DeletePhase(context.Background(), phaseToRemove.Phase.ID)
		require.Error(t, err)

		for _, ingredientToRemove := range phaseToRemove.Ingredients {
			err = userAccount.DeleteFormulaIngredient(context.Background(), ingredientToRemove.ID)
			require.NoError(t, err)
		}
		err = userAccount.DeletePhase(context.Background(), phaseToRemove.Phase.ID)
		require.NoError(t, err)
	}

	err = userAccount.DeleteFormula(context.Background(), formula.ID)
	require.NoError(t, err)

	err = userAccount.DeleteIngredientsByUserId(context.Background(), user.ID)
	require.NoError(t, err)

	err = userAccount.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func compareIngredientProperties(t *testing.T, fullFormulaIngredientParams models.UpdateFullFormulaIngredientParams, phaseIngredient FormulaIngredient, phase UpdatePhaseTxResult) {
	require.Equal(t, fullFormulaIngredientParams.IngredientId, phaseIngredient.IngredientID)
	require.Equal(t, fullFormulaIngredientParams.FormulaIngredientPercentage, phaseIngredient.Percentage)
	require.Equal(t, fullFormulaIngredientParams.FormulaIngredientCost, float32(phaseIngredient.Cost.Float64))
	require.Equal(t, phase.Phase.ID, phaseIngredient.PhaseID)
}

func comparePhaseProperties(t *testing.T, formulaPhaseParams models.UpdateFullFormulaPhaseParams, phase UpdatePhaseTxResult) {
	require.Equal(t, formulaPhaseParams.PhaseName, phase.Phase.Name)
	require.Equal(t, formulaPhaseParams.PhaseId, phase.Phase.ID)
	require.Equal(t, formulaPhaseParams.PhaseDescription, phase.Phase.Description)
}

func compareFormulaProperties(t *testing.T, fullFormulaUpdateParams models.UpdateFullFormulaParams, tx UpdateFormulaTxResult) {
	require.Equal(t, fullFormulaUpdateParams.FormulaId, tx.Formula.ID)
	require.Equal(t, fullFormulaUpdateParams.FormulaName, tx.Formula.Name)
	require.Equal(t, fullFormulaUpdateParams.UserId, tx.Formula.UserID)
	require.Equal(t, fullFormulaUpdateParams.FormulaDescription, tx.Formula.Description)
	require.Equal(t, fullFormulaUpdateParams.Weight, tx.Formula.DefaultAmount)
	require.Equal(t, fullFormulaUpdateParams.WeightOz, tx.Formula.DefaultAmountOz)
	require.Equal(t, len(fullFormulaUpdateParams.Phases), len(tx.Phases))
}

func ingredientToUpdateParams(formulaIngredientId int64, ingredient Ingredient) models.UpdateFullFormulaIngredientParams {
	ingredientParams := models.UpdateFullFormulaIngredientParams{
		FormulaIngredientId:         formulaIngredientId,
		IngredientId:                ingredient.ID,
		FormulaIngredientPercentage: float32(F.RandomFloat(1, 5, 10)),
		FormulaIngredientCost:       float32(F.RandomFloat(2, 50, 1500)),
		FormulaIngredientName:       ingredient.Name,
	}
	return ingredientParams
}

func CreateRandomFormula(t *testing.T, userId int64) Formula {

	amountInGrams := F.Float32(2, 100, 500)

	params := CreateFormulaParams{
		Name:            F.Music().Genre(),
		DefaultAmount:   amountInGrams,
		DefaultAmountOz: amountInGrams / 28.34,
		Description:     F.Lorem().Sentence(10),
		UserID:          userId,
	}

	formula, err := testQueries.CreateFormula(context.Background(), params)
	require.NoError(t, err)

	return formula
}

func CreateRandomPhase(t *testing.T, formulaId int64) Phase {
	params := CreatePhaseParams{
		Name:        F.Food().Vegetable(),
		Description: F.Lorem().Sentence(10),
		FormulaID:   formulaId,
	}

	phase, err := testQueries.CreatePhase(context.Background(), params)
	require.NoError(t, err)

	return phase
}

func CreateRandomFormulaIngredient(t *testing.T, ingredientId int64, phaseId int64) FormulaIngredient {

	params := CreateFormulaIngredientParams{
		IngredientID: ingredientId,
		Percentage:   float32(F.RandomFloat(1, 25, 50)),
		PhaseID:      phaseId,
		Description:  sql.NullString{},
	}

	formulaIngredient, err := testQueries.CreateFormulaIngredient(context.Background(), params)
	require.NoError(t, err)

	return formulaIngredient
}
