package main

import (
	json "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	 "github.com/swaggo/fiber-swagger" 

	conf "goex/config"
	"goex/controllers"
	_ "goex/docs"
)
// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @host localhost:3000
// @BasePath /
// @schemes http

func main() {
	app := fiber.New(fiber.Config{
		Prefork:     true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(logger.New())
	controllers.RegisterAll(app)
	app.Get("/swagger/*", fiberSwagger.WrapHandler) // default


	app.Listen(conf.GetPort())
}
