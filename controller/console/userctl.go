package console

import (
	"net/http"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 新建用户
func AddUserAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add user request failed"

		return
	}

	name := arg["name"].(string)
	user := service.User.GetUserByName(name)
	if nil == user {
		result.Code = util.CodeErr
		result.Msg = "the user should login first"

		return
	}

	session := util.GetSession(c)
	if err := service.User.AddUserToBlog(user.ID, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}
}

// 获取一个博客的所有作者
func GetUsersAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var users []*ConsoleUser
	userModels, pagination := service.User.GetBlogUsers(util.GetPage(c), session.BID)
	for _, userModel := range userModels {
		userBlog := service.User.GetUserBlog(userModel.ID, session.BID)
		users = append(users, &ConsoleUser{
			ID:           userModel.ID,
			Name:         userModel.Name,
			Nickname:     userModel.Nickname,
			Role:         userBlog.UserRole,
			URL:          blogURLSetting.Value + util.PathAuthors + "/" + userModel.Name,
			AvatarURL:    userModel.AvatarURL,
			ArticleCount: userBlog.UserArticleCount,
		})
	}

	result.Data = map[string]interface{}{
		"users":      users,
		"pagination": pagination,
	}
}
