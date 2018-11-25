package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"log"
	"runtime"

	"github.com/go-redis/redis"
)

var (
	flagListenAddr        = flag.String("listen", ":6020", "the listen address")
	flagRedisAddr         = flag.String("redis", "redis://localhost:6379/0", "redis server address")
	flagWorkers           = flag.Int("workers", runtime.NumCPU(), "the number of workers/cores")
	flagInitialVMPoolSize = flag.Int("pool", 100, "the initial value of the default VM pool size per each worker")
	flagAuthKey           = flag.String("auther", "", "the authentication key for adding or running any funk in the playground")
)

var (
	redisClient   *redis.Client
	funker        *FunkerManager
	redisFunksKey = "funker:funks"
	modules       = map[string]interface{}{
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
)

func init() {
	flag.Parse()

	redisOptions, err := redis.ParseURL(*flagRedisAddr)
	if err != nil {
		log.Fatal("[redis]", err.Error())
		return
	}

	redisClient = redis.NewClient(redisOptions)
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal("[redis]", err.Error())
	}

	runtime.GOMAXPROCS(*flagWorkers)

	funker = NewFunker()
}
