package main

import (
	"strings"

	"github.com/labstack/echo"
)

func jsRequestEnv(c echo.Context) map[string]interface{} {
	var body interface{}

	c.Bind(&body)

	req := c.Request()

	headers := map[string]interface{}{}
	for k, v := range req.Header {
		k = strings.ToLower(k)
		if len(v) > 1 {
			headers[k] = v
		} else {
			headers[k] = v[0]
		}
	}

	params := map[string]interface{}{}
	for k, v := range c.QueryParams() {
		k = strings.ToLower(k)
		if len(v) > 1 {
			params[k] = v
		} else {
			params[k] = v[0]
		}
	}

	if body == nil {
		b := map[string]interface{}{}
		for k, v := range req.Form {
			k = strings.ToLower(k)
			if len(v) > 1 {
				b[k] = v
			} else {
				b[k] = v[0]
			}
		}
		body = b
	}

	return map[string]interface{}{
		"uri":         req.URL.RequestURI(),
		"proto":       req.Proto,
		"method":      req.Method,
		"path":        req.URL.Path,
		"host":        req.Host,
		"https":       req.TLS != nil,
		"query":       params,
		"body":        body,
		"remote_addr": req.RemoteAddr,
		"real_ip":     c.RealIP(),
		"headers":     headers,
	}
}
