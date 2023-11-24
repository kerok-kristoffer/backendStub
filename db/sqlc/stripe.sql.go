// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: stripe.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createStripeEntry = `-- name: CreateStripeEntry :one
INSERT INTO stripe (
    id,
    user_id,
    stripe_customer_id,
    stripe_plan_id
) VALUES (
          $1, $2, $3, $4
         ) RETURNING id, user_id, stripe_customer_id, stripe_plan_id, created_at, updated_at
`

type CreateStripeEntryParams struct {
	ID               uuid.UUID      `json:"id"`
	UserID           int64          `json:"userID"`
	StripeCustomerID sql.NullString `json:"stripeCustomerID"`
	StripePlanID     uuid.UUID      `json:"stripePlanID"`
}

func (q *Queries) CreateStripeEntry(ctx context.Context, arg CreateStripeEntryParams) (Stripe, error) {
	row := q.db.QueryRowContext(ctx, createStripeEntry,
		arg.ID,
		arg.UserID,
		arg.StripeCustomerID,
		arg.StripePlanID,
	)
	var i Stripe
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StripeCustomerID,
		&i.StripePlanID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getStripeByUserId = `-- name: GetStripeByUserId :one
SELECT id, user_id, stripe_customer_id, stripe_plan_id, created_at, updated_at FROM stripe
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetStripeByUserId(ctx context.Context, userID int64) (Stripe, error) {
	row := q.db.QueryRowContext(ctx, getStripeByUserId, userID)
	var i Stripe
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StripeCustomerID,
		&i.StripePlanID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateStripeByUserId = `-- name: UpdateStripeByUserId :one
UPDATE stripe
SET (stripe_customer_id,
     stripe_plan_id,
     updated_at) =
    ($2, $3, CURRENT_TIMESTAMP)
WHERE user_id = $1
RETURNING id, user_id, stripe_customer_id, stripe_plan_id, created_at, updated_at
`

type UpdateStripeByUserIdParams struct {
	UserID           int64          `json:"userID"`
	StripeCustomerID sql.NullString `json:"stripeCustomerID"`
	StripePlanID     uuid.UUID      `json:"stripePlanID"`
}

func (q *Queries) UpdateStripeByUserId(ctx context.Context, arg UpdateStripeByUserIdParams) (Stripe, error) {
	row := q.db.QueryRowContext(ctx, updateStripeByUserId, arg.UserID, arg.StripeCustomerID, arg.StripePlanID)
	var i Stripe
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StripeCustomerID,
		&i.StripePlanID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
