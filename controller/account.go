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

//	@description	Create new account for a user.
//	@summary		Create new account for a user
//	@accept			json
//	@produce		json
//	@tags			account
//	@param			requestBody	body		request.AccountRequest	true	"Account type"
//	@success		201			{object}	model.Account
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account [POST]
func (receiver AccountController) Create(context *gin.Context) {
	var req request.AccountRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	var limit int
	var ok bool
	if limit, ok = util.AccountTypesLimit[req.Type]; !ok {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account type. Supported options are: 'checking', 'saving'"})
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
		if errors.Is(err, util.AlreadyExists) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.JSON(http.StatusCreated, bankAccount)
}

//	@description	Get accounts for a specific user.
//	@summary		Get accounts for a specific user
//	@produce		json
//	@tags			account
//	@param			type	path		string			true	"What accounts to get: 'open', 'closed', 'all'"
//	@success		200		{object}	[]model.Account	"An array of Account's"
//	@failure		400		{object}	response.ErrorResponse
//	@failure		500		{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/accounts/{type} [GET]
func (receiver AccountController) GetAll(context *gin.Context) {
	acc := receiver.get(context)
	if acc != nil {
		context.JSON(http.StatusOK, acc)
	}
}

func (receiver AccountController) depositWithdraw(context *gin.Context, deposit bool) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account id"})
		return
	}

	var req request.MonetaryRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Amount < 1 {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid amount, minimum is 1"})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	if deposit {
		err := receiver.DB.Deposit(bankAccount, req.Amount)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
			return
		}
	} else {
		err := receiver.DB.Withdraw(bankAccount, req.Amount)
		if err != nil {
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

//	@description	Deposit money to a specific account.
//	@summary		Deposit money to a specific account
//	@tags			account
//	@param			accountID	path	string					true	"Account ID"
//	@param			requestBody	body	request.MonetaryRequest	true	"Amount to deposit"
//	@success		204			"No Content"
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account/{accountID}/deposit [PATCH]
func (receiver AccountController) Deposit(context *gin.Context) {
	receiver.depositWithdraw(context, true)
}

//	@description	Withdraw money from a specific account.
//	@summary		Withdraw money from a specific account
//	@tags			account
//	@param			accountID	path	string					true	"Account ID"
//	@param			requestBody	body	request.MonetaryRequest	true	"Amount to withdraw"
//	@success		204			"No Content"
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account/{accountID}/withdraw [PATCH]
func (receiver AccountController) Withdraw(context *gin.Context) {
	receiver.depositWithdraw(context, false)
}

//	@description	Close a specific account.
//	@summary		Close a specific account
//	@tags			account
//	@param			accountID	path	string	true	"Account ID"
//	@success		204			"No Content"
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account/{accountID}/close [PATCH]
func (receiver AccountController) Close(context *gin.Context) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account id"})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	err := receiver.DB.Close(bankAccount)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}

//	@description	Delete a specific account.
//	@summary		Delete a specific account
//	@tags			account
//	@param			accountID	path	string	true	"Account ID"
//	@success		204			"No Content"
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account/{accountID} [DELETE]
func (receiver AccountController) Delete(context *gin.Context) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account id"})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	err := receiver.DB.Delete(bankAccount)
	if err != nil {
		if errors.Is(err, util.InvalidAccount) || errors.Is(err, util.OpenAccount) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}

//	@description	Get a specific account.
//	@summary		Get a specific account
//	@produce		json
//	@tags			account
//	@param			accountID	path		string	true	"Account ID"
//	@success		200			{object}	model.Account
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account/{accountID} [GET]
func (receiver AccountController) GetAccount(context *gin.Context) {
	accountID := context.Param("accountID")
	if !util.IsValidUUID(accountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account id"})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(context.MustGet("ID").(string)),
		SK: util.GetSK(accountID),
	}

	acc, err := receiver.DB.GetAccount(bankAccount)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if acc.PK == "" {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: util.InvalidAccount.Error()})
		return
	}
	context.JSON(http.StatusOK, acc)
}

func (receiver AccountController) get(context *gin.Context) []model.Account {
	t := context.Param("type")
	if !(t == "open" || t == "closed" || t == "all") {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid type, supported: 'open', 'closed', all"})
		return nil
	}

	acc, err := receiver.DB.GetAll(context.MustGet("ID").(string), t)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return nil
	}

	if len(acc) == 0 {
		context.Status(http.StatusNoContent)
		return nil
	}
	return acc
}

//	@description	Get all accounts with transactions for a given user.
//	@summary		Get all accounts with transactions for a given user
//	@produce		json
//	@tags			account
//	@param			type	path		string	true	"What accounts to get: 'open', 'closed', 'all'"
//	@success		200		{object}	[]model.Account
//	@failure		400		{object}	response.ErrorResponse
//	@failure		500		{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/account/{type}/transactions [GET]
func (receiver AccountController) GetAllWithTransactions(context *gin.Context) {
	acc := receiver.get(context)
	if acc == nil {
		return
	}

	var accTr []model.Account
	for _, v := range acc {
		tr, err := util.GetTransactions(strings.Split(v.SK, "#")[1], context.MustGet("token").(string))

		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
			return
		}

		account := v
		account.Transactions = tr
		accTr = append(accTr, account)
	}
	context.JSON(http.StatusOK, accTr)
}
