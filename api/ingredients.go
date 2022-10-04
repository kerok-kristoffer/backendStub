package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
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

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

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

	user, err := server.getAuthenticatedUser(ctx)
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

func (server *Server) getIngredientCount(ctx *gin.Context) {
	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	count, err := server.userAccount.GetIngredientCount(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, count)
}

type listIngredientsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listIngredients(ctx *gin.Context) {
	var req listIngredientsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getAuthenticatedUser(ctx)
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

type updateIngredientRequest struct {
	Id   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Inci string `json:"inci"`
}

func (server *Server) updateIngredient(ctx *gin.Context) {
	var req updateIngredientRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	validateIngredient, err := server.userAccount.GetIngredient(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if validateIngredient.UserID != user.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	params := db.UpdateIngredientParams{
		ID:         req.Id,
		Name:       req.Name,
		Inci:       req.Inci,
		Hash:       "",
		UserID:     user.ID,
		FunctionID: sql.NullInt64{},
	}

	ingredient, err := server.userAccount.UpdateIngredient(ctx, params)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newIngredientResponse(ingredient))
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
	Inci string `json:"inci"`
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

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

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

func (server *Server) deleteIngredient(ctx *gin.Context) {

	ingredientId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	validateIngredient, err := server.userAccount.GetIngredient(ctx, ingredientId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if validateIngredient.UserID != user.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.userAccount.DeleteIngredient(ctx, ingredientId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
