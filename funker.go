package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/labstack/echo"

	"github.com/dop251/goja"
	"github.com/go-redis/redis"
)

// FunkerManager - the funker runtime
type FunkerManager struct {
	redis       *redis.Client
	funksCache  map[string]string
	funksLocker sync.RWMutex
	vmpool      sync.Pool
}

// NewFunker - creates a new funker runtime
func NewFunker() *FunkerManager {
	f := new(FunkerManager)

	f.redis = redisClient
	f.funksCache = map[string]string{}
	f.funksLocker = sync.RWMutex{}
	f.vmpool = sync.Pool{
		New: func() interface{} {
			return goja.New()
		},
	}

	f.RefreshCache()

	for a := 0; a < *flagWorkers; a++ {
		go (func() {
			for i := 0; i < *flagInitialVMPoolSize; i++ {
				f.vmpool.Put(f.vmpool.New())
			}
		})()
	}

	return f
}

// RefreshCache - refreshes the funker cache to speedup executions
func (f *FunkerManager) RefreshCache() {
	f.funksLocker.Lock()
	defer f.funksLocker.Unlock()
	for name, code := range f.redis.HGetAll(redisFunksKey).Val() {
		_, err := goja.Compile(name, code, true)
		if err != nil {
			log.Println("[FunksInitializer]", err.Error())
			continue
		}
		f.funksCache[name] = code
	}
}

// AddFunk - add/override the specified funk
func (f *FunkerManager) AddFunk(name, code string, temp bool) error {
	f.funksLocker.Lock()
	defer f.funksLocker.Unlock()

	name = strings.ToLower(name)

	_, err := goja.Compile("", code, true)
	if err != nil {
		return err
	}

	if !temp {
		if _, err := f.redis.HSet(redisFunksKey, name, code).Result(); err != nil {
			return err
		}
	}

	f.funksCache[name] = code

	return nil
}

// DeleteFunk - Delete a funk
func (f *FunkerManager) DeleteFunk(name string) {
	f.funksLocker.Lock()
	defer f.funksLocker.Unlock()

	name = strings.ToLower(name)
	if f.redis.HExists(redisFunksKey, name).Val() {
		delete(f.funksCache, name)
		f.redis.HDel(redisFunksKey, name).Result()
	}
}

// CallFunk - executes the specified funk code
func (f *FunkerManager) CallFunk(ctx echo.Context, name string) (*jsExports, error) {
	code := f.funksCache[name]
	exports := &jsExports{
		DocType: "json",
		Status:  200,
		Headers: map[string]string{},
	}

	jsThis := map[string]interface{}{}

	jsThis = map[string]interface{}{
		"request": jsRequestEnv(ctx),
		"response": map[string]interface{}{
			"status": func(code int) interface{} {
				exports.Status = code
				return jsThis["response"]
			},
			"type": func(typ string) interface{} {
				exports.DocType = typ
				return jsThis["response"]
			},
			"headers": func(h map[string]string) interface{} {
				for k, v := range h {
					exports.Headers[k] = v
				}
				return jsThis["response"]
			},
			"send": func(data interface{}) interface{} {
				exports.Body = data
				return jsThis["response"]
			},
		},
		"module": func(name string) interface{} {
			mod, ok := modules[name]
			if !ok {
				return false
			}
			return mod
		},
	}

	vm := f.vmpool.Get().(*goja.Runtime)
	defer f.vmpool.Put(vm)

	vm.Set("context", jsThis)
	_, err := vm.RunString(fmt.Sprintf("%s.apply(context)", code))

	if nil != err {
		exports.Status = 500
		exports.Headers["Content-Type"] = "application/json"
		exports.Body = map[string]interface{}{
			"error": err.Error(),
		}
	}

	return exports, nil
}
