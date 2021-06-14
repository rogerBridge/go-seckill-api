package router

import (
	"go-seckill/internal/auth"
	"go-seckill/internal/controllers"
	"go-seckill/internal/controllers2"

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

	// purchase_limits table
	register(fasthttp.MethodPost, "/admin/createPurchaseLimit", controllers2.CreatePurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/queryPurchaseLimit", controllers2.QueryPurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/updatePurchaseLimit", controllers2.UpdatePurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/deletePurchaseLimit", controllers2.DeletePurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/loadGoodPurchaseLimit", controllers2.LoadGoodPurchaseLimit, auth.MiddleAuth)
	// syncGoodsFromMysql2Redis 在go-seckill初始化的时候就已经做到了, 不需要再做
	// register(fasthttp.MethodPost, "/admin/syncGoodsFromMysql2Redis", controllers.SyncGoodsFromMysql2Redis, auth.MiddleAuth)
	// 这个之后还是用rabbitmq-receiver来做吧, 每次redis库存扣减成功之后都发送消息, 让mysql也扣减
	register(fasthttp.MethodPost, "/admin/syncGoodsFromRedis2Mysql", controllers.SyncGoodsFromRedis2Mysql, auth.MiddleAuth)
	// goods table
	register(fasthttp.MethodGet, "/admin/goodList", controllers2.GoodsList, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/goodCreate", controllers2.CreateGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/goodUpdate", controllers2.UpdateGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/admin/goodDelete", controllers2.DeleteGood, auth.MiddleAuth)

	// use those API don't need auth :)
	register(fasthttp.MethodPost, "/user/register", controllers2.UserRegister, nil)
	register(fasthttp.MethodPost, "/user/login", controllers2.UserLogin, nil)

	// use those API need auth :)
	register(fasthttp.MethodPost, "/user/buy", controllers.Buy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/user/cancelBuy", controllers.CancelBuy, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/user/logout", controllers2.UserLogout, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/user/updatePassword", controllers2.UserUpdatePassword, auth.MiddleAuth)
	register(fasthttp.MethodPost, "/user/updateInfo", controllers2.UserUpdateInfo, auth.MiddleAuth)
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
