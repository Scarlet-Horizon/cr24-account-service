package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/db"
	"main/model"
	"main/request"
	"main/response"
	"main/util"
	"net/http"
	"strings"
	"time"
)

type AccountController struct {
	DB *db.AccountDB
}

// Create godoc
//
//	@Description	Create a new account for user.
//	@Summary		Create a new account for user
//	@Accept			json
//	@Produce		json
//	@Tags			account
//	@Param			requestBody	body		request.AccountRequest	true	"Account type"
//	@Success		201			{object}	model.Account
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/account [POST]
func (receiver AccountController) Create(context *gin.Context) {
	var req request.AccountRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	var limit int
	var ok bool
	if limit, ok = util.AccountTypesLimit[req.Type]; !ok {
		err := context.Error(errors.New("invalid account type. Supported options are: 'checking', 'saving'"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	bankAccount := model.Account{
		PK:       util.GetPK(context.MustGet("ID").(string)),
		SK:       util.GetSK(uuid.NewString()),
		Amount:   0,
		Limit:    limit,
		OpenDate: time.Now(),
		Type:     req.Type,
	}

	err := receiver.DB.Create(bankAccount)
	if err != nil {
		_ = context.Error(err)
		if errors.Is(err, util.AlreadyExists) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	//util.UploadAccount(bankAccount, context)
	context.JSON(http.StatusCreated, bankAccount)
}

// GetAll godoc
//
//	@Description	Get accounts for a specific user.
//	@Summary		Get accounts for a specific user
//	@Produce		json
//	@Tags			account
//	@Param			type	path		string			true	"What accounts to get: 'open', 'closed', 'all'"
//	@Success		200		{object}	[]model.Account	"An array of Account's"
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/accounts/{type} [GET]
func (receiver AccountController) GetAll(context *gin.Context) {
	acc := receiver.get(context)
	if acc != nil {
		context.JSON(http.StatusOK, acc)
	}
}

func (receiver AccountController) depositWithdraw(context *gin.Context, deposit bool) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		err := context.Error(errors.New("invalid account id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	var req request.MonetaryRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Amount < 1 {
		err := context.Error(errors.New("invalid amount, minimum is 1"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	if deposit {
		err := receiver.DB.Deposit(bankAccount, req.Amount)
		if err != nil {
			_ = context.Error(err)
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
			return
		}
	} else {
		err := receiver.DB.Withdraw(bankAccount, req.Amount)
		if err != nil {
			_ = context.Error(err)
			if errors.Is(err, util.InsufficientFounds) {
				context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
				return
			}
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
			return
		}
	}
	context.Status(http.StatusNoContent)
}

// Deposit godoc
//
//	@Description	Deposit money to a specific account.
//	@Summary		Deposit money to a specific account
//	@Tags			account
//	@Param			accountID	path	string					true	"Account ID"
//	@Param			requestBody	body	request.MonetaryRequest	true	"Amount to deposit"
//	@Success		204			"No Content"
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/account/{accountID}/deposit [PATCH]
func (receiver AccountController) Deposit(context *gin.Context) {
	receiver.depositWithdraw(context, true)
}

// Withdraw godoc
//
//	@Description	Withdraw money from a specific account.
//	@Summary		Withdraw money from a specific account
//	@Tags			account
//	@Param			accountID	path	string					true	"Account ID"
//	@Param			requestBody	body	request.MonetaryRequest	true	"Amount to withdraw"
//	@Success		204			"No Content"
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/account/{accountID}/withdraw [PATCH]
func (receiver AccountController) Withdraw(context *gin.Context) {
	receiver.depositWithdraw(context, false)
}

// Close godoc
//
//	@Description	Close a specific account.
//	@Summary		Close a specific account
//	@Tags			account
//	@Param			accountID	path	string	true	"Account ID"
//	@Success		204			"No Content"
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/account/{accountID}/close [PATCH]
func (receiver AccountController) Close(context *gin.Context) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		err := context.Error(errors.New("invalid account id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	err := receiver.DB.Close(bankAccount)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}

// Delete godoc
//
//	@Description	Delete a specific account.
//	@Summary		Delete a specific account
//	@Tags			account
//	@Param			accountID	path	string	true	"Account ID"
//	@Success		204			"No Content"
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/account/{accountID} [DELETE]
func (receiver AccountController) Delete(context *gin.Context) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		err := context.Error(errors.New("invalid account id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	err := receiver.DB.Delete(bankAccount)
	if err != nil {
		_ = context.Error(err)
		if errors.Is(err, util.InvalidAccount) || errors.Is(err, util.OpenAccount) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}

// GetAccount godoc
//
//	@Description	Get a specific account.
//	@Summary		Get a specific account
//	@Produce		json
//	@Tags			account
//	@Param			accountID	path		string	true	"Account ID"
//	@Success		200			{object}	model.Account
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/account/{accountID} [GET]
func (receiver AccountController) GetAccount(context *gin.Context) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		err := context.Error(errors.New("invalid account id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	acc, err := receiver.DB.GetAccount(bankAccount)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if acc.PK == "" {
		err := context.Error(util.InvalidAccount)
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, acc)
}

func (receiver AccountController) get(context *gin.Context) []model.Account {
	t := context.Param("type")
	if !(t == "open" || t == "closed" || t == "all") {
		err := context.Error(errors.New("invalid type, supported: 'open', 'closed', all"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return nil
	}

	acc, err := receiver.DB.GetAll(context.MustGet("ID").(string), t)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return nil
	}

	if len(acc) == 0 {
		context.Status(http.StatusNoContent)
		return nil
	}
	return acc
}

// GetAllWithTransactions godoc
//
//	@Description	Get all accounts with transactions for a given user.
//	@Summary		Get all accounts with transactions for a given user
//	@Produce		json
//	@Tags			account
//	@Param			type	path		string	true	"What accounts to get: 'open', 'closed', 'all'"
//	@Success		200		{object}	[]model.Account
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Security		JWT
//	@Param			Authorization	header	string	true	"Authorization"
//	@Router			/accounts/{type}/transactions [GET]
func (receiver AccountController) GetAllWithTransactions(context *gin.Context) {
	acc := receiver.get(context)
	if acc == nil {
		return
	}

	var accTr []model.Account
	for _, v := range acc {
		tr, err := util.GetTransactions(strings.Split(v.SK, "#")[1], context.MustGet("token").(string),
			context.GetString("Correlation"))

		if err != nil {
			_ = context.Error(err)
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
			return
		}

		account := v
		account.Transactions = tr
		accTr = append(accTr, account)
	}
	context.JSON(http.StatusOK, accTr)
}
