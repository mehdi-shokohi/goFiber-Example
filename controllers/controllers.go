package controllers
import(
"github.com/gofiber/fiber/v2"
"goex/controllers/users"
)
func RegisterAll(api fiber.Router) {

	usersController.RouteDecision(api)

}