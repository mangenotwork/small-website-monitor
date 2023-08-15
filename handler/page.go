package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/utils"
	"net/http"
	"small-website-monitor/global"
)

func LoginPage(c *gin.Context) {
	token, _ := c.Cookie(global.UserToken)
	if token != "" {
		j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
		if err := j.ParseToken(token); err == nil {
			c.Redirect(http.StatusFound, "/home")
			return
		}
	}
	c.HTML(200, "login.html", gin.H{
		"csrf": ginHelper.FormSetCSRF(c.Request),
	})
}

func HomePage(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{})
}