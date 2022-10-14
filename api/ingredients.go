package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

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

	var responseIngredients []ingredientResponse
	for _, ingredient := range ingredients {
		tagMaps, err := server.userAccount.ListIngredientTagsByIngredientId(ctx, ingredient.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		var tags []ingredientTagRequest
		for _, tagMap := range tagMaps {
			tags = append(tags, ingredientTagRequest{
				Id:   tagMap.IngredientTagsID,
				Name: tagMap.IngredientTag,
			})
		}
		responseIngredients = append(responseIngredients, newIngredientResponse(ingredient, tags))
	}

	ctx.JSON(http.StatusOK, responseIngredients)
}

type ingredientTagRequest struct {
	Id   int64  `json:"id"`
	Name string `json:"name" binding:"required"`
}

type updateIngredientRequest struct {
	Id   int64                  `json:"id" binding:"required"`
	Name string                 `json:"name" binding:"required"`
	Inci string                 `json:"inci"`
	Cost int                    `json:"cost"`
	Tags []ingredientTagRequest `json:"tags" binding:"required"`
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
		Cost:       sql.NullInt32{Int32: int32(req.Cost), Valid: true},
		UserID:     user.ID,
		FunctionID: sql.NullInt64{},
	}

	ingredient, err := server.userAccount.UpdateIngredient(ctx, params)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// update tags
	var tags []ingredientTagRequest
	for _, tag := range req.Tags {

		ingredientTag, err := server.userAccount.GetIngredientTagByName(ctx, db.GetIngredientTagByNameParams{
			Name:   tag.Name,
			UserID: user.ID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				ingredientTag, err = server.userAccount.CreateIngredientTag(ctx, db.CreateIngredientTagParams{
					Name:   strings.ToLower(tag.Name),
					UserID: user.ID,
				})
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
			} else {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}

		_, err = server.userAccount.GetTagMap(ctx, db.GetTagMapParams{
			IngredientID:    ingredient.ID,
			IngredientTagID: ingredientTag.ID,
		})

		if err == sql.ErrNoRows {
			_, err := server.userAccount.CreateIngredientTagMap(ctx, db.CreateIngredientTagMapParams{
				IngredientTagID: ingredientTag.ID,
				IngredientID:    req.Id,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

	}

	tagMaps, err := server.userAccount.ListIngredientTagsByIngredientId(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	// TODO: Remove existing ingredient_tag_map with current ingredient_id if not in tags table.
	for _, tagMap := range tagMaps {
		if !isInProvidedTags(req.Tags, tagMap.IngredientTag) {
			err := server.userAccount.DeleteIngredientTagMap(ctx, tagMap.IngredientTagMapsID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		tags = append(tags, ingredientTagRequest{
			Id:   tagMap.IngredientTagsID,
			Name: tagMap.IngredientTag,
		})
	}
	ctx.JSON(http.StatusOK, newIngredientResponse(ingredient, tags))
}

func isInProvidedTags(tags []ingredientTagRequest, str string) bool {
	for _, tag := range tags {
		if strings.ToLower(tag.Name) == str {
			return true
		}
	}
	return false
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
	Id   int64                  `json:"Id" binding:"required"`
	Name string                 `json:"Name" binding:"required"`
	Inci string                 `json:"Inci" binding:"required"`
	Cost int32                  `json:"cost"`
	Tags []ingredientTagRequest `json:"tags" binding:"required"`
}

func newIngredientResponse(ingredient db.Ingredient, tags []ingredientTagRequest) ingredientResponse {
	return ingredientResponse{
		Id:   ingredient.ID,
		Name: ingredient.Name,
		Inci: ingredient.Inci,
		Cost: ingredient.Cost.Int32,
		Tags: tags,
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

	var noTags []ingredientTagRequest
	ctx.JSON(http.StatusOK, newIngredientResponse(ingredient, noTags))
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
