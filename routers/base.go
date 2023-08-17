package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/utils"
	"net/http"
	"small-website-monitor/global"
	"small-website-monitor/handler"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func Routers() *gin.Engine {
	Router.StaticFS("/static", http.Dir("./static"))
	Router.Delims("{[", "]}")
	Router.LoadHTMLGlob("views/**/*")
	Login()
	Page()
	API()
	Router.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.HTML(http.StatusOK, "notfond.html", gin.H{})
	})
	return Router
}

func Login() {
	login := Router.Group("")
	login.Use(ginHelper.CSRFMiddleware())
	login.GET("/", handler.LoginPage)
	login.POST("/login", handler.Login)
}

func Page() {
	pg := Router.Group("")
	pg.Use(AuthPG())
	pg.GET("/home", handler.HomePage)
	// TODO 监测器
	// TODO 工具
	// TODO 使用说明
}

func API() {
	api := Router.Group("/api").Use(AuthAPI())
	api.GET("/out", handler.Out)
	api.GET("/test", ginHelper.Handle(handler.CaseT))
	api.POST("/website/add", ginHelper.Handle(handler.WebsiteAdd))
	api.GET("/website/list", ginHelper.Handle(handler.WebsiteList))            //
	api.GET("/mail/init", ginHelper.Handle(handler.MailInit))                  // 是否设置邮件
	api.POST("/mail/conf", ginHelper.Handle(handler.MailConf))                 // 设置邮件配置
	api.GET("/mail/info", ginHelper.Handle(handler.MailInfo))                  // 获取邮件配置信息
	api.POST("/mail/sendTest", ginHelper.Handle(handler.MailSendTest))         // 测试发生邮件
	api.POST("/point/add/:hostId", ginHelper.Handle(handler.WebsitePointAdd))  // 添加监测点
	api.GET("/point/list/:hostId", ginHelper.Handle(handler.WebsitePointList)) // 获取监测点
	api.POST("/point/del/:hostId", ginHelper.Handle(handler.WebsitePointDel))  // 删除监测点
	// TODO 获取当前站点采集的URI列表
	// TODO 更新URI列表
	// TODO 查看日志
	// TODO 设置
	// TODO 删除
	// TODO 图表

	// 测试
	api.GET("/test/case1", ginHelper.Handle(handler.Case1))

}

// AuthPG 权限验证中间件
func AuthPG() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(global.UserToken)
		j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
		if err := j.ParseToken(token); err == nil {
			c.Next()
			return
		}
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}
}

// AuthAPI 权限验证中间件
func AuthAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(global.UserToken)
		j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
		if err := j.ParseToken(token); err == nil {
			c.Next()
			return
		}
		ginHelper.AuthErrorOut(c)
		c.Abort()
		return

	}
}
