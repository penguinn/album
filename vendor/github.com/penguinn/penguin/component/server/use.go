package server

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/penguinn/penguin/component/config"
	"github.com/penguinn/penguin/component/db"
	"github.com/penguinn/penguin/component/log"
	"github.com/penguinn/penguin/component/router"
	"reflect"
	"sync"
	"github.com/penguinn/penguin/component/mongo"
	"github.com/penguinn/penguin/component/redis"
)

const (
	useShouldNotBeAPointer   = "Use(SessionComponent) Should Be Value Type"
	componentsCannotUseTwice = "Components Can not Use Twice"
)

//这个key是放在UseComps这个map里的，用来判断重复等
const (
	CompConfigName = "ConfigComponent"
	CompLogName    = "LogComponent"
	CompRouterName = "RouterComponent"
	CompDBName     = "DBComponent"
	CompRedisName  = "RedisComponent"
	CompMongoName  = "MongoComponent"
)

//这个key是作为配置文件的关键字
const (
	CompDBConfigKey     = "mysql"
	CompRedisConfigKey  = "redis"
	CompServerConfigKey = "server"
	CompJWTConfigKey    = "jwt"
	CompLogConfigKey    = "log"
	CompMQConfigKey     = "mq"
	CompMongoConfigKey  = "mongo"
)

type Component interface {
	Init(options ...interface{}) error
}

type compSettingWrapper struct {
	K  string
	CM interface{}
}

var (
	useMutex sync.RWMutex
	useComps = map[string]Component{}

	compConfigMapping = map[string]compSettingWrapper{
		CompLogName:    {CompLogConfigKey, log.LogConfig{}},
		CompRouterName: {CompServerConfigKey, router.RouterConfig{}},
		CompDBName:     {CompDBConfigKey, db.DBConfig{}},
		CompMongoName:  {CompMongoConfigKey, mongo.MongoConfig{}},
		CompRedisName:  {CompRedisConfigKey, redis.ReidsConfig{}},
	}
)

//通过反射直接在配置文件里面捞文件，然后初始化，最后放到useComps里面
func Use(c Component, options ...interface{}) {

	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		panic(useShouldNotBeAPointer)
	}

	compKey := t.Name()

	if CompUsed(compKey) {
		panic(componentsCannotUseTwice)
	}

	//先判断没有参数传入，且不是配置设置，且设置了配置
	if len(options) == 0 && compKey != CompConfigName && CompUsed(CompConfigName) {
		compConfig, ok := compConfigMapping[compKey]

		if ok {
			cmt := reflect.TypeOf(compConfig.CM)
			cm := reflect.New(cmt).Interface()
			if err := mapstructure.WeakDecode(config.Get(compConfig.K), cm); err == nil {
				options = append(options, cm)
			} else {
				fmt.Println(compKey, " Parse Error : ", err)
			}
		}
	}
	useMutex.Lock()
	defer useMutex.Unlock()

	if err := c.Init(options...); err != nil {
		fmt.Println(compKey, err)
		panic(err)
	}
	useComps[compKey] = c
}

func CompUsed(comp string) bool {
	useMutex.RLock()
	defer useMutex.RUnlock()
	_, ok := useComps[comp]
	return ok
}
