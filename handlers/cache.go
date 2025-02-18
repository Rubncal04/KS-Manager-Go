package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/labstack/echo/v4"
)

type HandlerCache struct {
	Cache db.RedisRepo
}

func NewCache(cache *db.RedisRepo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			cacheKey := c.Request().RequestURI
			method := c.Request().Method

			if method == "POST" || method == "PUT" || method == "DELETE" {
				cache.DeleteCache(cacheKey)
			} else {
				cachedData, err := cache.GetCache(cacheKey)
				if err == nil && cachedData != "" {
					var response interface{}
					if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
						c.Response().Header().Set("X-Cache-Status", "HIT")
						return c.JSON(http.StatusOK, response)
					}
				}
			}

			recorder := newResponseRecorder(c.Response())
			c.Response().Writer = recorder

			if err := next(c); err != nil {
				c.Error(err)
				return err
			}

			if recorder.statusCode == http.StatusOK && (method != "POST" && method != "PUT" && method != "DELETE") {
				_, err := cache.SetCache(cacheKey, recorder.body.String())
				if err != nil {
					log.Printf("Error storing data in cache: %v", err)
				}
			}

			c.Response().Header().Set("X-Cache-Status", "MISS")
			return nil
		}
	}
}

type ResponseRecorder struct {
	echo.Response
	body       *strings.Builder
	statusCode int
}

func newResponseRecorder(res *echo.Response) *ResponseRecorder {
	return &ResponseRecorder{
		Response:   *res,
		body:       &strings.Builder{},
		statusCode: res.Status,
	}
}

func (r *ResponseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.Response.WriteHeader(code)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.Response.Write(b)
}
