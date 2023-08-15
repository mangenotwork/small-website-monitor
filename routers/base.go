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
}

func API() {
	api := Router.Group("/api").Use(AuthAPI())
	api.GET("/out", handler.Out)
	api.GET("/test", ginHelper.Handle(handler.CaseT))
	api.POST("/website/add", ginHelper.Handle(handler.WebsiteAdd))
	api.GET("/website/list", ginHelper.Handle(handler.WebsiteList))   // TODO
	api.GET("/mail/init", ginHelper.Handle(handler.MailInit))         // 是否设置邮件
	api.POST("/mail/conf", ginHelper.Handle(handler.MailConf))        // 设置邮件配置
	api.GET("/mail/info", ginHelper.Handle(handler.MailInfo))         // 获取邮件配置信息
	api.GET("/mail/sendTest", ginHelper.Handle(handler.MailSendTest)) // 测试发生邮件
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
