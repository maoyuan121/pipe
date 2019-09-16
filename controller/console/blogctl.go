package console

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

// 切换博客
func BlogSwitchAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	blogID, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr

		return
	}

	session := util.GetSession(c)
	userID := session.UID

	userBlogs := service.User.GetUserBlogs(userID)
	if 1 > len(userBlogs) {
		result.Code = util.CodeErr
		result.Msg = "switch blog failed"

		return
	}

	role := -1
	for _, userBlog := range userBlogs {
		if userBlog.ID == uint64(blogID) {
			role = userBlog.UserRole

			break
		}
	}

	if -1 == role {
		result.Code = util.CodeErr
		result.Msg = "switch blog failed"

		return
	}

	result.Data = role

	session.URole = role
	session.BID = uint64(blogID)
	session.Save(c)
}

// 检查 pipe 版本信息，是否可以更新
func CheckVersionAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	rhyResult := map[string]interface{}{}
	request := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	_, _, errs := request.Get("https://rhythm.b3log.org/version/pipe/latest/"+model.Version).
		Set("User-Agent", model.UserAgent).Timeout(30*time.Second).
		Retry(3, 5*time.Second).EndStruct(&rhyResult)
	if nil != errs {
		result.Code = util.CodeErr
		result.Msg = errs[0].Error()

		return
	}

	data := map[string]interface{}{}
	data["version"] = rhyResult["pipeVersion"]
	data["download"] = rhyResult["pipeDownload"]
	result.Data = data
}
