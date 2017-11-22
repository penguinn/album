package router

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/penguinn/penguin/component/controller"
	"github.com/penguinn/penguin/utils"
	"net/http"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"

	PathTag = "path"
	PermTag = "perm"
)

var (
	GlobalRouterConfig RouterConfig

	GlobalRouter *gin.Engine

	registerControllerShouldBeValueType = errors.New("Register Controller Should be Value Type")

	actionType = reflect.TypeOf(func(*gin.Context) {})

	routePerms = make(map[interface{}]string)

	methodSupport = map[string]func(*gin.RouterGroup, string, func(*gin.Context)){
		GET:     func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.GET(p, h) },
		POST:    func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.POST(p, h) },
		PUT:     func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.PUT(p, h) },
		DELETE:  func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.DELETE(p, h) },
		HEAD:    func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.HEAD(p, h) },
		PATCH:   func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.PATCH(p, h) },
		OPTIONS: func(r *gin.RouterGroup, p string, h func(*gin.Context)) { r.OPTIONS(p, h) },
	}

	actionParser = func(f reflect.StructField, v reflect.Value, r *gin.RouterGroup) bool {

		if f.Type != actionType {
			return false
		}

		p := f.Tag.Get(PathTag)
		if p == "-" {
			return false
		}

		if p == "" {
			p = "/"
		}
		name := f.Name
		h := v.Interface().(func(*gin.Context))
		perm := strings.TrimSpace(f.Tag.Get(PermTag))
		if perm != "" {
			routePerms[utils.NameOfFunction(h)] = perm
		}
		return guessMethod(r, 3, name, h, p) || guessMethod(r, 4, name, h, p) || guessMethod(r, 5, name, h, p)

	}

	guessMethod = func(r *gin.RouterGroup, n int, name string, h func(*gin.Context), p string) bool {
		nameB := []byte(name)
		if len(nameB) < n {
			return false
		}
		m := nameB[:n]

		if sm, ok := methodSupport[strings.ToUpper(string(m))]; ok {
			sm(r, p, h)
			return true
		}
		return false
	}
)

type Request struct {
}

type Response struct {
}

type RouterConfig struct {
	Addr  string
	Mode  string
	Pprof bool
	Host  string
}

type RouterComponent struct {
}

func (c RouterComponent) Init(options ...interface{}) error {

	if len(options) != 0 {
		conf, ok := options[0].(*RouterConfig)
		if ok {
			GlobalRouterConfig = *conf
		} else {
			panic(errors.New("Cant find GlobalRouter config!"))
		}
	} else {
		panic(errors.New("Cant find GlobalRouter config!"))
	}

	gin.SetMode(GlobalRouterConfig.Mode)

	GlobalRouter = gin.Default()

	return nil
}

func Use(middlewares ...gin.HandlerFunc) {
	GlobalRouter.Use(middlewares...)
}

func GroupUse(g string, middlewares ...gin.HandlerFunc) {
	GlobalRouter.Group(g).Use(middlewares...)
}

func RegisterController(c controller.Controller, middlewares ...gin.HandlerFunc) error {

	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		return registerControllerShouldBeValueType
	}

	v := reflect.ValueOf(c)
	fN := t.NumField()

	GlobalRouter.RouterGroup.Use(middlewares...)

	for i := 0; i < fN; i++ {
		actionParser(t.Field(i), v.Field(i), &GlobalRouter.RouterGroup)
	}
	return nil
}

func RegisterControllerGroup(c controller.Controller, g string, middlewares ...gin.HandlerFunc) error {
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		return registerControllerShouldBeValueType
	}

	v := reflect.ValueOf(c)
	fN := t.NumField()

	r := &GlobalRouter.RouterGroup
	r = r.Group(g)
	r.Use(middlewares...)

	for i := 0; i < fN; i++ {
		actionParser(t.Field(i), v.Field(i), r)
	}
	return nil
}

func RoutePerms() map[interface{}]string {
	return routePerms
}

func RoutePerm(f gin.HandlerFunc) string {
	perm, ok := routePerms[utils.NameOfFunction(f)]
	if !ok {
		return ""
	}
	return perm
}

func BindRequest(c *gin.Context, obj interface{}) error {
	return c.ShouldBindWith(obj, binding.Default(c.Request.Method, c.ContentType()))
}

//映射几个gin的方法出来
func StaticFile(relativePath, filepath string) gin.IRoutes {
	return GlobalRouter.StaticFile(relativePath, filepath)
}

func Static(relativePath, root string) gin.IRoutes {
	return GlobalRouter.Static(relativePath, root)
}

func StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return GlobalRouter.StaticFS(relativePath, fs)
}
