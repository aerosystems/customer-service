package main

import "github.com/aerosystems/customer-service/cmd/app"

// @title Customer Service
// @version 1.0.1
// @description A part of microservice infrastructure, who responsible for customer user entity.

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host gw.verifire.dev/customer
// @schemes https
// @BasePath /
func main() {
	app.Execute()
}
