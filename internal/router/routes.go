package router

import (
	"go-seckill/internal/auth"
	"go-seckill/internal/controllers2"
	"go-seckill/internal/utils"

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

// 提前分配好内存, 一般一个应用的话, 1024个API是足够的
var routes = make([]Route, 0, 1024)

// 路由中间件注册
func register(method, pattern string, handler fasthttp.RequestHandler, middle func(handler fasthttp.RequestHandler) fasthttp.RequestHandler) {
	routes = append(routes, Route{method, pattern, handler, middle})
}

const ApiVersion = utils.API_VERSION

func init() {
	// use this token needing permission
	// 用户和管理员的权限是应该有区分的, 这里并没有做什么区分, 后面要修改的

	// purchase_limits table
	register(fasthttp.MethodPost, ApiVersion+"/admin/createPurchaseLimit", controllers2.CreatePurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/queryPurchaseLimit", controllers2.QueryPurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/queryPurchaseLimits", controllers2.QueryPurchaseLimits, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/updatePurchaseLimit", controllers2.UpdatePurchaseLimit, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/deletePurchaseLimit", controllers2.DeletePurchaseLimit, auth.MiddleAuth)
	//register(fasthttp.MethodPost, "/admin/loadGoodPurchaseLimit", controllers2.LoadGoodPurchaseLimit, auth.MiddleAuth)
	// syncGoodsFromMysql2Redis 在go-seckill初始化的时候就已经做到了, 不需要再做
	// register(fasthttp.MethodPost, "/admin/syncGoodsFromMysql2Redis", controllers-bak.SyncGoodsFromMysql2Redis, auth.MiddleAuth)
	// 这个之后还是用rabbitmq-receiver来做吧, 每次redis库存扣减成功之后都发送消息, 让mysql也扣减
	// 还是写成定时同步吧, 每隔60s同步redis中goods数据到mysql.goods中
	//register(fasthttp.MethodPost, "/admin/syncGoodsFromRedis2Mysql", controllers_bak.SyncGoodsFromRedis2Mysql, auth.MiddleAuth)

	// goods table

	register(fasthttp.MethodGet, ApiVersion+"/admin/goodList", controllers2.GoodList, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/goodCreate", controllers2.CreateGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/goodUpdate", controllers2.UpdateGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/goodDelete", controllers2.DeleteGood, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/admin/register", controllers2.AdminRegister, nil)

	// users table
	register(fasthttp.MethodPost, ApiVersion+"/user/register", controllers2.UserRegister, nil)
	register(fasthttp.MethodPost, ApiVersion+"/user/login", controllers2.UserLogin, nil)
	register(fasthttp.MethodPost, ApiVersion+"/user/logout", controllers2.UserLogout, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/user/updatePassword", controllers2.UserUpdatePassword, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/user/updateInfo", controllers2.UserUpdateInfo, auth.MiddleAuth)

	// orders table
	register(fasthttp.MethodPost, ApiVersion+"/user/order/buy", controllers2.Buy, auth.MiddleAuth)
	register(fasthttp.MethodPost, ApiVersion+"/user/order/cancelBuy", controllers2.CancelBuy, auth.MiddleAuth)
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
