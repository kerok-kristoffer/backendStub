package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/kerok-kristoffer/formulating/db/models"
)

// UserAccount Interface representing SQL and Mock version
// Mock is generated as per below automatically in Makefile
//go:generate mockgen -package mockdb -destination ../mock/user_account.go github.com/kerok-kristoffer/formulating/db/sqlc UserAccount
type UserAccount interface { // todo kerok - rename interface at some point?
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	UpdateFullFormulaTx(ctx context.Context, arg models.UpdateFullFormulaParams) (UpdateFormulaTxResult, error)
}

type PhaseModel struct {
}

type FormulaModel struct {
	Name   string `json:"name" binding:"required"`
	Phases []*PhaseModel
}

type SQLUserAccount struct {
	// todo corresponds to store in tut
	*Queries
	db       *sql.DB
	Formulas []*FormulaModel
}

func NewUserAccount(db *sql.DB) UserAccount {
	return &SQLUserAccount{
		db:      db,
		Queries: New(db),
	}
}

// DB transaction execution example, not in use currently
func (userAccount *SQLUserAccount) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := userAccount.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromUserID int64 `json:"fromUserID"`
	ToUserID   int64 `json:"toUserID"`
	Amount     int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer   Transfer `json:"transfer"`
	FromUserID User     `json:"fromUserID"`
	ToUserID   User     `json:"toUserID"`
	FromEntry  Entry    `json:"fromEntry"`
	ToEntry    Entry    `json:"toEntry"`
}

type UpdatePhaseTxResult struct {
	Phase       Phase               `json:"phase"`
	Ingredients []FormulaIngredient `json:"ingredients"`
}

type UpdateFormulaTxResult struct {
	Formula Formula               `json:"formula"`
	Phases  []UpdatePhaseTxResult `json:"phases"`
}

func (userAccount *SQLUserAccount) UpdateFullFormulaTx(ctx context.Context, arg models.UpdateFullFormulaParams) (UpdateFormulaTxResult, error) {
	var result UpdateFormulaTxResult

	updateId := uuid.New()
	err := userAccount.execTx(ctx, func(q *Queries) error {
		var formulaPhases = new([]UpdatePhaseTxResult)
		for _, phase := range arg.Phases {

			phaseTxResult, err := addOrUpdatePhase(q, phase, ctx, arg, updateId)
			if err != nil {
				return err
			}
			phase.PhaseId = phaseTxResult.ID

			var phaseIngredients = new([]FormulaIngredient)
			for _, ingredient := range phase.Ingredients {
				ingredientTxResult, err := addOrUpdateFormulaIngredient(q, ingredient, phase, ctx, updateId)
				if err != nil {
					return err
				}
				*phaseIngredients = append(*phaseIngredients, ingredientTxResult)
			}

			updatePhaseTxResult := UpdatePhaseTxResult{
				Phase:       phaseTxResult,
				Ingredients: *phaseIngredients,
			}
			*formulaPhases = append(*formulaPhases, updatePhaseTxResult)
		}

		phases, err := q.ListPhasesByFormulaId(ctx, arg.FormulaId)
		if err != nil {
			return err
		}
		err = DeleteDiscardedIngredientsAndPhases(phases, q, ctx, updateId)
		if err != nil {
			return err
		}

		formulaTxResult, err := q.UpdateFormula(ctx, UpdateFormulaParams{
			ID:            arg.FormulaId,
			Name:          arg.FormulaName,
			DefaultAmount: arg.Weight,
			Description:   arg.FormulaDescription,
			UserID:        arg.UserId,
		})
		if err != nil {
			return err
		}

		result = UpdateFormulaTxResult{
			Formula: formulaTxResult,
			Phases:  *formulaPhases,
		}
		return nil
	})
	return result, err
}

func DeleteDiscardedIngredientsAndPhases(phases []Phase, q *Queries, ctx context.Context, updateId uuid.UUID) error {
	for _, phase := range phases { // TODO add consolidating SQL with join on formulaID -> phase.formulaId instead of
		err := q.DeleteFormulaIngredientsNotInUpdate(ctx, DeleteFormulaIngredientsNotInUpdateParams{
			PhaseID:  phase.ID,
			UpdateID: updateId,
		})
		if err != nil {
			return err
		}
		if phase.UpdateID != updateId {
			err := q.DeletePhase(ctx, phase.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addOrUpdatePhase(q *Queries, phase models.UpdateFullFormulaPhaseParams, ctx context.Context, arg models.UpdateFullFormulaParams, updateId uuid.UUID) (Phase, error) {
	var err error
	var phaseTxResult Phase
	isAddedPhase := phase.PhaseId < 1
	if isAddedPhase {
		phaseTxResult, err = q.CreatePhase(ctx, CreatePhaseParams{
			Name:        phase.PhaseName,
			Description: phase.PhaseDescription,
			FormulaID:   arg.FormulaId,
			UpdateID:    updateId,
		})

		if err != nil {
			return Phase{}, err
		}
	} else {
		phaseTxResult, err = q.UpdatePhase(ctx, UpdatePhaseParams{
			ID:          phase.PhaseId,
			Name:        phase.PhaseName,
			Description: phase.PhaseDescription,
			FormulaID:   arg.FormulaId,
			UpdateID:    updateId,
		})
		if err != nil {
			return Phase{}, err
		}
	}
	return phaseTxResult, nil
}

func addOrUpdateFormulaIngredient(q *Queries, ingredient models.UpdateFullFormulaIngredientParams, phase models.UpdateFullFormulaPhaseParams, ctx context.Context, updateId uuid.UUID) (FormulaIngredient, error) {
	var ingredientTxResult FormulaIngredient
	var err error
	isNewIngredient := ingredient.FormulaIngredientId == 0
	if isNewIngredient {
		formulaIngredientParams := CreateFormulaIngredientParams{
			IngredientID: ingredient.IngredientId,
			Percentage:   ingredient.FormulaIngredientPercentage,
			PhaseID:      phase.PhaseId,
			Description:  sql.NullString{},
			UpdateID:     updateId, // TODO potentially keep dangling ingredients in a update post in sql as backup?
		}
		ingredientTxResult, err = q.CreateFormulaIngredient(ctx, formulaIngredientParams)
	} else {
		params := UpdateFormulaIngredientParams{
			ID:           ingredient.FormulaIngredientId,
			IngredientID: ingredient.IngredientId,
			Percentage:   ingredient.FormulaIngredientPercentage,
			PhaseID:      phase.PhaseId,
			Cost:         sql.NullFloat64{Float64: float64(ingredient.FormulaIngredientCost), Valid: true},
			Description:  sql.NullString{},
			UpdateID:     updateId,
		}
		ingredientTxResult, err = q.UpdateFormulaIngredient(ctx, params)
	}
	if err != nil {
		return FormulaIngredient{}, err
	}
	return ingredientTxResult, nil
}

func (userAccount *SQLUserAccount) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// TODO kerok - trying out transactions, no current need in project, keep for future reference example
	err := userAccount.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromUserID: sql.NullInt64{Int64: arg.FromUserID, Valid: true},
			ToUserID:   sql.NullInt64{Int64: arg.ToUserID, Valid: true},
			Amount:     sql.NullInt64{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			UserID: sql.NullInt64{Int64: arg.FromUserID, Valid: true},
			Amount: sql.NullInt64{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			UserID: sql.NullInt64{Int64: arg.ToUserID, Valid: true},
			Amount: sql.NullInt64{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}
		// TODO update accounts - probably will skip this since I'm not really implementing transfers in this way.

		return nil
	})

	return result, err
}
