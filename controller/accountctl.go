package controller

import (
	"net/http"

	"github.com/b3log/gulu"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 退出登录
func logoutAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: -1,
	})
	session.Clear()
	if err := session.Save(); nil != err {
		logger.Errorf("saves session failed: " + err.Error())
	}
}
