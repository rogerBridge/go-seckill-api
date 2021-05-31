package router

import (
	"redisplay/auth"
	"redisplay/controllers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

/*
	自定义route结构体
*/
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
	register(fasthttp.MethodPost, "/syncGoodsFromMysql2Redis", controllers.SyncGoodsFromMysql2Redis, nil)
	register(fasthttp.MethodPost, "/syncGoodsFromRedis2Mysql", controllers.SyncGoodsFromRedis2Mysql, nil)

	// use those interface need auth :)
	register(fasthttp.MethodGet, "/goodsList", controllers.GoodsList, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/addGood", controllers.AddGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/modifyGood", controllers.ModifyGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/deleteGood", controllers.DeleteGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/buy", controllers.Buy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/cancelBuy", controllers.CancelBuy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/logout", controllers.Logout, auth.MiddleAuth)

	register(fasthttp.MethodPost, "/register", controllers.Register, nil)
	register(fasthttp.MethodPost, "/login", controllers.Login, nil)
}

// ThisRouter 通过遍历[]Route, 将需要中间件处理的和不需要中间件处理的分开处置 :)
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
