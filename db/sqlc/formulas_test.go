package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestUpdateFormulaFailRollbackOnFaultyFormulaID(t *testing.T) {
	userAccount := NewUserAccount(testDB)
	user := CreateRandomUser(t)

	ingredientCount := 6
	originalIngredients := createIngredients(t, user.ID, ingredientCount)

	formula := CreateRandomFormula(t, user.ID)

	phaseCount := 2
	originalPhases := createPhases(t, formula.ID, phaseCount)
	randomlyAssignIngredientsToPhases(t, originalIngredients, originalPhases)

	originalFormula, err := userAccount.GetFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, originalFormula)
	log.Println(originalFormula)

	updateIngredients := createIngredients(t, user.ID, 4)
	var updateFullFormulaIngredientParams = make([]UpdateFullFormulaIngredientParams, 0)
	for _, ingredient := range updateIngredients {
		updateFullFormulaIngredientParams = append(updateFullFormulaIngredientParams, ingredientToUpdateParams(0, ingredient))
	}
	updatePhases := generatePhaseParamsForUpdate(updateFullFormulaIngredientParams)
	faultyFormulaId := -1
	expectedFormulaParams := UpdateFullFormulaParams{
		FormulaId:          int64(faultyFormulaId),
		FormulaName:        F.Music().Author(),
		Weight:             25,
		WeightOz:           25 / 28.34,
		FormulaDescription: F.Lorem().Sentence(6),
		Phases:             updatePhases,
		UserId:             user.ID,
	}

	txResultFormula, err := userAccount.UpdateFullFormulaTx(context.Background(), expectedFormulaParams)

	require.Error(t, err)
	require.Empty(t, txResultFormula)
	// TODO assert that the updated Formula corresponds to the diff, not original

	updatedFormula, err := userAccount.GetFullFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFormula)
	actualFormula, err := userAccount.GetFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actualFormula)

	require.Equal(t, originalFormula, actualFormula)

	for _, phaseToRemove := range originalPhases {
		err = userAccount.DeletePhase(context.Background(), phaseToRemove.ID)
		require.Error(t, err)
	}

	for _, row := range updatedFormula {
		err := userAccount.DeleteFormulaIngredient(context.Background(), row.FormulaIngredientID)
		require.NoError(t, err)
	}

	for _, phaseToRemove := range originalPhases {

		err = userAccount.DeletePhase(context.Background(), phaseToRemove.ID)
		require.NoError(t, err)
	}

	err = userAccount.DeleteFormula(context.Background(), formula.ID)
	require.NoError(t, err)

	err = userAccount.DeleteIngredientsByUserId(context.Background(), user.ID)
	require.NoError(t, err)

	err = userAccount.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}
