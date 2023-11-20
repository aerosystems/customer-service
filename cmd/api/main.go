package main

import (
	"fmt"
	"github.com/aerosystems/customer-service/internal/handlers"
	"github.com/aerosystems/customer-service/internal/middleware"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/aerosystems/customer-service/internal/repository"
	RPCServer "github.com/aerosystems/customer-service/internal/rpc_server"
	RPCServices "github.com/aerosystems/customer-service/internal/rpc_services"
	"github.com/aerosystems/customer-service/internal/services"
	GormPostgres "github.com/aerosystems/customer-service/pkg/gorm_postgres"
	"github.com/aerosystems/customer-service/pkg/logger"
	RPCClient "github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"os"
)

const (
	webPort = 80
	rpcPort = 5001
)

// @title Customer Service
// @version 1.0.0
// @description A part of microservice infrastructure, who responsible for customer user entity.

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host gw.verifire.com/customer
// @schemes https
// @BasePath /
func main() {
	log := logger.NewLogger(os.Getenv("HOSTNAME"))

	clientGORM := GormPostgres.NewClient(logrus.NewEntry(log.Logger))
	if err := clientGORM.AutoMigrate(models.Customer{}); err != nil {
		log.Fatal(err)
	}

	customerRepo := repository.NewCustomerRepo(clientGORM)

	projectRPCClient := RPCClient.NewClient("tcp", "project-service:5001")
	projectRPC := RPCServices.NewProjectRPC(projectRPCClient)

	subsRPCClient := RPCClient.NewClient("tcp", "subs-service:5001")
	subsRPC := RPCServices.NewSubsRPC(subsRPCClient)

	userService := services.NewCustomerServiceImpl(customerRepo, projectRPC, subsRPC)

	baseHandler := handlers.NewBaseHandler(os.Getenv("APP_ENV"), log.Logger, userService)

	rpcServer := RPCServer.NewCustomerServer(rpcPort, log.Logger, userService)

	accessTokenService := services.NewAccessTokenServiceImpl(os.Getenv("ACCESS_SECRET"))

	oauthMiddleware := middleware.NewOAuthMiddlewareImpl(accessTokenService)
	basicAuthMiddleware := middleware.NewBasicAuthMiddlewareImpl(os.Getenv("BASIC_AUTH_DOCS_USERNAME"), os.Getenv("BASIC_AUTH_DOCS_PASSWORD"))

	app := NewApp(baseHandler, oauthMiddleware, basicAuthMiddleware)
	e := app.NewRouter()
	middleware.AddLog(e, log.Logger)
	middleware.AddCORS(e)

	errChan := make(chan error)

	go func() {
		log.Infof("starting customer-service HTTP server on port %d\n", webPort)
		errChan <- e.Start(fmt.Sprintf(":%d", webPort))
	}()

	go func() {
		log.Infof("starting customer-service RPC server on port %d\n", rpcPort)
		errChan <- rpc.Register(rpcServer)
		errChan <- rpcServer.Listen(rpcPort)
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
