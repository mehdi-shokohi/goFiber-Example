package main

import (
	"github.com/gofiber/fiber/v2"

	conf "goex/config"
	"goex/controllers"
)

func main() {
	app := fiber.New()

	controllers.RegisterAll(app)

	app.Listen(conf.GetPort())
}
