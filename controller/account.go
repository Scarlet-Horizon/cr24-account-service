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
//	@param			requestBody	body		request.CreateAccountRequest	true	"User ID and account type"
//	@success		201			{object}	model.Account
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@router			/account [POST]
func (receiver AccountController) Create(c *gin.Context) {
	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUUID(req.UserID) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user id"})
		return
	}

	var limit int
	var ok bool
	if limit, ok = util.AccountTypesLimit[req.Type]; !ok {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account type. Supported options are: 'checking', 'saving'"})
		return
	}

	bankAccount := model.Account{
		PK:       util.GetPK(req.UserID),
		SK:       util.GetSK(uuid.NewString()),
		Amount:   0,
		Limit:    limit,
		OpenDate: time.Now(),
		Type:     req.Type,
	}

	err := receiver.DB.Create(bankAccount)
	if err != nil {
		if errors.Is(err, util.AlreadyExists) {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bankAccount)
}

//	@description	Get all accounts for a specific user.
//	@summary		Get all accounts for a specific user
//	@produce		json
//	@tags			account
//	@param			userID	path		string			true	"User ID"
//	@success		200		{object}	[]model.Account	"An array of Account's"
//	@failure		400		{object}	response.ErrorResponse
//	@failure		500		{object}	response.ErrorResponse
//	@router			/accounts/{userID} [GET]
func (receiver AccountController) GetAll(c *gin.Context) {
	userID := c.Param("userID")

	if !util.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid user id"})
		return
	}

	acc, err := receiver.DB.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if len(acc) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, acc)
}

//	@description	Deposit money to a specific account.
//	@summary		Deposit money to a specific account
//	@tags			account
//	@param			accountID	path	string					true	"Account ID"
//	@param			requestBody	body	request.MonetaryRequest	true	"User ID and amount to deposit"
//	@success		204			""
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@router			/account/{accountID}/deposit [PATCH]
func (receiver AccountController) Deposit(c *gin.Context) {
	accountID := c.Param("accountID")

	if !util.IsValidUUID(accountID) {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account id"})
		return
	}

	var req request.MonetaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Amount < 1 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid amount, minimum is 1"})
		return
	}

	bankAccount := model.Account{
		PK: util.GetPK(req.UserID),
		SK: util.GetSK(accountID),
	}

	err := receiver.DB.Deposit(bankAccount, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
