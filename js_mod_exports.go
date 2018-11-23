package main

import (
	"strings"

	"github.com/labstack/echo"
)

type jsExports struct {
	DocType string
	Status  int
	Headers map[string]string
	Body    interface{}
}

func (x *jsExports) echoify(c echo.Context) error {
	for k, v := range x.Headers {
		c.Response().Header().Add(k, v)
	}
	switch strings.ToLower(x.DocType) {
	case "json":
		return c.JSON(x.Status, x.Body)
	case "html":
		return c.HTML(x.Status, x.Body.(string))
	}

	return c.JSON(500, map[string]interface{}{
		"success": false,
		"error":   "invalid exports.DocType specified",
	})
}
