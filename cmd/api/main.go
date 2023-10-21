package main

import (
	"fmt"
	"github.com/aerosystems/user-service/internal/handlers"
	"github.com/aerosystems/user-service/internal/middleware"
	"github.com/aerosystems/user-service/internal/models"
	"github.com/aerosystems/user-service/internal/repository"
	RPCServer "github.com/aerosystems/user-service/internal/rpc_server"
	"github.com/aerosystems/user-service/internal/services"
	GormPostgres "github.com/aerosystems/user-service/pkg/gorm_postgres"
	"github.com/aerosystems/user-service/pkg/logger"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"os"
)

const (
	webPort = 80
	rpcPort = 5001
)

// @title User Service
// @version 1.0.0
// @description A part of microservice infrastructure, who responsible for user entity.

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host gw.verifire.com/user
// @schemes https
// @BasePath /
func main() {
	log := logger.NewLogger(os.Getenv("HOSTNAME"))

	clientGORM := GormPostgres.NewClient(logrus.NewEntry(log.Logger))
	_ = clientGORM.AutoMigrate(models.User{})

	userRepo := repository.NewUserRepo(clientGORM)

	userService := services.NewUserServiceImpl(userRepo)

	baseHandler := handlers.NewBaseHandler(userService)
	rpcServer := RPCServer.NewUserServer(rpcPort, log.Logger, userService)

	app := Config{
		BaseHandler: baseHandler,
	}

	e := app.NewRouter()
	middleware.AddMiddleware(e, log.Logger)

	errChan := make(chan error)

	go func() {
		log.Infof("starting user-service HTTP server on port %d\n", webPort)
		errChan <- e.Start(fmt.Sprintf(":%d", webPort))
	}()

	go func() {
		log.Infof("starting user-service RPC server on port %d\n", rpcPort)
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
