package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"

	"github.com/labstack/echo"

	"github.com/dop251/goja"
	"github.com/go-redis/redis"
)

// FunkerManager - the funker runtime
type FunkerManager struct {
	redis         *redis.Client
	globalRuntime *goja.Runtime
	funksCache    map[string]*goja.Program
	funksLocker   sync.RWMutex
	vmpool        sync.Pool
}

// NewFunker - creates a new funker runtime
func NewFunker() *FunkerManager {
	f := new(FunkerManager)

	f.redis = redisClient
	f.globalRuntime = goja.New()
	f.funksCache = map[string]*goja.Program{}
	f.funksLocker = sync.RWMutex{}
	f.vmpool = sync.Pool{
		New: func() interface{} {
			return f.newJsRuntime()
		},
	}

	f.RefreshCache()

	for a := 0; a < *flagWorkers; a++ {
		go (func() {
			for i := 0; i < *flagInitialVMPoolSize; i++ {
				f.vmpool.Put(f.newJsRuntime())
			}
		})()
	}

	return f
}

// newJSVM - init a new js runtime
func (f *FunkerManager) newJsRuntime() *goja.Runtime {
	vm := goja.New()
	for name, mod := range jsModules() {
		vm.Set(name, mod)
	}
	vm.RunString(fmt.Sprintf("vm = {id: '%s'}", xid.New().String()))
	vm.RunString(jsScripts())
	return vm
}

// RefreshCache - refreshes the funker cache to speedup executions
func (f *FunkerManager) RefreshCache() {
	f.funksLocker.Lock()
	defer f.funksLocker.Unlock()

	for name, fn := range f.redis.HGetAll(redisFunksKey).Val() {
		prog, err := goja.Compile(name, fn, true)
		if err != nil {
			continue
		}
		f.funksCache[name] = prog
	}
}

// AddFunk - add/override the specified funk
func (f *FunkerManager) AddFunk(name, code string, cacheTTL int64) error {
	f.funksLocker.Lock()
	defer f.funksLocker.Unlock()

	name = strings.ToLower(name)

	prog, err := goja.Compile("", code, true)
	if err != nil {
		return err
	}

	if _, err := f.redis.HSet(redisFunksKey, name, code).Result(); err != nil {
		return err
	}

	f.redis.IncrBy(redisFunksTTLKey+name, cacheTTL).Result()

	f.funksCache[name] = prog

	return nil
}

// DeleteFunk - Delete a funk
func (f *FunkerManager) DeleteFunk(name string) {
	f.funksLocker.Lock()
	defer f.funksLocker.Unlock()

	name = strings.ToLower(name)
	delete(f.funksCache, name)
	f.redis.HDel(redisFunksKey, name).Result()
	f.redis.Del(redisFunksTTLKey + name)
}

// CallFunk - executes the specified funk with the specified scope/context
func (f *FunkerManager) CallFunk(ctx echo.Context, name string) (*jsExports, error) {
	name = strings.ToLower(name)
	cacheKey := redisCachePrefix + base64.StdEncoding.EncodeToString([]byte(ctx.Request().Header.Get("Authorization")+ctx.Request().RequestURI))
	cacheTTL, _ := f.redis.Get(redisFunksTTLKey + name).Int()
	if cacheTTL > 0 && f.redis.Exists(cacheKey).Val() > 0 {
		var d jsExports
		json.Unmarshal([]byte(f.redis.Get(cacheKey).Val()), &d)
		return &d, nil
	}

	f.funksLocker.RLock()
	defer f.funksLocker.RUnlock()

	name = strings.ToLower(name)
	code := f.funksCache[name]

	if nil == code {
		return nil, errors.New("funk not found")
	}

	res, err := f.EvalFunk(ctx, code)
	if err != nil {
		return nil, err
	}

	d, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	if cacheTTL > 0 {
		f.redis.Set(cacheKey, string(d), time.Duration(cacheTTL)*time.Second).Result()
	}

	return res, nil
}

// EvalFunk - executes the specified funk code
func (f *FunkerManager) EvalFunk(ctx echo.Context, prog *goja.Program) (*jsExports, error) {
	vm := (f.vmpool.Get()).(*goja.Runtime)
	defer f.vmpool.Put(vm)

	exports := &jsExports{
		DocType: "json",
		Status:  200,
		Headers: map[string]string{},
	}

	vm.Set("env", jsRequestEnv(ctx))
	vm.Set("exports", exports)

	_, err := vm.RunProgram(prog)
	if err != nil {
		return nil, err
	}

	return exports, nil
}
