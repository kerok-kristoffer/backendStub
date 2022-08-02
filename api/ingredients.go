package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/token"
)

type listIngredientsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listIngredients(ctx *gin.Context) {
	var req listIngredientsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListIngredientsByUserIdParams{
		UserID: sql.NullInt64{Int64: int64(authPayLoad.ID.ID()), Valid: true},
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	ingredients, err := server.userAccount.ListIngredientsByUserId(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ingredients)
}

type addIngredientRequest struct {
	Name string `json:"Name" binding:"required,alphanum"`
}

type ingredientResponse struct {
	Name string `json:"Name" binding:"required,alphanum"`
}

func newIngredientResponse(ingredient db.Ingredient) ingredientResponse {
	return ingredientResponse{
		Name: ingredient.Name,
	}
}

func (server Server) addIngredient(ctx *gin.Context) {
	var req addIngredientRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateIngredientParams{
		Name:   req.Name,
		Hash:   "", // todo kerok - remove parameter from Sqlc generator
		UserID: sql.NullInt64{Int64: int64(authPayLoad.ID.ID()), Valid: true},
	}

	ingredient, err := server.userAccount.CreateIngredient(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newIngredientResponse(ingredient))

}
