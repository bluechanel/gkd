package gkd

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup	// 继承RouterGroup的属性
	router *router	// 路由器
	groups []*RouterGroup	// 路由组数组
}

// 路由组
type RouterGroup struct {
	prefix string  // 前缀
	middleWares []HandlerFunc	// 中间件
	parent *RouterGroup	// 支持分组嵌套
	engine *Engine	// 所有的分组共享一个engine实例
}
// New is the constructor of gkd.Engine
// 相当于调用的时候初始化
func New() *Engine {
	engine := &Engine{router:NewRouter()}	// 创建engine对象
	engine.RouterGroup = &RouterGroup{	//
		engine: engine,
	}	// 初始化RouterGroup
	engine.groups = []*RouterGroup{engine.RouterGroup}	// 初始化engine里面的路由分组
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine	// 初始化RouteGroup中的engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,	// 组装前缀
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)	// 将新的带有prefix前缀group放入engine的group数组中
	return newGroup
}
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)// 组装pattern
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 使用use方法将中间件添加到group中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middleWares = append(group.middleWares, middlewares...)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleWares...)
		}
	}
	context := newContext(w,req)
	context.handlers = middlewares
	engine.router.handle(context)
}
