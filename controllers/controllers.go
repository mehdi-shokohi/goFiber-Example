package controllers

import (
	"github.com/gofiber/fiber/v2"

	usersController "goex/controllers/users"
)

func RegisterAll(api fiber.Router) {

	usersController.RouteDecision(api)

}
