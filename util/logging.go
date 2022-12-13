package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
)

func logging(level string, context *gin.Context) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("time=\"%s\"", time.Now().Format("2006-01-02 15-04-05")))
	sb.WriteString(" id=" + uuid.NewString())
	sb.WriteString(" level=" + level)
	sb.WriteString(" path=" + context.Request.RequestURI)

	correlation := context.GetHeader("Correlation")
	if correlation == "" {
		correlation = "nil"
	}
	sb.WriteString(" correlation=" + correlation)
	sb.WriteString(" ip=" + context.ClientIP())

	var token string
	values := strings.Split(context.GetHeader("Authorization"), "Bearer ")

	if len(values) == 2 {
		token = values[1]

		to, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err == nil && to.Valid {
			if claims, ok := to.Claims.(jwt.MapClaims); ok {
				token = claims["sub"].(string)
			}
		}
	} else {
		token = "nil"
	}
	sb.WriteString(" auth=" + token)

	return sb.String()
}

func Info(context *gin.Context) string {
	return logging("info", context)
}

func Error(err string, context *gin.Context) string {
	return logging("error", context) + " msg=" + fmt.Sprintf("\"%s\"", err)
}
