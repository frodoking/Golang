/**
 * Copyright 2015 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2015-04-27 00:53
 * description :
 * history :
 */
package main

import (
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
	"github.com/atnet/gof/storage"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"net/http"
	"strings"
)

// 实现gof.App接口，以支持配置(config)、DB、ORM、日志(logger)、
// 模版(template)、存储(storage)的访问。同时可以扩展,以支持自定义
// 的需求。App接口在应用上下文中可以获取到。
type HttpApp struct {
	config      *gof.Config
	dbConnector db.Connector
	template    *gof.Template
	logger      log.ILogger
}

// 配置，支持从文件中加载。用于文件存储某些简单数据。
func (this *HttpApp) Config() *gof.Config {
	if this.config == nil {
		this.config = gof.NewConfig()
		this.config.Set("SYS_NAME", "GOF DEMO")
		this.config.Set("MYSQL_HOST", "127.0.0.1")
		this.config.Set("MYSQL_PORT", 3306)
		this.config.Set("MYSQL_MAXCONN", 1000)
		this.config.Set("MYSQL_USR", "root")
		this.config.Set("MYSQL_PWD", "")
		this.config.Set("MYSQL_DBNAME", "")
	}
	return this.config
}

// 数据库连接器、Connector.ORM()可以访问ORM
func (this *HttpApp) Db() db.Connector {
	if this.dbConnector == nil {
		source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf-8&loc=Local",
			this.Config().GetString("MYSQL_USR"),
			this.Config().GetString("MYSQL_PWD"),
			this.Config().GetString("MYSQL_HOST"),
			this.Config().GetInt("MYSQL_PORT"),
			this.Config().GetString("MYSQL_DBNAME"),
		)
		this.dbConnector = db.NewSimpleConnector("mysql", source, this.Log(),
			this.Config().GetInt("MYSQL_MAXCONN"))
	}
	return this.dbConnector
}

// 模板
func (this *HttpApp) Template() *gof.Template {
	if this.template == nil {
		this.template = &gof.Template{
			Init: func(v *gof.TemplateDataMap) {
				(*v)["SYS_NAME"] = this.Config().GetString("SYS_NAME")
			},
		}
	}
	return this.template
}

// 是否调试
func (this *HttpApp) Debug() bool {
	return false
}

// 存储,通常用来存储全局的变量。或缓存一些共享的数据。
func (this *HttpApp) Storage() gof.Storage {
	// 使用一个Redis存储数据
	return storage.NewRedisStorage(nil)
}

// 日志
func (this *HttpApp) Log() log.ILogger {
	if this.logger == nil {
		var flag int = 0
		if this.Debug() {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		// 可创建文件的io.Writer,将日志存储到文件。
		this.logger = log.NewLogger(nil, " O2O", flag)
	}
	return this.logger
}

//获取Http请求代理处理程序
func getInterceptor(a gof.App, routes web.Route) *web.Interceptor {
	var in = web.NewInterceptor(a, getHttpExecFunc(routes))
	// 处理发生的异常
	in.Except = web.HandleDefaultHttpExcept
	in.Before = nil
	in.After = nil
	return in
}

// 获取执行方法
func getHttpExecFunc(routes web.Route) web.RequestHandler {
	return func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		switch host, f := r.Host, strings.HasPrefix; {
		//静态文件，处理http://static.domain.com的请求。
		case f(host, "static."):
			http.ServeFile(w, r, "./static"+r.URL.Path)
		default:
			routes.Handle(ctx)
		}
	}
}

// 控制器（任意类型）
type testController struct {
}

// 控制器可以实现mvc.Filter接口
var _ mvc.Filter = new(testController)

// 控制器生成器
func testControllerGenerator() mvc.Controller {
	return &testController{}
}

// 请求时执行函数，返回true继续执行，返回false即中止。
// 可以做验证如：判断用户登陆这类的逻辑。
func (this *testController) Requesting(ctx *web.Context) bool {
	return true
}

// 请求结束时执行。
// 可以做统一设置Header,gzip压缩，页面执行时间此类逻辑
func (this *testController) RequestEnd(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("\r\nRequest End."))
}

// Index动作，GET请求。来源URL:/test/index
// Index为默认的动作，即也可以通过/test访问呢。
// 如果为POST请求，约定在动作名称后添加"_post",即"Index_post"
func (this *testController) Index(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("\r\nRequesting....."))
}

func main() {
	app := &HttpApp{}

	// MVC路由表
	routes := mvc.NewRoute(nil)

	// 注册控制器,test与/test/index中的test对应
	routes.NormalRegister("test", testControllerGenerator)

	// 注册单例控制器
	routes.Register("test_all", &testController{})

	// 除了控制器，可以添加自定义的路由规则（正则表达式)
	routes.Add("^/[0-9]$", func(ctx *web.Context) {
		// 直接输出内容
		ctx.ResponseWriter.Write([]byte("数字路径"))
		return
	})

	// 默认页路由
	routes.Add("/", func(ctx *web.Context) {
		//		sysName := ctx.App.Config().GetString("SYS_NAME")
		//		ctx.ResponseWriter.Write([]byte("Hello," + sysName + "!"))
		//		return

		// 使用模板
		data := gof.TemplateDataMap{
			"变量名": "变量值",
		}
		ctx.App.Template().Execute(ctx.ResponseWriter, data, "template.html")
		return

		// 使用会话
		ctx.Session().Set("user", "jarrysix")
		ctx.Session().Save()

		// 使用DB和ORM
		db := ctx.App.Db()
		orm := db.GetOrm()
		_ = orm.Version()
	})

	// 使用一个拦截器，来拦截请求。
	// 拦截器里可以决定，访问对应的路由表。
	// 一个系统通常有多个子系统，每个子系统对应一个路由表。
	var in = getInterceptor(app, routes)
	go http.ListenAndServe(":8080", in)

	log.Println("[ OK] - web is listening on port :8080.")
	var ch = make(chan int, 1)
	<-ch
}
