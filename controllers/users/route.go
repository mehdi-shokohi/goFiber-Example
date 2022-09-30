package usersController

import (
	"github.com/gofiber/fiber/v2"

	goexJWT "goex/middlewares/jwt"
)

import  "goex/middlewares/routerGroups"
const RouteContext = "/user"

func RouteDecision(api fiber.Router) {
	api.Post("/user/login", UserLoginHandler)
	apiGroup := api.Group(RouteContext,)
	apiGroup.Use(userGroups.SendErrorUserGroups(),  goexJWT.New())
	apiGroup.Get("/data", Getuserdata)
	apiGroup.Post("/", AddNewUser)
	apiGroup.Delete("/", RemoveUser)

}
