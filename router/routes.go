package router

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go_redis/auth"
	"go_redis/controllers"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    fasthttp.RequestHandler
}

var routes = make([]Route, 0)

func register(method, pattern string, handler fasthttp.RequestHandler) {
	routes = append(routes, Route{method, pattern, handler})
}

func init() {
	register(fasthttp.MethodPost, "/syncGoodsLimit", controllers.SyncGoodsLimit)
	register(fasthttp.MethodPost, "/goodsList", controllers.GoodsList)
	register(fasthttp.MethodPost, "/syncGoodsFromMysql2Redis", controllers.SyncGoodsFromMysql2Redis)
	register(fasthttp.MethodPost, "/syncGoodsFromRedis2Mysql", controllers.SyncGoodsFromRedis2Mysql)
	register(fasthttp.MethodPost, "/buy", auth.MiddleAuth(controllers.Buy))
	register(fasthttp.MethodPost, "/cancelBuy", controllers.CancelBuy)
	register(fasthttp.MethodPost, "/login", controllers.Login)
	register(fasthttp.MethodPost, "/logout", controllers.Logout)
	register(fasthttp.MethodPost, "/register", controllers.Register)
}

//
func ThisRouter() *router.Router{
	r := router.New()
	for _, route := range routes {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
	return r
}