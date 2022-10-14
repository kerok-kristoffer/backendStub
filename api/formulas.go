package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kerok-kristoffer/formulating/db/models"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/token"
	"log"
	"net/http"
)

type FormulaIngredient struct {
	Id             int64   `json:"id"`
	IngredientId   int64   `json:"ingredient_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Inci           string  `json:"inci" binding:"required"`
	Percentage     int64   `json:"percentage" binding:"required"`
	WeightInGrams  float32 `json:"weightInGrams" binding:"required"`
	WeightInOunces float32 `json:"weightInOunces" binding:"required"`
	Description    string  `json:"description" binding:"required"`
}

type Phase struct {
	ID                 int64               `json:"id" binding:"required"`
	Name               string              `json:"name" binding:"required"`
	FormulaIngredients []FormulaIngredient `json:"ingredients" binding:"required"`
	Description        string              `json:"description" binding:"required"`
}

type updateFormulaRequest struct {
	Id            int64                                 `json:"id" binding:"required"`
	Phases        []models.UpdateFullFormulaPhaseParams `json:"phases" binding:"required"`
	Name          string                                `json:"name" binding:"required"`
	TotalWeight   float32                               `json:"totalWeight" binding:"required"`
	TotalWeightOz float32                               `json:"totalWeightInOunces" binding:"required"`
	Description   string                                `json:"description" binding:"required"`
}

func (server Server) updateFormula(ctx *gin.Context) {
	var req updateFormulaRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	formula, err := server.userAccount.GetFormula(ctx, req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if formula.UserID != user.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, err = server.userAccount.UpdateFullFormulaTx(ctx, models.UpdateFullFormulaParams{
		FormulaId:          req.Id,
		FormulaName:        req.Name,
		Weight:             req.TotalWeight,
		WeightOz:           req.TotalWeightOz,
		FormulaDescription: req.Description,
		Phases:             req.Phases,
		UserId:             user.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fullFormula, err := server.userAccount.GetFullFormula(ctx, formula.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	response := generateFormulaViewModel(fullFormula)

	ctx.JSON(http.StatusOK, response)

}

type addFormulaRequest struct {
	Phases        []Phase `json:"phases" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	TotalWeight   float32 `json:"totalWeight" binding:"required"`
	TotalWeightOz float32 `json:"totalWeightInOunces" binding:"required"`
	Description   string  `json:"description" binding:"required"`
}

