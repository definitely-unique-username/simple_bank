package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Ammount       int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := s.validAccount(ctx, req.FromAccountId, req.Currency)

	if !valid {
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)

	if authPayload.UserID != fromAccount.Owner {
		err := errors.New("from account not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))

		return
	}

	_, valid = s.validAccount(ctx, req.ToAccountId, req.Currency)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Ammount,
	}

	result, err := s.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (s *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
