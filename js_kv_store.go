package main

func jsKVStore() map[string]interface{} {
	return map[string]interface{}{
		"has": func(ns, k string) bool {
			ns = "funker:kv:" + ns
			return redisClient.HExists(ns, k).Val()
		},
		"set": func(ns, k, v string) {
			ns = "funker:kv:" + ns
			redisClient.HSet(ns, k, v).Result()
		},
		"incr": func(ns, k string, by int64) {
			ns = "funker:kv:" + ns
			redisClient.HIncrBy(ns, k, by).Result()
		},
		"get": func(ns, k string) string {
			ns = "funker:kv:" + ns
			return redisClient.HGet(ns, k).Val()
		},
		"delete": func(ns, k string) {
			ns = "funker:kv:" + ns
			redisClient.HDel(ns, k).Result()
		},
		"getAll": func(ns string) map[string]string {
			ns = "funker:kv:" + ns
			return redisClient.HGetAll(ns).Val()
		},
		"deleteAll": func(ns string) {
			ns = "funker:kv:" + ns
			redisClient.Del(ns).Result()
		},
	}
}
