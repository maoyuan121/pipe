package console

import (
	"net/http"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 检查是否登录
// 没登录的话返回 unauthenticated request
func LoginCheck(c *gin.Context) {
	session := util.GetSession(c)
	if 0 == session.UID {
		result := gulu.Ret.NewResult()
		result.Code = util.CodeAuthErr
		result.Msg = "unauthenticated request"
		c.AbortWithStatusJSON(http.StatusOK, result)

		return
	}

	c.Next()
}
