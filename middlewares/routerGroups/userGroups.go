package userGroups

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"

	conf "goex/config"
)

var poolSender = sync.Pool{New: func() interface{} {
	return new(Send)
}}

type Send struct {
	C            *fiber.Ctx
	Status       bool
	ErrorMessage *Error
	ResponseCode int
	Data         interface{}
}
type SendResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Error  *Error      `json:"error"`
}

func (provider *Send) Send() {
	response := new(SendResponse)

	provider.C.Response().Header.Set("Content-Type", "application/json")
	//if provider.ResponseCode != 0 {
	provider.C.Response().Header.SetStatusCode(200)
	//}
	//else{
	//	if provider.Status==true {
	//		provider.C.Response.Header.SetStatusCode(200)
	//		response.Status = true
	//	}
	//}
	response.Status = provider.Status
	response.Error = provider.ErrorMessage
	response.Data = provider.Data

	v, err := jsoniter.Marshal(response)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprint(provider.C, string(v))
	if err != nil {
		fmt.Println("Problem In Response - Internal Error")
	}

}

type Error struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func releasePoolSender(response *Send) {
	Clear(response)
	poolSender.Put(response)
}
func Sender(connection *fiber.Ctx, status int, responseCode int, error *Error, data interface{}) {
	response := poolSender.Get().(*Send)
	defer releasePoolSender(response)
	switch status {
	case 1:
		response.Status = true
	case 2:
		response.Status = false

	case 3:
		response.Status = false

	}

	response.ResponseCode = 200
	if error != nil {
		response.ErrorMessage = error
	}
	response.C = connection
	response.Data = data
	response.Send()
	//connection.Abort()
}
func SendErrorUserGroups() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {

				fmt.Println(r)

				if e, ok := r.(error); ok {
					Sender(c, conf.InterErrorCode, 200, &Error{
						Code:    400,
						Message: e.Error(),
						Data:    nil,
					}, nil)

				} else if v, ok := r.(map[string]error); ok {
					t := ""
					for eKey, errMessage := range v {
						t += eKey + errMessage.Error() + " , "
					}
					Sender(c, conf.InterErrorCode, 200, &Error{
						Code:    400,
						Message: t,
						Data:    nil,
					}, nil)
				}
			}

		}()
		c.Next()
		return nil
	}

}
