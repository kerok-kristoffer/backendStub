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
	ingredient2 := createRandomIngredient(t, user.ID)
	formula := CreateRandomFormula(t, user.ID)
	phase := CreateRandomPhase(t, formula.ID)
	formulaIngredient := CreateRandomFormulaIngredient(t, ingredient.ID, phase.ID)

	ingredientParams := models.UpdateFullFormulaIngredientParams{
		FormulaIngredientId:         formulaIngredient.ID,
		IngredientId:                ingredient.ID,
		FormulaIngredientPercentage: F.Int32Between(5, 10),
		FormulaIngredientName:       ingredient.Name,
	}
	ingredientParams2 := models.UpdateFullFormulaIngredientParams{
		IngredientId:                ingredient2.ID,
		FormulaIngredientPercentage: F.Int32Between(5, 10),
		FormulaIngredientName:       ingredient.Name,
	}
	ingredients := []models.UpdateFullFormulaIngredientParams{
		ingredientParams,
		ingredientParams2,
	}

	phaseParams := models.UpdateFullFormulaPhaseParams{
		PhaseName:        F.Genre().Name(),
		PhaseDescription: F.Lorem().Sentence(5),
		PhaseId:          phase.ID,
		Ingredients:      ingredients,
	}
	phases := []models.UpdateFullFormulaPhaseParams{
		phaseParams,
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
	// TODO kerok: Add asserts on all components
	// TODO kerok: need to adapt update to include cases where ingredients are added, removed, changed phase, etc.
	// TODO kerok: Add more test scenarios
	// TODO kerok: delete created test components after tests
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
		Percentage:   F.Int32Between(25, 50),
		PhaseID:      phaseId,
		Description:  sql.NullString{},
	}

	formulaIngredient, err := testQueries.CreateFormulaIngredient(context.Background(), params)
	require.NoError(t, err)

	return formulaIngredient
}
