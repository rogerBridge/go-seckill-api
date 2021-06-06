package router

import (
	"go-seckill/internal/auth"
	"go-seckill/internal/controllers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//自定义route结构体
type Route struct {
	Method     string
	Pattern    string
	Handler    fasthttp.RequestHandler
	Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler
}

// 提前分配好内存, 一般一个应用的话, 256个API是足够的
var routes = make([]Route, 0, 256)

// 路由中间件注册
func register(method, pattern string, handler fasthttp.RequestHandler, middle func(handler fasthttp.RequestHandler) fasthttp.RequestHandler) {
	routes = append(routes, Route{method, pattern, handler, middle})
}

func init() {
	// use this token needing permission
	// 用户和管理员的权限是应该有区分的, 这里并没有做什么区分, 后面要修改的
	register(fasthttp.MethodPost, "/syncGoodsLimit", controllers.SyncGoodsLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/syncGoodsFromMysql2Redis", controllers.SyncGoodsFromMysql2Redis, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/syncGoodsFromRedis2Mysql", controllers.SyncGoodsFromRedis2Mysql, auth.MiddleAuth)
	register(fasthttp.MethodGet, "/goodsList", controllers.GoodsList, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/addGood", controllers.AddGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/modifyGood", controllers.ModifyGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/deleteGood", controllers.DeleteGood, auth.MiddleAuth)

	// use those API don't need auth :)
	register(fasthttp.MethodPost, "/register", controllers.Register, nil)
	register(fasthttp.MethodPost, "/login", controllers.Login, nil)

	// use those API need auth :)
	register(fasthttp.MethodPost, "/buy", controllers.Buy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/cancelBuy", controllers.CancelBuy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/logout", controllers.Logout, auth.MiddleAuth)
}

// ThisRouter 通过遍历[]Route, 将需要中间件处理的和不需要中间件处理的分开处置 :)
func ThisRouter() *router.Router {
	r := router.New()
	for _, route := range routes {
		if route.Middleware != nil {
			r.Handle(route.Method, route.Pattern, route.Middleware(route.Handler))
		} else {
			r.Handle(route.Method, route.Pattern, route.Handler)
		}
	}
	return r
}
