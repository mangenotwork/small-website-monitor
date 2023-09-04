package routers

import (
	"github.com/gin-contrib/gzip"
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
	Router.Use(gzip.Gzip(gzip.DefaultCompression))
	Router.StaticFS("/static", http.Dir("./static"))
	Router.Delims("{[", "]}")
	Router.LoadHTMLGlob("views/**/*")
	Login()
	Page()
	API()
	return Router
}

func Login() {
	login := Router.Group("")
	login.Use(ginHelper.CSRFMiddleware())
	login.GET("/", handler.LoginPage)
	login.POST("/login", handler.Login)
}

func Page() {
	// 404 && 405 && err page
	Router.NoRoute(handler.NotFond)
	Router.NoMethod(handler.NotFond)
	// page group
	pg := Router.Group("")
	pg.Use(AuthPG())
	pg.GET("/home", handler.HomePage)
	pg.GET("/monitor", handler.MonitorPage)
	pg.GET("/tool", handler.ToolPage)
	pg.GET("/mysql", handler.MysqlMonitorPage)
	pg.GET("/redis", handler.RedisMonitorPage)
	pg.GET("/sqlserver", handler.SqlServerMonitorPage)
	pg.GET("/operation", handler.OperationLogPage)
	pg.GET("/instructions", handler.InstructionsPage)
}

func API() {
	api := Router.Group("/api").Use(AuthAPI())
	api.GET("/out", handler.Out)
	api.GET("/test", ginHelper.Handle(handler.CaseT))
	api.POST("/website/add", ginHelper.Handle(handler.WebsiteAdd))
	api.GET("/website/list", ginHelper.Handle(handler.WebsiteList))                  //
	api.GET("/website/delete/:hostId", ginHelper.Handle(handler.WebsiteDelete))      //
	api.GET("/website/info/:hostId", ginHelper.Handle(handler.WebsiteInfo))          // 获取当前站点采集的URI列表
	api.POST("/website/edit", ginHelper.Handle(handler.WebsiteEdit))                 // 监测设置
	api.GET("/website/chart/:hostId", ginHelper.Handle(handler.WebsiteChart))        // 图表
	api.GET("/website/alert/:hostId", ginHelper.Handle(handler.WebsiteAlertList))    // 报警信息
	api.GET("/website/alert/del/:hostId", ginHelper.Handle(handler.WebsiteAlertDel)) // 报警信息
	api.GET("/log/list/:hostId", ginHelper.Handle(handler.WebsiteLogList))           // 日志列表
	api.GET("/log/upload/:hostId", ginHelper.Handle(handler.WebsiteLogUpload))       // 日志文件下载
	api.GET("/mail/init", ginHelper.Handle(handler.MailInit))                        // 是否设置邮件
	api.POST("/mail/conf", ginHelper.Handle(handler.MailConf))                       // 设置邮件配置
	api.GET("/mail/info", ginHelper.Handle(handler.MailInfo))                        // 获取邮件配置信息
	api.POST("/mail/sendTest", ginHelper.Handle(handler.MailSendTest))               // 测试发生邮件
	api.POST("/point/add/:hostId", ginHelper.Handle(handler.WebsitePointAdd))        // 添加监测点
	api.GET("/point/list/:hostId", ginHelper.Handle(handler.WebsitePointList))       // 获取监测点
	api.POST("/point/del/:hostId", ginHelper.Handle(handler.WebsitePointDel))        // 删除监测点
	api.GET("/alert/list", ginHelper.Handle(handler.AlertList))                      // 获取报警通知
	api.GET("/alert/clear", ginHelper.Handle(handler.AlertClear))                    // 清空报警通知
	api.GET("/monitor/err/list", ginHelper.Handle(handler.MonitorErrList))           // 获取监控平台错误日志
	api.GET("/monitor/err/clear", ginHelper.Handle(handler.MonitorErrClear))         // 清空监控平台错误日志
	api.GET("/monitor/log/:hostId", ginHelper.Handle(handler.MonitorLog))            // 查看日志
	api.GET("/alert/count/:hostId", ginHelper.Handle(handler.AlertCount))            // 获取报警通知数量
	api.GET("/slave/info", ginHelper.Handle(handler.MonitorSlaveInfo))               // 获取平台端系统信息
	api.POST("/conf/save", ginHelper.Handle(handler.MonitorConfSave))                // 设置监测器
	// 测试
	api.GET("/test/case1", ginHelper.Handle(handler.Case1))
	api.GET("/test/case2", ginHelper.Handle(handler.Case2))

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
