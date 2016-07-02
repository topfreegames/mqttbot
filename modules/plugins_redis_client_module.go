package modules

import (
	"fmt"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/layeh/gopher-luar"
	"github.com/topfreegames/mqttbot/logger"

	"github.com/spf13/viper"
	"github.com/yuin/gopher-lua"
)

var redisPool *redis.Pool
var once sync.Once

// RedisModuleLoader loads the Redis module
func RedisModuleLoader(L *lua.LState) int {
	loadDefaultConfigurations()
	InitRedisPool()
	mod := L.SetFuncs(L.NewTable(), redisClientModuleExports)
	L.Push(mod)
	return 1
}

var redisClientModuleExports = map[string]lua.LGFunction{
	"execute": ExecuteCommand,
}

// InitRedisPool starts the redis pool
func InitRedisPool() {
	once.Do(func() {
		loadDefaultConfigurations()
		logger.Logger.Info(fmt.Sprintf("Redis address: %s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")))
		redisPool = redis.NewPool(func() (redis.Conn, error) {
			if viper.GetString("redis.password") != "" {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("redis.host"),
					viper.GetInt("redis.port")), redis.DialPassword(viper.GetString("redis.password")))
				if err != nil {
					logger.Logger.Fatal(err)
				}
				return c, err
			}
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("redis.host"),
				viper.GetInt("redis.port")))
			if err != nil {
				logger.Logger.Fatal(err)
			}
			return c, err
		}, viper.GetInt("redis.maxPoolSize"))
	})
}

func loadDefaultConfigurations() {
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.maxPoolSize", 10)
}

// ExecuteCommand executes the command given by the Lua script
func ExecuteCommand(L *lua.LState) int {
	command := L.Get(1)
	argNum := L.Get(2)
	var args []interface{}
	for i := 3; i <= 2+int(lua.LVAsNumber(argNum)); i++ {
		args = append(args, L.Get(i).String())
	}
	logger.Logger.Debug(fmt.Sprintf("redismod command %s and args %s", command, args))
	L.Pop(2 + int(lua.LVAsNumber(argNum)))
	status, err := redisPool.Get().Do(command.String(), args...)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s", err)))
		L.Push(L.ToNumber(1))
		return 2
	}
	L.Push(lua.LNil)
	L.Push(luar.New(L, status))
	return 2
}
