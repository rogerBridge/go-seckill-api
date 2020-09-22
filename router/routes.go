package router

import (
	"go_redis/auth"
	"go_redis/controllers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Route struct {
	Method  string
	Pattern string
	Handler fasthttp.RequestHandler
	Middle  func(handler fasthttp.RequestHandler) fasthttp.RequestHandler
}

var routes = make([]Route, 0)

// 路由中间件注册
func register(method, pattern string, handler fasthttp.RequestHandler, middle func(handler fasthttp.RequestHandler) fasthttp.RequestHandler) {
	routes = append(routes, Route{method, pattern, handler, middle})
}

func init() {
	register(fasthttp.MethodPost, "/syncGoodsLimit", controllers.SyncGoodsLimit, nil)
	register(fasthttp.MethodGet, "/goodsList", controllers.GoodsList, nil)
	register(fasthttp.MethodPost, "/addGood", controllers.AddGood, nil)
	register(fasthttp.MethodPost, "/modifyGood", controllers.ModifyGood, nil)
	register(fasthttp.MethodPost, "/deleteGood", controllers.DeleteGood, nil)
	register(fasthttp.MethodPost, "/syncGoodsFromMysql2Redis", controllers.SyncGoodsFromMysql2Redis, nil)
	register(fasthttp.MethodPost, "/syncGoodsFromRedis2Mysql", controllers.SyncGoodsFromRedis2Mysql, nil)
	register(fasthttp.MethodPost, "/buy", controllers.Buy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/cancelBuy", controllers.CancelBuy, nil)
	register(fasthttp.MethodPost, "/login", controllers.Login, nil)
	register(fasthttp.MethodPost, "/logout", controllers.Logout, nil)
	register(fasthttp.MethodPost, "/register", controllers.Register, nil)
}

func ThisRouter() *router.Router {
	r := router.New()
	for _, route := range routes {
		if route.Middle != nil {
			r.Handle(route.Method, route.Pattern, route.Middle(route.Handler))
		} else {
			r.Handle(route.Method, route.Pattern, route.Handler)
		}
	}
	return r
}
