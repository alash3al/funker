package main

import (
	"github.com/labstack/echo"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if (*flagAuthKey != "" && c.Request().Header.Get("Authorization") == *flagAuthKey) || (*flagAuthKey == "") {
			return next(c)
		}
		return c.JSON(401, map[string]interface{}{
			"success": false,
			"error":   "authentication required",
		})
	}
}
