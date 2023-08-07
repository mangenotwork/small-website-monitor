package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"small-website-monitor/routers"
)

func main() {
	conf.InitConf("./conf/")
	//mysqlClient.InitMysqlGorm()     // 初始化 mysql
	gin.SetMode(gin.DebugMode)
	s := routers.Routers()
	s.Run(":" + conf.Conf.Default.HttpServer.Prod)
}
