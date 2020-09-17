// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package fiber

import (
	"fmt"
	"net/http"
)

func GetServiceCode(statusCode int) int {
	if statusCode == 200 {
		return 0
	}
	return statusCode
}

// Give the difference of reponse body or reponse data
type ReplyBody Map

// Type sets the Content-Type HTTP header to the MIME type specified by the file extension.
func (c *Ctx) SetType(extension string, charset ...string) *Ctx {
	return c.Type(extension, charset...)
}

// Status sets the HTTP status for the response.
func (c *Ctx) SetStatus(status int) *Ctx {
	return c.Status(status)
}

// Send formatted string
func (c *Ctx) Printf(format string, args ...interface{}) error {
	c.SendString(fmt.Sprintf(format, args...))
	return nil
}

func (c *Ctx) Jsonify(format string, args ...interface{}) error {
	c.SetType("json")
	return c.Printf(format, args...)
}

func (c *Ctx) Errorf(servCode int, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	body := Map{"code": servCode, "message": msg}
	c.SetStatus(http.StatusOK)
	return c.JSON(body)
}

func (c *Ctx) Abort(code int, data interface{}) error {
	c.SetStatus(code)
	if data != nil {
		return c.JSON(data)
	}
	return nil
}

func (c *Ctx) Deny(msg string) error {
	code := http.StatusForbidden
	servCode := GetServiceCode(code)
	return c.Abort(code, Map{"code": servCode, "message": msg})
}

func (c *Ctx) Reply(data interface{}, metas ...int64) error {
	var body Map
	if rbody, ok := data.(ReplyBody); ok { // Map as response body
		body = Map(rbody)
	} else { // Map as data in response body
		servCode := GetServiceCode(http.StatusOK)
		body = Map{"code": servCode, "data": data}
	}
	if len(metas) >= 1 {
		body["total"] = metas[0]
	}
	return c.JSON(body)
}
