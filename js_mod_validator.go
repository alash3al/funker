package main

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
)

func jsValidator() map[string]interface{} {
	return map[string]interface{}{
		"validate": func(data map[string]interface{}, validators map[string][]string) map[string]interface{} {
			invalid, result := 0, map[string][]string{}
			for k, rules := range validators {
				result[k] = []string{}
				value, exists := data[k]
				valuestr := strings.TrimSpace(fmt.Sprintf("%v", value))
				for _, r := range rules {
					if r == "required" && !exists || valuestr == "" {
						invalid++
						result[k] = append(result[k], "required")
					} else if ruler, ok := govalidator.TagMap[r]; ok && !ruler(valuestr) {
						invalid++
						result[k] = append(result[k], r)
					} else {
						parts := strings.SplitN(r, ":", 2)
						if len(parts) < 2 {
							parts = append(parts, "")
						}
						r, args := parts[0], parts[1]
						args = strings.TrimSpace(args)
						if ruler, ok := govalidator.ParamTagMap[r]; ok {
							if !ruler(valuestr, strings.Split(args, ",")...) {
								invalid++
								result[k] = append(result[k], r)
							}
						}
					}
				}
			}
			return map[string]interface{}{
				"errors": invalid,
				"fields": result,
			}
		},
	}
}
