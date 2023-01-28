package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/kerok-kristoffer/formulating/db/sqlc"
	"github.com/kerok-kristoffer/formulating/token"
	"github.com/kerok-kristoffer/formulating/util"
	"github.com/kerok-kristoffer/formulating/util/access"
	"github.com/lib/pq"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"net/http"
	"time"
)

type getUserRequest struct {
	UserID int64 `uri:"id" binding:"required,min=1"`
}

type applySubscriptionRequest struct {
	PriceID    string `uri:"price_id" json:"price_id" binding:"required,min=1"`
	SuccessUrl string `json:"success_url" uri:"success_url" binding:"required,min=1"`
	CancelUrl  string `json:"cancel_url" uri:"cancel_url" binding:"required,min=1"`
}

func (server *Server) applySubscription(ctx *gin.Context) {

	authenticatedUser, err := server.getAuthenticatedUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	var req applySubscriptionRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: before pushing backend and frontend
	// TODO test new request parameters, email, succ and cancel,
	// TODO run migrations on live db

	// TODO verify that redirects work on live page
	// TODO add restrictions on access depending on sub level

	stripe.Key = server.config.StripeKey

	priceId := req.PriceID // todo: add to front-end, comes from product page on Stripe dashboard
	stripeTestCustomerEmail := authenticatedUser.Email
	stripeTestCustomerName := authenticatedUser.UserName

	stripePlanInternal, err := server.userAccount.GetStripePlanByUserAccess(ctx, access.NONE)
	stripeUserData, err := server.userAccount.GetStripeByUserId(ctx, authenticatedUser.ID)
	if err == sql.ErrNoRows {
		stripeUserData, err = server.userAccount.CreateStripeEntry(ctx, db.CreateStripeEntryParams{
			ID:           uuid.New(),
			UserID:       authenticatedUser.ID,
			StripePlanID: stripePlanInternal.ID, // todo kerok fix StripePlanId implementation!
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	stripeCustomerParams := &stripe.CustomerParams{
		Email: stripe.String(stripeTestCustomerEmail),
		Name:  stripe.String(stripeTestCustomerName),
	}

	var stripeCustomer *stripe.Customer
	if stripeUserData.StripeCustomerID.Valid == false {
		stripeCustomer, err = customer.New(stripeCustomerParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		stripeCustomer, err = customer.Get(stripeUserData.StripeCustomerID.String, nil)
		if err != nil {
			stripeCustomer, err = customer.New(stripeCustomerParams)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	stripeUserData, err = server.userAccount.UpdateStripeByUserId(ctx, db.UpdateStripeByUserIdParams{
		UserID:           authenticatedUser.ID,
		StripeCustomerID: sql.NullString{String: stripeCustomer.ID, Valid: true},
		StripePlanID:     stripePlanInternal.ID,
	}) // Todo kerok : temporary using random generated stripePlanId since we only have one
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // Todo kerok : this really shouldn't be allowed to happen, do we need a rollback if it does?!
		return
	}
	// todo kerok : stripe plan update needs to be done after a success is retrieved from Stripe
	// todo kerok : there should be support for checking if a Customer's purchase went through in the API
	// todo kerok : We should probably check this and update status every time we authenticate the user

	stripeCheckoutSessionParams := &stripe.CheckoutSessionParams{ // todo kerok : set up the sub portal link on front-end
		Customer:   stripe.String(stripeCustomer.ID),
		SuccessURL: stripe.String(req.SuccessUrl), // todo kerok : set up success and cancel pages on front-end
		CancelURL:  stripe.String(req.CancelUrl),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
	}

	checkoutSession, err := session.New(stripeCheckoutSessionParams)

	// todo kerok : add database table for subscriptions, tied to user.ID *DONE*
	// todo kerok: set up subLvl table to keep track of the different products
	// todo: perhaps keep the list of available subs to fetch from the front-end here.
	// todo kerok : needed for storing things like stripe.CustomerID, subscriptionLevel, etc...

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, stripeCheckoutRedirectResponse{
		ID:  checkoutSession.ID,
		URL: checkoutSession.URL,
	})
	//ctx.Redirect(http.StatusSeeOther, checkoutSession.URL)
}

type stripeCheckoutRedirectResponse struct {
	ID  string `form:"id"`
	URL string `form:"url"`
}

func (server *Server) getSubscriptions(ctx *gin.Context) {

}

func (server *Server) getUserAccount(ctx *gin.Context) {
	// todo kerok - refactor this and server.go - concept of UserAccount from tut might not fit my purposes?
	// in tut#22, at around 18m, implement corresponding middleware authentication for routes listing ingredients, recipes, etc. (instead of accounts like in tut)
	// listing users might use a different middleware for admin, etc.
	// could add a listFollowers or listFriends when adding social media functionality to site.
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.userAccount.GetUser(ctx, req.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("not authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUsersRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	users, err := server.userAccount.ListUsers(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type createUserRequest struct {
	UserName string `json:"userName" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	UserName  string    `json:"userName" binding:"required,alphanum"`
	Email     string    `json:"email" binding:"required,email"`
	FullName  string    `json:"fullName" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		UserName:  user.UserName,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
	}
}

func (server Server) createUser(ctx *gin.Context) {

	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tester, err := server.userAccount.GetTesterByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("email not registered as tester")))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		UserName: req.UserName,
		Email:    req.Email,
		FullName: req.FullName,
		Hash:     hashedPassword,
	}

	user, err := server.userAccount.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.userAccount.UpdateTester(ctx, db.UpdateTesterParams{
		ID:     tester.ID,
		Email:  tester.Email,
		UserID: sql.NullInt64{Int64: user.ID, Valid: true},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response, done := server.loginValidatedUser(ctx, user)
	if done {
		return // TODO refactor this ugly thing
	}

	ctx.JSON(http.StatusOK, response)
}

type loginUserRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionId             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) { // todo kerok - add tests for login api endpoint
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.userAccount.GetUserByUserName(ctx, req.UserName)
	if err != nil {
		user, err = server.userAccount.GetUserByUserEmail(ctx, req.UserName)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = util.CheckPassword(req.Password, user.Hash)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	response, done := server.loginValidatedUser(ctx, user)
	if done {
		return
	}
	ctx.JSON(http.StatusOK, response)

}

func (server Server) loginValidatedUser(ctx *gin.Context, user db.User) (loginUserResponse, bool) {
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.UserName,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return loginUserResponse{}, true
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.UserName, server.config.RefreshTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return loginUserResponse{}, true
	}
	// TODO add access_level from stripe_plans table, either to session or straight to user
	userSession, err := server.userAccount.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserName:     user.UserName,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return loginUserResponse{}, true
	}

	response := loginUserResponse{
		SessionId:             userSession.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  newUserResponse(user),
	}
	return response, false
}
