package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"io"
	"log"
	"main/controller"
	"main/db"
	_ "main/docs"
	"main/env"
	"main/util"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@title			cr24 Account API
//	@version		1.0
//	@description	API for account management for cr24 project
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	David Slatinek
//	@contact.url	https://github.com/david-slatinek

//	@accept		json
//	@produce	json
//	@schemes	http

//	@license.name	GNU General Public License v3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.html

//	@securityDefinitions.apikey	JWT
//@in header
//@name Authorization

//	@host		localhost:8080
//	@BasePath	/api/v1
func main() {
	rand.Seed(time.Now().UnixNano())

	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env variables: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Fatalf("failed to load SDK config, %s", err)
	}

	accountController := controller.AccountController{
		DB: &db.AccountDB{
			Client: dynamodb.NewFromConfig(cfg),
		},
	}

	gin.SetMode(os.Getenv("GIN_MODE"))

	if gin.Mode() == gin.ReleaseMode {
		gin.DisableConsoleColor()

		f, err := os.Create("gin.log")
		if err != nil {
			gin.DefaultWriter = io.MultiWriter(f)
		}
	}

	router := gin.Default()
	router.Use(util.CORS).Use(util.Info)
	api := router.Group("api/v1").Use(util.ValidateToken)
	{
		api.POST("/account", accountController.Create)

		api.GET("/accounts/:type", accountController.GetAll)
		api.GET("/accounts/:type/transactions", accountController.GetAllWithTransactions)
		api.GET("/account/:accountID", accountController.GetAccount)

		api.PATCH("/account/:accountID/deposit", accountController.Deposit)
		api.PATCH("/account/:accountID/withdraw", accountController.Withdraw)
		api.PATCH("/account/:accountID/close", accountController.Close)

		api.DELETE("/account/:accountID", accountController.Delete)
	}
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Println("server is up at: " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ListenAndServe() error: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown() error: %s\n", err)
	}

	log.Println("shutting down")
	os.Exit(0)
}
