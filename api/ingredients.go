package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/token"
)

type addFunctionRequest struct {
	Name string
}

func (server *Server) addIngredientFunction(ctx *gin.Context) {
	var req addFunctionRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.userAccount.GetUserByUserName(ctx, authPayLoad.Username)

	arg := db.CreateIngredientFunctionParams{
		Name:   req.Name,
		UserID: user.ID,
	}

	function, err := server.userAccount.CreateIngredientFunction(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, function)
}

func (server *Server) listIngredientFunctions(ctx *gin.Context) {

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.userAccount.GetUserByUserName(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	params := db.ListIngredientFunctionsByUserIdParams{
		UserID: user.ID,
		Limit:  50, // TODO add actual pagination if needed
		Offset: 0,
	}

	functions, err := server.userAccount.ListIngredientFunctionsByUserId(ctx, params)
	makeFunctionsViewModel(functions)
}

type functionsResponse struct {
	Id   int64  `json:"Id" binding:"required"`
	Name string `json:"Name" binding:"required"`
}

func makeFunctionsViewModel(functions []db.IngredientFunction) {
	var viewModels []functionsResponse
	for _, function := range functions {
		viewModels = append(viewModels, functionsResponse{
			Id:   function.ID,
			Name: function.Name,
		})
	}
}

type listIngredientsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listIngredients(ctx *gin.Context) {
	var req listIngredientsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.userAccount.GetUserByUserName(ctx, authPayLoad.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := ingredientByUserParams(user, req)

	ingredients, err := server.userAccount.ListIngredientsByUserId(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, makeViewModel(ingredients))
}

func ingredientByUserParams(user db.User, req listIngredientsRequest) db.ListIngredientsByUserIdParams {
	arg := db.ListIngredientsByUserIdParams{
		UserID: user.ID,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	return arg
}

func makeViewModel(ingredients []db.Ingredient) []ingredientResponse {
	var viewModelIngredients []ingredientResponse
	for i := range ingredients {
		viewModelIngredients = append(viewModelIngredients, ingredientResponse{
			Id:   ingredients[i].ID,
			Name: ingredients[i].Name,
			Inci: ingredients[i].Inci,
		})
	}
	return viewModelIngredients
}

type addIngredientRequest struct {
	Name string `json:"name" binding:"required"`
	Inci string `json:"inci" binding:"required"`
}

type ingredientResponse struct {
	Id   int64  `json:"Id" binding:"required"`
	Name string `json:"Name" binding:"required"`
	Inci string `json:"Inci" binding:"required"`
}

func newIngredientResponse(ingredient db.Ingredient) ingredientResponse {
	return ingredientResponse{
		Id:   ingredient.ID,
		Name: ingredient.Name,
		Inci: ingredient.Inci,
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
	user, err := server.userAccount.GetUserByUserName(ctx, authPayLoad.Username)

	arg := db.CreateIngredientParams{
		Name:   req.Name,
		Inci:   req.Inci,
		Hash:   "",
		UserID: user.ID,
	}

	ingredient, err := server.userAccount.CreateIngredient(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newIngredientResponse(ingredient))
}
