package console

import (
	"net/http"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 更新当前用户的一些信息
func UpdateAccountAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update account request failed"

		return
	}

	b3Key := arg["b3key"].(string)
	avatarURL := arg["avatarURL"].(string)

	session := util.GetSession(c)
	user := service.User.GetUserByName(session.UName)
	user.B3Key = b3Key
	user.AvatarURL = avatarURL
	if err := service.User.UpdateUser(user); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}
	session.UB3Key = b3Key
	session.UAvatar = avatarURL
	session.Save(c)
}

// 获取当前用户信息
func GetAccountAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	data := map[string]interface{}{}
	data["name"] = session.UName
	data["avatarURL"] = session.UAvatar
	data["b3Key"] = session.UB3Key

	result.Data = data
}
