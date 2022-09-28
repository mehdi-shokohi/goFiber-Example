package usersController

import (
	"github.com/gofiber/fiber/v2"

	goexJWT "goex/middlewares/jwt"
)

const RouteContext = "/user"

func RouteDecision(api fiber.Router) {
	apiGroup := api.Group(RouteContext)
	apiGroup.Post("/login", UserLoginHandler)
	apiGroup.Get("/data", goexJWT.New(), Getuserdata)
	apiGroup.Post("/", goexJWT.New(), AddNewUser)
	apiGroup.Delete("/",goexJWT.New(), RemoveUser)

}