func (server Server) addFormula(ctx *gin.Context) {

	var req addFormulaRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	formulaParams := createFormulaParamsFromRequest(req, user)
	formula, err := server.userAccount.CreateFormula(ctx, formulaParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	var response formulaResponse
	phases := req.Phases
	for i := range req.Phases {
		phaseRequest := req.Phases[i]
		phaseParams := createPhaseParamsFromRequest(phaseRequest, formula.ID)
		savedPhase, err := server.userAccount.CreatePhase(ctx, phaseParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		var ingredients []FormulaIngredient

		for _, ingredient := range phaseRequest.FormulaIngredients {
			formulaIngredientParams := CreateIngredientParamsFromRequest(ingredient, savedPhase)
			_, err := server.userAccount.CreateFormulaIngredient(ctx, formulaIngredientParams)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return // todo: rollback?
			}
			ingredients = append(ingredients, ingredient)
		}
		response = newFormulaResponse(formula, phases)
	}
	ctx.JSON(http.StatusOK, response)
}

func (server Server) getAuthenticatedUser(ctx *gin.Context) (db.User, error) {
	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.userAccount.GetUserByUserName(ctx, authPayLoad.Username)
	return user, err
}

type formulaResponse struct {
	ID            int64   `json:"id" binding:"required"`
	Phases        []Phase `json:"phases"`
	Name          string  `json:"name"`
	TotalWeight   float64 `json:"totalWeight"`
	TotalWeightOz float64 `json:"totalWeightInOunces"`
	Description   string  `json:"description"`
}

type listFormulasRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (server Server) listFormulas(ctx *gin.Context) {
	var req listFormulasRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	listFormulaParams := db.ListFormulasByUserIdParams{
		UserID: user.ID,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	formulas, err := server.userAccount.ListFormulasByUserId(ctx, listFormulaParams)
	if err != nil {
		return
	}

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	formulaViewModels, err := makeFormulaViewModels(formulas, server, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, formulaViewModels)
}

func makeFormulaViewModels(formulas []db.Formula, server Server, ctx *gin.Context) ([]formulaResponse, error) {
	var formulaViewModels []formulaResponse
	for _, formula := range formulas {
		fullFormulaIngredients, err := server.userAccount.GetFullFormula(ctx, formula.ID)
		if err != nil {
			return nil, err
		}

		if len(fullFormulaIngredients) > 0 {
			formulaResponse := generateFormulaViewModel(fullFormulaIngredients)
			formulaViewModels = append(formulaViewModels, formulaResponse)
		}
	}
	return formulaViewModels, nil
}

func generateFormulaViewModel(fullFormulaIngredients []db.GetFullFormulaRow) formulaResponse {
	var formulaPhases = make(map[int64]Phase)
	var phase Phase

	for _, ingredient := range fullFormulaIngredients {
		phase, formulaPhases = getOrCreatePhase(formulaPhases, ingredient)
		formulaIngredientModel := FormulaIngredient{
			Id:           ingredient.FormulaIngredientID,
			IngredientId: ingredient.IngredientID,
			Name:         ingredient.IngredientName,
			Inci:         ingredient.Inci,
			Percentage:   int64(ingredient.Percentage),
		}
		phase.FormulaIngredients = append(phase.FormulaIngredients, formulaIngredientModel)
		formulaPhases[ingredient.PhaseID] = phase
	}

	formulaPhaseModels := new([]Phase)
	for _, phase := range formulaPhases {
		*formulaPhaseModels = append(*formulaPhaseModels, phase)
	}

	fullFormulaInfo := fullFormulaIngredients[0]
	formulaResponse := formulaResponse{
		ID:            fullFormulaInfo.FormulaID,
		Phases:        *formulaPhaseModels,
		Name:          fullFormulaInfo.FormulaName,
		TotalWeight:   float64(fullFormulaInfo.DefaultAmount),
		TotalWeightOz: float64(fullFormulaInfo.DefaultAmountOz),
		Description:   fullFormulaInfo.Description,
	}
	return formulaResponse
}

func getOrCreatePhase(phases map[int64]Phase, ingredient db.GetFullFormulaRow) (Phase, map[int64]Phase) {
	phase, exists := phases[ingredient.PhaseID]
	if exists {
		return phase, phases
	}

	formulaIngredients := new([]FormulaIngredient)
	phase = Phase{
		ID:                 ingredient.PhaseID,
		Name:               ingredient.PhaseName,
		FormulaIngredients: *formulaIngredients,
		Description:        ingredient.PhaseDescription,
	}

	phases[ingredient.PhaseID] = phase
	return phase, phases

}

func newFormulaResponse(formula db.Formula, phases []Phase) formulaResponse {

	return formulaResponse{
		ID:            formula.ID,
		Phases:        phases,
		Name:          formula.Name,
		TotalWeight:   float64(formula.DefaultAmount),
		TotalWeightOz: float64(formula.DefaultAmountOz),
		Description:   formula.Description,
	}
}

func CreateIngredientParamsFromRequest(ingredient FormulaIngredient, savedPhase db.Phase) db.CreateFormulaIngredientParams {
	formulaIngredientParams := db.CreateFormulaIngredientParams{
		IngredientID: ingredient.IngredientId,
		Percentage:   int32(ingredient.Percentage),
		PhaseID:      savedPhase.ID,
		Description:  sql.NullString{String: ingredient.Description, Valid: true},
	}
	return formulaIngredientParams
}

func createPhaseParamsFromRequest(phase Phase, id int64) db.CreatePhaseParams {
	phaseParams := db.CreatePhaseParams{
		Name:        phase.Name,
		Description: phase.Description,
		FormulaID:   id,
	}
	return phaseParams
}

func createFormulaParamsFromRequest(req addFormulaRequest, user db.User) db.CreateFormulaParams {
	formulaParams := db.CreateFormulaParams{
		Name:            req.Name,
		DefaultAmount:   req.TotalWeight,
		DefaultAmountOz: req.TotalWeightOz,
		Description:     req.Description,
		UserID:          user.ID,
	}
	return formulaParams
}
