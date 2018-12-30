package main

import (
	"fmt"
	"io/ioutil"

	"github.com/rs/xid"

	"github.com/dop251/goja"

	"github.com/labstack/echo"
)

// routeAddFunk - adding a funk to the registry
func routeAddFunk(c echo.Context) error {
	allBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil || len(allBytes) < 1 {
		return c.JSON(500, map[string]interface{}{
			"success": false,
			"error":   "empty request body",
		})
	}

	// cacheTTL, _ := strconv.Atoi(c.QueryParam("cache"))
	code := fmt.Sprintf("(%s)", string(allBytes))

	if err := funker.AddFunk(c.Param("funkName"), code, false); err != nil {
		return c.JSON(500, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(200, map[string]interface{}{
		"success": true,
		"payload": code,
	})
}

// routeDeleteFunk - remove a funk
func routeDeleteFunk(c echo.Context) error {
	funker.DeleteFunk(c.Param("funkName"))
	return c.JSON(200, map[string]interface{}{
		"success": true,
	})
}

// routeCallFunk - call a funk
func routeCallFunk(c echo.Context) error {
	res, err := funker.CallFunk(c, c.Param("funkName"))
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return res.echoify(c)
}

// routeEvalFunk - the playground
func routeEvalFunk(c echo.Context) error {
	allBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil || len(allBytes) < 1 {
		return c.JSON(500, map[string]interface{}{
			"success": false,
			"error":   "empty request body",
		})
	}

	code := fmt.Sprintf("(%s)", string(allBytes))

	_, err = goja.Compile("playground", code, true)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	name := "playground_" + xid.New().String()

	funker.AddFunk(name, code, true)

	res, err := funker.CallFunk(c, name)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return res.echoify(c)
}

// routeHome - the home
func routeHome(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Hi, I'm funker, Let's funkify * from your.mind ;)",
	})
}
