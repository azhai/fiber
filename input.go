// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package fiber

import (
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
)

const (
	COOKIE_TOKEN_KEY = "access_token"
	HEADER_TOKEN_KEY = "x-token"
)

func PeekForm(key string, mf *multipart.Form, err error) string {
	if err == nil && mf.Value != nil {
		vv := mf.Value[key]
		if len(vv) > 0 {
			return vv[0]
		}
	}
	return ""
}

func FirstStrArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func FirstIntArg(args []int) int {
	if len(args) > 0 {
		return args[0]
	}
	return 0
}

func FirstFloatArg(args []float64) float64 {
	if len(args) > 0 {
		return args[0]
	}
	return 0.0
}

func (c *Ctx) Token() (token string) {
	if token = c.CookieStr(COOKIE_TOKEN_KEY); token == "" {
		token = c.HeaderStr(HEADER_TOKEN_KEY)
	}
	return
}

// Common method read data of GET/POST/PARAM/HEADER/COOKIE
func (c *Ctx) Read(key, val string, methods ...string) (bool, string) {
	var value string
	req := c.Request()
	for _, m := range methods {
		switch strings.ToUpper(m) {
		case "COOKIE", "COOKIES":
			value = getString(req.Header.Cookie(key))
		case "GET":
			value = getString(req.URI().QueryArgs().Peek(key))
		case "HEAD", "HEADER":
			value = getString(req.Header.Peek(key))
		case "PARAM", "PARAMS":
			value = c.Params(key)
		case "POST", "PUT":
			value = c.FormValue(key)
		}
		if value != "" {
			return true, value
		}
	}
	return false, val
}

func (c *Ctx) ReadInt(key string, val int, args ...string) (bool, int) {
	if has, value := c.Read(key, "", args...); has {
		if v, err := strconv.Atoi(value); err == nil {
			return true, v
		}
	}
	return false, val
}

func (c *Ctx) ReadFloat(key string, val float64, args ...string) (bool, float64) {
	if has, value := c.Read(key, "", args...); has {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return true, v
		}
	}
	return false, val
}

func (c *Ctx) Contains(key, expect string, args ...string) (has, fit bool) {
	var value string
	if has, value = c.Read(key, "", args...); has {
		fit = strings.Contains(value, expect)
	}
	return
}

func (c *Ctx) GetStr(key string, args ...string) string {
	_, value := c.Read(key, FirstStrArg(args), "GET")
	return value
}

func (c *Ctx) GetInt(key string, args ...int) int {
	if has, value := c.ReadInt(key, 0, "GET"); has {
		return value
	}
	return FirstIntArg(args)
}

func (c *Ctx) GetFloat(key string, args ...float64) float64 {
	if has, value := c.ReadFloat(key, 0, "GET"); has {
		return value
	}
	return FirstFloatArg(args)
}

func (c *Ctx) PostStr(key string, args ...string) string {
	_, value := c.Read(key, FirstStrArg(args), "POST")
	return value
}

func (c *Ctx) PostInt(key string, args ...int) int {
	if has, value := c.ReadInt(key, 0, "POST"); has {
		return value
	}
	return FirstIntArg(args)
}

func (c *Ctx) PostFloat(key string, args ...float64) float64 {
	if has, value := c.ReadFloat(key, 0, "POST"); has {
		return value
	}
	return FirstFloatArg(args)
}

func (c *Ctx) PostAll() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	values, err := url.ParseQuery(string(c.Body()))
	if err != nil {
		return data, err
	}
	for key, vals := range values {
		data[key] = strings.Join(vals, ",")
	}
	return data, nil
}

// Read the POST first, if empty then read GET
func (c *Ctx) FetchStr(key string, args ...string) string {
	_, value := c.Read(key, FirstStrArg(args), "POST", "GET")
	return value
}

func (c *Ctx) FetchInt(key string, args ...int) int {
	if has, value := c.ReadInt(key, 0, "POST", "GET"); has {
		return value
	}
	return FirstIntArg(args)
}

func (c *Ctx) FetchFloat(key string, args ...float64) float64 {
	if has, value := c.ReadFloat(key, 0, "POST", "GET"); has {
		return value
	}
	return FirstFloatArg(args)
}

func (c *Ctx) ParamStr(key string, args ...string) string {
	_, value := c.Read(key, FirstStrArg(args), "PARAM")
	return value
}

func (c *Ctx) ParamInt(key string, args ...int) int {
	if has, value := c.ReadInt(key, 0, "PARAM"); has {
		return value
	}
	return FirstIntArg(args)
}

func (c *Ctx) ParamFloat(key string, args ...float64) float64 {
	if has, value := c.ReadFloat(key, 0, "PARAM"); has {
		return value
	}
	return FirstFloatArg(args)
}

func (c *Ctx) HeaderStr(key string, args ...string) string {
	_, value := c.Read(key, FirstStrArg(args), "HEADER")
	return value
}

func (c *Ctx) CookieStr(key string, args ...string) string {
	_, value := c.Read(key, FirstStrArg(args), "COOKIE")
	return value
}
