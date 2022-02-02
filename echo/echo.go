package echo

import (
	"time"

	"github.com/fiam/reco"
	"github.com/labstack/echo/v4"
)

func Middleware(renderers ...reco.Renderer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			r = r.WithContext(reco.SetDebug(r.Context(), c.Echo().Debug))
			r = r.WithContext(reco.SetStarted(r.Context(), time.Now()))
			w := c.Response().Writer
			defer reco.HTTPReco(w, r, renderers...)
			return next(c)
		}
	}
}
