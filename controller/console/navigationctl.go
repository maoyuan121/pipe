package console

import (
	"net/http"
	"strconv"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 获取导航集合
func GetNavigationsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	navigationModels, pagination := service.Navigation.ConsoleGetNavigations(util.GetPage(c), session.BID)

	var navigations []*ConsoleNavigation
	for _, navigationModel := range navigationModels {
		comment := &ConsoleNavigation{
			ID:         navigationModel.ID,
			Title:      navigationModel.Title,
			URL:        navigationModel.URL,
			IconURL:    navigationModel.IconURL,
			OpenMethod: navigationModel.OpenMethod,
			Number:     navigationModel.Number,
		}

		navigations = append(navigations, comment)
	}

	data := map[string]interface{}{}
	data["navigations"] = navigations
	data["pagination"] = pagination
	result.Data = data
}

// 获取指定导航的详细信息
func GetNavigationAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr

		return
	}

	data := service.Navigation.ConsoleGetNavigation(uint64(id))
	if nil == data {
		result.Code = util.CodeErr

		return
	}

	result.Data = data
}

// 删除一个导航
func RemoveNavigationAction(c *gin.Context) {
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

	if err := service.Navigation.RemoveNavigation(uint64(id), blogID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 更新导航
func UpdateNavigationAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}

	navigation := &model.Navigation{Model: model.Model{ID: uint64(id)}}
	if err := c.BindJSON(navigation); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update navigation request failed"

		return
	}

	session := util.GetSession(c)
	navigation.BlogID = session.BID

	if err := service.Navigation.UpdateNavigation(navigation); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 新建导航
func AddNavigationAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	navigation := &model.Navigation{}
	if err := c.BindJSON(navigation); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add navigation request failed"

		return
	}

	navigation.BlogID = session.BID
	if err := service.Navigation.AddNavigation(navigation); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}
