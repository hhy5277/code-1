package log

import (
	"github.com/TeaWeb/code/teaproxy"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"time"
)

// 是否正在获取请求内容
var requestBodyFetching = false
var requestBodyTime = time.Now()

// 初始化
func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 路由设置
		server.
			EndAll().
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantLog,
			}).
			Helper(new(Helper)).
			Prefix("/log").
			Get("", new(IndexAction)).
			Get("/get", new(GetAction)).
			Get("/responseHeader/:logId", new(ResponseHeaderAction)).
			Get("/requestHeader/:logId", new(RequestHeaderAction)).
			Get("/cookies/:logId", new(CookiesAction)).
			GetPost("/runtime", new(RuntimeAction)).
			EndAll()

		// 请求Hook
		hook := &teaproxy.RequestHook{
			BeforeRequest: ProcessBeforeRequest,
			AfterRequest:  nil,
		}
		teaproxy.AddRequestHook(hook)
	})
}

// 请求Hook
func ProcessBeforeRequest(req *teaproxy.Request, writer *teaproxy.ResponseWriter) bool {
	if requestBodyFetching && time.Since(requestBodyTime).Seconds() < 5 {
		req.SetIsWatching(true)
	} else {
		requestBodyFetching = false
		req.SetIsWatching(false)
	}
	return true
}
