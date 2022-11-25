package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/db"
	"main/model"
	"main/request"
	"main/util"
	"net/http"
	"time"
)

type AccountController struct {
	DB *db.AccountDB
}

func (receiver AccountController) Create(c *gin.Context) {
	var req request.CreateAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !util.IsValidUUID(req.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var limit int
	var ok bool
	if limit, ok = util.AccountTypesLimit[req.Type]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account type. Supported options are: 'checking', 'saving'"})
		return
	}

	bankAccount := model.Account{
		PK:       req.UserID,
		SK:       "ACCOUNT#" + uuid.NewString(),
		Amount:   0,
		Limit:    limit,
		OpenDate: time.Now(),
		Type:     req.Type,
	}
	c.JSON(http.StatusCreated, bankAccount)
}