func TestUpdateFormulaFailRollbackOnFaultyIngredientId(t *testing.T) {
	userAccount := NewUserAccount(testDB)
	user := CreateRandomUser(t)

	ingredientCount := 6
	originalIngredients := createIngredients(t, user.ID, ingredientCount)

	formula := CreateRandomFormula(t, user.ID)

	phaseCount := 2
	originalPhases := createPhases(t, formula.ID, phaseCount)
	randomlyAssignIngredientsToPhases(t, originalIngredients, originalPhases)

	originalFormula, err := userAccount.GetFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, originalFormula)
	log.Println(originalFormula)

	updateIngredients := createIngredients(t, user.ID, 4)
	//updatePhases := createPhases(t, formula.ID, phaseCount)
	var updateFullFormulaIngredientParams = make([]UpdateFullFormulaIngredientParams, 0)
	for _, ingredient := range updateIngredients {
		updateFullFormulaIngredientParams = append(updateFullFormulaIngredientParams, ingredientToUpdateParams(0, ingredient))
	}
	updatePhases := generatePhaseParamsForUpdate(updateFullFormulaIngredientParams)
	updatePhases[0].Ingredients[0].FormulaIngredientId = -1
	expectedFormulaParams := UpdateFullFormulaParams{
		FormulaId:          formula.ID,
		FormulaName:        F.Music().Author(),
		Weight:             25,
		WeightOz:           25 / 28.34,
		FormulaDescription: F.Lorem().Sentence(6),
		Phases:             updatePhases,
		UserId:             user.ID,
	}

	txResultFormula, err := userAccount.UpdateFullFormulaTx(context.Background(), expectedFormulaParams)

	require.Error(t, err)
	require.Empty(t, txResultFormula)

	updatedFormula, err := userAccount.GetFullFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFormula)
	actualFormula, err := userAccount.GetFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actualFormula)

	require.Equal(t, originalFormula, actualFormula)

	for _, phaseToRemove := range originalPhases {
		err = userAccount.DeletePhase(context.Background(), phaseToRemove.ID)
		require.Error(t, err)
	}

	for _, row := range updatedFormula {
		err := userAccount.DeleteFormulaIngredient(context.Background(), row.FormulaIngredientID)
		require.NoError(t, err)
	}

	for _, phaseToRemove := range originalPhases {

		err = userAccount.DeletePhase(context.Background(), phaseToRemove.ID)
		require.NoError(t, err)
	}

	err = userAccount.DeleteFormula(context.Background(), formula.ID)
	require.NoError(t, err)

	err = userAccount.DeleteIngredientsByUserId(context.Background(), user.ID)
	require.NoError(t, err)

	err = userAccount.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func TestUpdateFormulaWithRandomChangesSuccess(t *testing.T) {
	userAccount := NewUserAccount(testDB)
	user := CreateRandomUser(t)

	ingredientCount := 15
	originalIngredients := createIngredients(t, user.ID, ingredientCount)

	formula := CreateRandomFormula(t, user.ID)

	phaseCount := 3
	originalPhases := createPhases(t, formula.ID, phaseCount)
	randomlyAssignIngredientsToPhases(t, originalIngredients, originalPhases)

	originalFormula, err := userAccount.GetFormula(context.Background(), formula.ID)
	require.NoError(t, err)
	require.NotEmpty(t, originalFormula)
	log.Println(originalFormula)

	updateIngredients := createIngredients(t, user.ID, 4)
	var updateFullFormulaIngredientParams = make([]UpdateFullFormulaIngredientParams, 0)
	for _, ingredient := range updateIngredients {
		updateFullFormulaIngredientParams = append(updateFullFormulaIngredientParams, ingredientToUpdateParams(0, ingredient))
	}

	updatePhases := generatePhaseParamsForUpdate(updateFullFormulaIngredientParams)
	expectedFormulaParams := UpdateFullFormulaParams{
		FormulaId:          formula.ID,
		FormulaName:        F.Music().Author(),
		Weight:             25,
		WeightOz:           25 / 28.34,
		FormulaDescription: F.Lorem().Sentence(6),
		Phases:             updatePhases,
		//Phases: nil,
		UserId: user.ID,
	}
	txResultFormula, err := userAccount.UpdateFullFormulaTx(context.Background(), expectedFormulaParams)
	require.NoError(t, err)
	require.NotEmpty(t, txResultFormula)

	updatedFormula, err := userAccount.GetFullFormula(context.Background(), expectedFormulaParams.FormulaId)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFormula)
	actualFormula, err := userAccount.GetFormula(context.Background(), expectedFormulaParams.FormulaId)
	require.NoError(t, err)
	require.NotEmpty(t, actualFormula)

	updatedPhases := GenerateFormulaViewModelPhases(updatedFormula)

	log.Println(updatedPhases)
	for i, phase := range updatedPhases {
		log.Println(phase.Name, " should equal: ", expectedFormulaParams.Phases[i].PhaseName)
	}

	formulaResponse := FormulaResponse{
		ID:            actualFormula.ID,
		Phases:        updatedPhases,
		Name:          actualFormula.Name,
		TotalWeight:   float64(actualFormula.DefaultAmount),
		TotalWeightOz: float64(actualFormula.DefaultAmountOz),
		Description:   actualFormula.Description,
		CreatedAt:     actualFormula.CreatedAt.Format("06-01-02"),
		UpdatedAt:     actualFormula.UpdatedAt.Format("06-01-02"),
	}

	compareFormulaProperties(t, expectedFormulaParams, actualFormula)
	validatePhasesEquals(t, expectedFormulaParams, formulaResponse)

	for _, phaseToRemove := range txResultFormula.Phases {
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

func generatePhaseParamsForUpdate(ingredientParams []UpdateFullFormulaIngredientParams) []UpdateFullFormulaPhaseParams {
	phaseCount := 3
	phases := make([]UpdateFullFormulaPhaseParams, 0)

	for i := 0; i < phaseCount; i++ {
		phases = append(phases, UpdateFullFormulaPhaseParams{
			PhaseName:        F.Food().Vegetable(),
			PhaseDescription: F.Lorem().Sentence(8),
			PhaseId:          0,
			Ingredients:      ingredientParams[i : i+2],
		})
	}

	return phases

}

func randomlyAssignIngredientsToPhases(t *testing.T, ingredients []Ingredient, phases []Phase) []FormulaIngredient {
	rand.Seed(time.Now().UnixNano())

	var formulaIngredients = make([]FormulaIngredient, 0)
	for _, ingredient := range ingredients {
		randomIndex := rand.Intn(len(phases))
		formulaIngredients = append(formulaIngredients, CreateRandomFormulaIngredient(t, ingredient.ID, phases[randomIndex].ID))
	}
	return formulaIngredients
}

func createPhases(t *testing.T, formulaID int64, count int) []Phase {

	var phases = make([]Phase, 0)
	for i := 0; i < count; i++ {
		phases = append(phases, CreateRandomPhase(t, formulaID))
	}
	return phases
}

// TODO create helper method to generate Ingredients like this for all tests.
func createIngredients(t *testing.T, userID int64, count int) []Ingredient {

	var ingredients = make([]Ingredient, 0)
	for i := 0; i < count; i++ {
		ing := createRandomIngredient(t, userID)
		ingredients = append(ingredients, ing)
	}
	return ingredients

}

func TestUpdateFormulaTx(t *testing.T) {
	userAccount := NewUserAccount(testDB)
	user := CreateRandomUser(t) // TODO add test with incorrect User? - probably should be caught at permissions level and never get here, but why not?

	ingredient := createRandomIngredient(t, user.ID)
	ingredientToMoveOnUpdate := createRandomIngredient(t, user.ID)
	ingredientToRemoveOnUpdate := createRandomIngredient(t, user.ID)
	ingredientToAddOnUpdate := createRandomIngredient(t, user.ID)
	newPhaseAddedIngredient := createRandomIngredient(t, user.ID)
	// TODO check all required parameters and make asserts for them being missing causing a rollback

	formula := CreateRandomFormula(t, user.ID)
	phase := CreateRandomPhase(t, formula.ID)

	formulaIngredient := CreateRandomFormulaIngredient(t, ingredient.ID, phase.ID)
	phaseMoveFormulaIngredient := CreateRandomFormulaIngredient(t, ingredientToMoveOnUpdate.ID, phase.ID)
	CreateRandomFormulaIngredient(t, ingredientToRemoveOnUpdate.ID, phase.ID)

	ingredientParams := ingredientToUpdateParams(formulaIngredient.ID, ingredient)
	addIngredientParams := ingredientToUpdateParams(0, ingredientToAddOnUpdate)
	moveIngredientParams := ingredientToUpdateParams(phaseMoveFormulaIngredient.ID, ingredientToMoveOnUpdate)
	ingredients := []UpdateFullFormulaIngredientParams{
		ingredientParams,
		addIngredientParams,
	}

	addPhaseIngredientParams := ingredientToUpdateParams(0, newPhaseAddedIngredient)
	newPhaseIngredients := []UpdateFullFormulaIngredientParams{
		addPhaseIngredientParams,
		moveIngredientParams,
	}

	originalPhaseParams := UpdateFullFormulaPhaseParams{
		PhaseName:        F.Genre().Name(),
		PhaseDescription: F.Lorem().Sentence(5),
		PhaseId:          phase.ID,
		Ingredients:      ingredients,
	}
	addPhaseParams := UpdateFullFormulaPhaseParams{
		PhaseName:        F.Genre().Name(),
		PhaseDescription: F.Lorem().Sentence(7),
		PhaseId:          0,
		Ingredients:      newPhaseIngredients,
	}
	phases := []UpdateFullFormulaPhaseParams{
		originalPhaseParams,
		addPhaseParams,
	}

	expectedFormulaParams := UpdateFullFormulaParams{
		FormulaId:          formula.ID,
		FormulaName:        F.Music().Author(),
		Weight:             50,
		WeightOz:           50 / 28.34,
		FormulaDescription: F.Lorem().Sentence(5),
		Phases:             phases,
		UserId:             user.ID,
	}

	txResultFormula, err := userAccount.UpdateFullFormulaTx(context.Background(), expectedFormulaParams)
	require.NoError(t, err)
	require.NotEmpty(t, txResultFormula)

	updatedFormula, err := userAccount.GetFullFormula(context.Background(), expectedFormulaParams.FormulaId)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFormula)
	actualFormula, err := userAccount.GetFormula(context.Background(), expectedFormulaParams.FormulaId)
	require.NoError(t, err)
	require.NotEmpty(t, actualFormula)

	updatedPhases := GenerateFormulaViewModelPhases(updatedFormula)

	formulaResponse := FormulaResponse{
		ID:            actualFormula.ID,
		Phases:        updatedPhases,
		Name:          actualFormula.Name, // TODO actually generate a Formula Object instead of getting it from FullFormulaRows.
		TotalWeight:   float64(actualFormula.DefaultAmount),
		TotalWeightOz: float64(actualFormula.DefaultAmountOz),
		Description:   actualFormula.Description,
		CreatedAt:     actualFormula.CreatedAt.Format("06-01-02"),
		UpdatedAt:     actualFormula.UpdatedAt.Format("06-01-02"),
	}

	compareFormulaProperties(t, expectedFormulaParams, actualFormula)
	validatePhasesEquals(t, expectedFormulaParams, formulaResponse)
	// TODO kerok: Add more test scenarios

	for _, phaseToRemove := range txResultFormula.Phases {
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

func validatePhasesEquals(t *testing.T, expectedFormulaProperties UpdateFullFormulaParams, txResultFormulaProperties FormulaResponse) {

	for i, txPhase := range txResultFormulaProperties.Phases {
		formulaPhaseParams := expectedFormulaProperties.Phases[i]
		comparePhaseProperties(t, formulaPhaseParams, txPhase)
		fullFormulaIngredientParams := formulaPhaseParams.Ingredients
		require.Equal(t, len(fullFormulaIngredientParams), len(txPhase.FormulaIngredients))
		for j, phaseIngredient := range txPhase.FormulaIngredients {
			compareIngredientProperties(t, fullFormulaIngredientParams[j], phaseIngredient, txPhase)
		}
	}
}

func compareIngredientProperties(t *testing.T, fullFormulaIngredientParams UpdateFullFormulaIngredientParams, phaseIngredient FullFormulaIngredient, phase FullFormulaPhase) {
	require.Equal(t, fullFormulaIngredientParams.IngredientId, phaseIngredient.IngredientId)
	require.Equal(t, fullFormulaIngredientParams.FormulaIngredientPercentage, phaseIngredient.Percentage)
	require.Equal(t, fullFormulaIngredientParams.FormulaIngredientCost, phaseIngredient.Cost)
}

func comparePhaseProperties(t *testing.T, formulaPhaseParams UpdateFullFormulaPhaseParams, phase FullFormulaPhase) {
	require.Equal(t, formulaPhaseParams.PhaseName, phase.Name)
	require.Equal(t, len(formulaPhaseParams.Ingredients), len(phase.FormulaIngredients))
	if formulaPhaseParams.PhaseId == 0 {
		require.NotEqual(t, formulaPhaseParams.PhaseId, phase.ID)
	} else {
		require.Equal(t, formulaPhaseParams.PhaseId, phase.ID)
	}
	require.Equal(t, formulaPhaseParams.PhaseDescription, phase.Description)
}

func compareFormulaProperties(t *testing.T, expectedFormula UpdateFullFormulaParams, actualFormulaResponse Formula) {
	require.Equal(t, expectedFormula.FormulaId, actualFormulaResponse.ID)
	require.Equal(t, expectedFormula.FormulaName, actualFormulaResponse.Name)
	//require.Equal(t, expectedFormula.UserId, actualFormulaResponse.UserID)
	require.Equal(t, expectedFormula.FormulaDescription, actualFormulaResponse.Description)
	require.Equal(t, expectedFormula.Weight, actualFormulaResponse.DefaultAmount)
	require.Equal(t, expectedFormula.WeightOz, actualFormulaResponse.DefaultAmountOz)
}

func ingredientToUpdateParams(formulaIngredientId int64, ingredient Ingredient) UpdateFullFormulaIngredientParams {
	ingredientParams := UpdateFullFormulaIngredientParams{
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
		Name:        F.Food().Fruit(),
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
