package controller

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/b3log/pipe/i18n"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
	"github.com/vinta/pangu"
)

// 分类列表
func showCategoriesAction(c *gin.Context) {
	dataModel := getDataModel(c)
	blogID := getBlogID(c)
	locale := getLocale(c)
	categoryModels := service.Category.GetCategories(math.MaxInt8, blogID)
	var themeCategories []*model.ThemeCategory
	for _, categoryModel := range categoryModels {
		var themeTags []*model.ThemeTag
		tagStrs := strings.Split(categoryModel.Tags, ",")
		for _, tagTitle := range tagStrs {
			themeTag := &model.ThemeTag{
				Title: tagTitle,
				URL:   getBlogURL(c) + util.PathTags + "/" + tagTitle,
			}
			themeTags = append(themeTags, themeTag)
		}

		themeCategory := &model.ThemeCategory{
			Title:        categoryModel.Title,
			URL:          getBlogURL(c) + util.PathCategories + categoryModel.Path,
			Description:  categoryModel.Description,
			Tags:         themeTags,
			ArticleCount: 8, // TODO: category article count
		}
		themeCategories = append(themeCategories, themeCategory)
	}

	dataModel["Categories"] = themeCategories
	dataModel["Title"] = i18n.GetMessage(locale, "categories") + " - " + dataModel["Title"].(string)

	c.HTML(http.StatusOK, getTheme(c)+"/categories.html", dataModel)
}

// 显示某个分类下的文章列表页面
func showCategoryArticlesArticlesAction(c *gin.Context) {
	dataModel := getDataModel(c)
	blogID := getBlogID(c)
	locale := getLocale(c)
	session := util.GetSession(c)
	page := util.GetPage(c)
	categoryPath := strings.SplitAfter(c.Request.URL.Path, util.PathCategories)[1]
	categoryModel := service.Category.GetCategoryByPath(categoryPath, blogID)
	if nil == categoryModel {
		notFound(c)

		return
	}
	articleListStyleSetting := service.Setting.GetSetting(model.SettingCategoryPreference, model.SettingNamePreferenceArticleListStyle, blogID)
	articleModels, pagination := service.Article.GetCategoryArticles(categoryModel.ID, page, blogID)
	var articles []*model.ThemeArticle
	for _, articleModel := range articleModels {
		var themeTags []*model.ThemeTag
		tagStrs := strings.Split(articleModel.Tags, ",")
		for _, tagStr := range tagStrs {
			themeTag := &model.ThemeTag{
				Title: tagStr,
				URL:   getBlogURL(c) + util.PathTags + "/" + tagStr,
			}
			themeTags = append(themeTags, themeTag)
		}

		authorModel := service.User.GetUser(articleModel.AuthorID)
		if nil == authorModel {
			logger.Errorf("not found author of article [id=%d, authorID=%d]", articleModel.ID, articleModel.AuthorID)

			continue
		}

		author := &model.ThemeAuthor{
			Name:      authorModel.Name,
			URL:       getBlogURL(c) + util.PathAuthors + "/" + authorModel.Name,
			AvatarURL: authorModel.AvatarURL,
		}

		mdResult := util.Markdown(articleModel.Content)
		abstract := template.HTML("")
		thumbnailURL := mdResult.ThumbURL
		if strconv.Itoa(model.SettingPreferenceArticleListStyleValueTitleAbstract) == articleListStyleSetting.Value {
			abstract = template.HTML(mdResult.AbstractText)
		}
		if "" != articleModel.Abstract {
			abstract = template.HTML(articleModel.Abstract)
		}
		if strconv.Itoa(model.SettingPreferenceArticleListStyleValueTitleContent) == articleListStyleSetting.Value {
			abstract = template.HTML(mdResult.ContentHTML)
			thumbnailURL = ""
		}
		article := &model.ThemeArticle{
			ID:             articleModel.ID,
			Abstract:       abstract,
			Author:         author,
			CreatedAt:      articleModel.CreatedAt.Format("2006-01-02"),
			CreatedAtYear:  articleModel.CreatedAt.Format("2006"),
			CreatedAtMonth: articleModel.CreatedAt.Format("01"),
			CreatedAtDay:   articleModel.CreatedAt.Format("02"),
			Title:          pangu.SpacingText(articleModel.Title),
			Tags:           themeTags,
			URL:            getBlogURL(c) + articleModel.Path,
			Topped:         articleModel.Topped,
			ViewCount:      articleModel.ViewCount,
			CommentCount:   articleModel.CommentCount,
			ThumbnailURL:   thumbnailURL,
			Editable:       session.UID == authorModel.ID,
		}

		articles = append(articles, article)
	}
	dataModel["Articles"] = articles
	dataModel["Pagination"] = pagination
	dataModel["Category"] = &model.ThemeCategory{
		Title:        categoryModel.Title,
		ArticleCount: pagination.RecordCount,
	}
	dataModel["Title"] = categoryModel.Title + " - " + i18n.GetMessage(locale, "categories") + " - " + dataModel["Title"].(string)

	c.HTML(http.StatusOK, getTheme(c)+"/category-articles.html", dataModel)
}
