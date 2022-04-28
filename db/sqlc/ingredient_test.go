package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateIngredient(t *testing.T) {
	user := CreateRandomUser(t)

	arg := CreateIngredientParams{
		Name:   F.Food().Vegetable(),
		UserID: sql.NullInt64{Int64: user.ID, Valid: true},
	}

	ing, err := testQueries.CreateIngredient(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, ing)
	require.Equal(t, arg.Name, ing.Name)
	require.Equal(t, user.ID, ing.UserID.Int64)

	require.NotZero(t, ing.UserID)
	require.NotZero(t, ing.CreatedAt)

	err = testQueries.DeleteIngredient(context.Background(), ing.ID)
	require.NoError(t, err)
	err = testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func TestDeleteIngredientWithConstraint(t *testing.T) {
	user := CreateRandomUser(t)
	ingredient := createRandomIngredient(t, user.ID)

	err := testQueries.DeleteUser(context.Background(), user.ID) // todo kerok - Update below to more reliable catch of error if message changes?
	require.Error(t, err, "pq: update or delete on table \"users\" violates foreign key constraint \"ingredients_user_id_fkey\" on table \"ingredients\"")

	err = testQueries.DeleteIngredient(context.Background(), ingredient.ID)
	require.NoError(t, err)
	ingredient, err = testQueries.GetIngredient(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, ingredient)

	err = testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func TestGetIngredient(t *testing.T) {
	user := CreateRandomUser(t)
	ingredient1 := createRandomIngredient(t, user.ID)

	ingredient2, err := testQueries.GetIngredient(context.Background(), ingredient1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, ingredient2)

	require.Equal(t, ingredient1.ID, ingredient2.ID)
	require.Equal(t, ingredient1.Name, ingredient2.Name)

	require.WithinDurationf(t, ingredient1.CreatedAt, ingredient2.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

	err = testQueries.DeleteIngredient(context.Background(), ingredient2.ID)
	require.NoError(t, err)
	err = testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}

func TestUpdateIngredientName(t *testing.T) {

	user := CreateRandomUser(t)
	ingredient := createRandomIngredient(t, user.ID)

	params := UpdateIngredientNameParams{
		ID:   ingredient.ID,
		Name: F.Food().Vegetable(),
	}

	updatedIngredient, err := testQueries.UpdateIngredientName(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedIngredient)

	require.Equal(t, updatedIngredient.ID, ingredient.ID)
	require.Equal(t, updatedIngredient.Name, params.Name)
	require.Equal(t, updatedIngredient.UserID, ingredient.UserID)
	require.WithinDuration(t, updatedIngredient.CreatedAt, ingredient.CreatedAt, time.Second, "Error, created_at timestamps not within 1sec")

}

func TestListIngredients(t *testing.T) {
	user := CreateRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomIngredient(t, user.ID)
	}

	arg := ListIngredientsParams{
		Limit:  5,
		Offset: 5,
	}

	ingredients, err := testQueries.ListIngredients(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, ingredients, 5)
	for _, ingredient := range ingredients {
		require.NotEmpty(t, ingredient)
	}

	params := ListIngredientsByUserIdParams{
		UserID: sql.NullInt64{Int64: user.ID, Valid: true},
		Limit:  5,
		Offset: 5,
	}

	ingredientsByUserId, err := testQueries.ListIngredientsByUserId(context.Background(), params)
	require.Len(t, ingredientsByUserId, 5)
	for _, ingredientByUserId := range ingredientsByUserId {
		require.NotEmpty(t, ingredientByUserId)
	}

	err = testQueries.DeleteIngredientsByUserId(context.Background(), sql.NullInt64{
		Int64: user.ID, Valid: true,
	})
	require.NoError(t, err)

	shouldBeEmptyIngredientsByUserId, err := testQueries.ListIngredientsByUserId(context.Background(), params)
	require.NoError(t, err)
	require.Empty(t, shouldBeEmptyIngredientsByUserId)
}

func createRandomIngredient(t *testing.T, userId int64) Ingredient {

	arg := CreateIngredientParams{
		Name:   F.Food().Vegetable(),
		UserID: sql.NullInt64{Int64: userId, Valid: true},
	}

	ing, err := testQueries.CreateIngredient(context.Background(), arg)
	require.NoError(t, err)
	return ing
}
