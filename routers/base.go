package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"small-website-monitor/handler"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func Routers() *gin.Engine {
	Router.StaticFS("/static", http.Dir("./static"))

	//模板
	// 自定义模板方法
	//Router.SetFuncMap(template.FuncMap{
	//	"formatAsDate": FormatAsDate,
	//})

	Router.Delims("{[", "]}")

	Router.LoadHTMLGlob("views/**/*")

	Page()

	return Router
}

func Page() {
	Router.GET("/", handler.LoginPage)
	Router.GET("/home", handler.HomePage)
}
