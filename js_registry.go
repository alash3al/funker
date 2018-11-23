package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func jsModules() map[string]interface{} {
	return map[string]interface{}{
		"fetch":        jsFetch,
		"crypto":       jsCrypto(),
		"localStorage": jsKVStore(),
		"uniqid": func(l ...int) string {
			if len(l) < 1 {
				l = []int{15}
			}
			b := make([]byte, l[0])
			rand.Read(b)
			return hex.EncodeToString(b)
		},
		"base64": map[string]interface{}{
			"encode": func(s string) string {
				return base64.StdEncoding.EncodeToString([]byte(s))
			},
			"decode": func(s string) string {
				b, _ := base64.StdEncoding.DecodeString(s)
				return string(b)
			},
		},
	}
}

func jsScripts() string {
	return strings.Join([]string{
		jsUnderscore(),
	}, ";")
}
