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
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var limit int
	var ok bool
	if limit, ok = util.AccountTypesLimit[req.Type]; !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "invalid account type. Supported options are: 'checking', 'saving'",
		})
	}

	bankAccount := model.Account{
		PK:       req.UserID,
		SK:       "ACCOUNT#" + uuid.NewString(),
		Amount:   0,
		Limit:    limit,
		OpenDate: time.Now(),
		Type:     req.Type,
	}
	c.IndentedJSON(http.StatusCreated, bankAccount)
}
