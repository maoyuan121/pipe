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

// 更新分类
func UpdateCategoryAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}

	category := &model.Category{Model: model.Model{ID: uint64(id)}}
	if err := c.BindJSON(category); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update category request failed"

		return
	}

	session := util.GetSession(c)
	category.BlogID = session.BID

	if err := service.Category.UpdateCategory(category); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取指定的一个分类信息
func GetCategoryAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}

	data := service.Category.ConsoleGetCategory(id)
	if nil == data {
		result.Code = util.CodeErr

		return
	}

	result.Data = data
}

// 获取分类集合
func GetCategoriesAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	categoryModels, pagination := service.Category.ConsoleGetCategories(util.GetPage(c), session.BID)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var categories []*ConsoleCategory
	for _, categoryModel := range categoryModels {
		categories = append(categories, &ConsoleCategory{
			ID:          categoryModel.ID,
			Title:       categoryModel.Title,
			URL:         blogURLSetting.Value + util.PathCategories + categoryModel.Path,
			Description: categoryModel.Description,
			Number:      categoryModel.Number,
			Tags:        categoryModel.Tags,
		})
	}

	data := map[string]interface{}{}
	data["categories"] = categories
	data["pagination"] = pagination
	result.Data = data
}

// 新建分类
func AddCategoryAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	category := &model.Category{}
	if err := c.BindJSON(category); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add category request failed"

		return
	}

	category.BlogID = session.BID
	if err := service.Category.AddCategory(category); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 删除分类
func RemoveCategoryAction(c *gin.Context) {
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
	if err := service.Category.RemoveCategory(id, blogID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}
