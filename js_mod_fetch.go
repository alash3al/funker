package main

import (
	"errors"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

func jsFetch(url string, opts map[string]interface{}) (resp interface{}) {
	var method string
	var headers map[string]string
	var body interface{}
	var redirects int
	var timeout int
	var proxy string

	if opts["method"] == nil {
		method = "GET"
	} else {
		method = strings.ToUpper(opts["method"].(string))
	}

	if opts["headers"] != nil {
		hdrs := opts["headers"].(map[string]interface{})
		headers = map[string]string{}
		for k, v := range hdrs {
			headers[k] = v.(string)
		}
	}

	if opts["body"] != nil {
		body = opts["body"]
	}

	if opts["timeout"] == nil {
		timeout = 1
	} else {
		timeout = opts["timeout"].(int)
	}

	if opts["redirects"] == nil {
		redirects = 1
	} else {
		redirects = opts["redirects"].(int)
	}

	if opts["proxy"] != nil {
		proxy = opts["proxy"].(string)
	}

	c := resty.New()
	c.AllowGetMethodPayload = true
	c.SetDoNotParseResponse(false)
	c.SetRedirectPolicy(resty.FlexibleRedirectPolicy(redirects))
	c.SetTimeout(time.Duration(timeout) * time.Second)

	if proxy != "" {
		c.SetProxy(proxy)
	}

	r := c.R()

	if body != nil {
		r.SetBody(body)
	}

	if headers != nil {
		r.SetHeaders(headers)
	}

	res, err := r.Execute(method, url)
	if err != nil {
		return errors.New(err.Error())
	}

	respHeaders := map[string]interface{}{}
	for k, v := range res.Header() {
		k = strings.ToLower(k)
		if len(v) > 1 {
			respHeaders[k] = v
		} else {
			respHeaders[k] = v[0]
		}
	}

	resp = map[string]interface{}{
		"code":    res.StatusCode(),
		"headers": respHeaders,
		"body":    string(res.Body()),
	}

	return resp
}
