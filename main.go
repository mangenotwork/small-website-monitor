package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"small-website-monitor/business"
	"small-website-monitor/model"
	"small-website-monitor/routers"
)

func main() {
	conf.InitConf("./conf/")
	model.DB.Init()
	business.Monitor()
	gin.SetMode(gin.ReleaseMode)
	s := routers.Routers()
	s.Run(":" + conf.Conf.Default.HttpServer.Prod)

	//server := &http.Server{
	//	Addr:           ":" + conf.Conf.Default.HttpServer.Prod,
	//	Handler:        routers.Routers(),
	//	ReadTimeout:    10 * time.Second,
	//	WriteTimeout:   10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//server.ListenAndServe()
}
