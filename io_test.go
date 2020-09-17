// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package fiber

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectBodies = []string{
	`{"code":200,"data":{"param":"john"}}`,
	`{"code":200,"data":{"page":1,"size":20}}`,
}

func AppTestReq(app *App, method, url string, body io.Reader) (resp *http.Response, err error) {
	return app.Test(httptest.NewRequest(method, url, body))
}

func CreateTestApp() *App {
	app := New()
	route := app.Group("/test")
	route.Get("/:param", func(c *Ctx) (err error) {
		param := c.Params("param")
		data := Map{"param": param}
		c.JSON(Map{"code": 200, "data": data})
		return
	})
	route.Get("/page/:page/:size", func(c *Ctx) (err error) {
		page, size := 1, 20
		if pageStr := c.FormValue("page"); pageStr != "" {
			page, err = strconv.Atoi(pageStr)
		}
		if sizeStr := c.FormValue("size"); sizeStr != "" {
			size, err = strconv.Atoi(sizeStr)
		}
		data := Map{"page": page, "size": size}
		c.JSON(Map{"code": 200, "data": data})
		return
	})
	return app
}

//func CreateSeniorApp() *App {
//	app := New()
//	route := app.Group("/test")
//	route.Get("/:param", func(c *Ctx) (err error) {
//		param := c.ParamStr("param")
//		c.Reply(Map{"param": param})
//		return
//	})
//	route.Get("/page/:page/:size", func(c *Ctx) (err error) {
//		page, size := c.FetchInt("page", 1), c.FetchInt("size", 20)
//		c.Reply(Map{"page": page, "size": size})
//		return
//	})
//	return app
//}

func Test_Request01_ParamStr(t *testing.T) {
	app := CreateTestApp()
	resp, err := AppTestReq(app, "GET", "/test/john", nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	assert.Greater(t, n, 0)
	assert.Equal(t, string(body[:n]), expectBodies[0])
}

func Test_Request02_FetchInt(t *testing.T) {
	app := CreateTestApp()
	resp, err := AppTestReq(app, "GET", "/test/page/3/7", nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	assert.Greater(t, n, 0)
	assert.Equal(t, string(body[:n]), expectBodies[1])
}
