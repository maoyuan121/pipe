package console

import (
	"math"
	"net/http"
	"strconv"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 获取标签集合
func GetTagsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var tags []*ConsoleTag
	tagModels := service.Tag.GetTags(math.MaxInt64, session.BID)
	for _, tagModel := range tagModels {
		tags = append(tags, &ConsoleTag{
			Title: tagModel.Title,
			URL:   blogURLSetting.Value + util.PathTags + "/" + tagModel.Title,
		})
	}

	result.Data = tags
}

// 获取标签集合，带分页
func GetTagsPageAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	tagModels, pagination := service.Tag.ConsoleGetTags(c.Query("key"), util.GetPage(c), session.BID)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var tags []*ConsoleTag
	for _, tagModel := range tagModels {
		item := &ConsoleTag{
			ID:    tagModel.ID,
			Title: tagModel.Title,
			URL:   blogURLSetting.Value + util.PathTags + "/" + tagModel.Title,
		}
		tags = append(tags, item)
	}
	data := map[string]interface{}{}
	data["tags"] = tags
	data["pagination"] = pagination
	result.Data = data
}

// 删除指定的一个标签，这个标签下必须没有文章
func RemoveTagsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}
	session := util.GetSession(c)
	blogID := session.BID

	if err := service.Tag.RemoveTag(id, blogID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}

}
