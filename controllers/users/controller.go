package usersController

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"goex/config/messages"
	"goex/db/dbHelper"
	"goex/db/mongoHelper"
	"goex/entity/User"
	goexJWT "goex/middlewares/jwt"
)

// @Login and get Token
// @Login.
// @Tags Login
// @Produce json
// @Param LoginBody body User.UserLogin true "enter username and password ,eg. admin , 1234"
// @Success 200 {object} map[string]interface{}
// @Router /user/login [post]
func UserLoginHandler(c *fiber.Ctx) error {
	loginForm := new(User.UserLogin)

	if err := c.BodyParser(loginForm); err != nil {

		return c.JSON(fiber.Map{"data": messages.InvalidInputForm})
	}
	// userModel:=new(User.Model)
	db := mongoHelper.NewMongo(c.Context(), &User.Model{})
	_, err := db.FindOne(&bson.D{{Key: "username", Value: loginForm.Username}, {Key: "password", Value: loginForm.Password}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(fiber.Map{"data": "user not found "})
		} else {
			return c.JSON(fiber.Map{"data": err.Error()})
		}
	}

	userModel := db.Model.(*User.Model)
	claims := jwt.MapClaims{
		"name":   userModel.Username,
		"admin":  userModel.Admin,
		"expire": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := goexJWT.GetPrivateKey()
	if err != nil {
		panic(err)
	}
	t, err := token.SignedString(privateKey)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"Authorization": "Bearer " + t})

}

// Extract Jwt Session Data godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Get User Data
// @Accept */*
// @Security Bearer
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /user/data [get]
func Getuserdata(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// panic(errors.New("error"))  // For Error Handler On RouterGroup
	return c.JSON(fiber.Map{"data": claims})
}
// @Add new User
// @Add New User.
// @Tags Add user
// @Produce json
// @Security Bearer
// @Param RegisterNewUser body User.RegisterForm true "Enter New User Data"
// @Success 200 {object} map[string]interface{}
// @Router /user [post]
func AddNewUser(c *fiber.Ctx) error {
	inputForm := new(User.RegisterForm)
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if admin, ok := claims["admin"].(bool); ok {
		if admin == false {
			return c.JSON(fiber.Map{"data": messages.HAVNTGRANT})
		}
		if err := c.BodyParser(inputForm); err != nil {

			return c.JSON(fiber.Map{"data": messages.InvalidInputForm})
		}
	}
	findedUser := new(User.Model)
	r := dbHelper.FindOneGo(c.Context(), &bson.D{{Key: "username", Value: inputForm.Username}}, findedUser)
	if err := <-r; err != nil {
		fmt.Print(err)
	}
	if findedUser.Username != "" {
		return c.SendString("User Exists")
	}
	userModel := new(User.Model)
	userModel.Username = inputForm.Username
	userModel.Password = inputForm.Password
	userModel.Admin = inputForm.Admin
	userModel.FirstName = inputForm.FirstName
	userModel.LastName = inputForm.LastName
	result := dbHelper.SaveGo(c.Context(), userModel)
	if err := <-result; err != nil {
		return c.JSON(fiber.Map{"data": "error in savin record"})
	}

	return c.JSON(fiber.Map{"data": userModel})
}

func RemoveUser(c *fiber.Ctx) error {
	inputForm := new(User.RegisterForm)
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if admin, ok := claims["admin"].(bool); ok {
		if !admin {
			return c.JSON(fiber.Map{"data": messages.HAVNTGRANT})
		}
		if err := c.BodyParser(inputForm); err != nil {

			return c.JSON(fiber.Map{"data": messages.InvalidInputForm})
		}
	}
	var deletedCount int64
	result := dbHelper.DeleteManyGo(c.Context(), "users", &bson.D{{"username", inputForm.Username}}, &deletedCount)
	if err := <-result; err != nil {
		return c.JSON(fiber.Map{"data": "error occurred"})

	}
	return c.JSON(fiber.Map{"data": deletedCount})
}
