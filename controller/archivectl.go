package controller

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/b3log/pipe/i18n"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 归档列表页面
func showArchivesAction(c *gin.Context) {
	dataModel := getDataModel(c)
	blogID := getBlogID(c)
	locale := getLocale(c)
	var themeArchives []*model.ThemeArchive
	archiveModels := service.Archive.GetArchives(blogID)
	for _, archiveModel := range archiveModels {
		archive := &model.ThemeArchive{
			Title:        i18n.GetMessagef(locale, "archiveYearMonth", archiveModel.Year, archiveModel.Month),
			URL:          getBlogURL(c) + util.PathArchives + "/" + archiveModel.Year + "/" + archiveModel.Month,
			ArticleCount: archiveModel.ArticleCount,
		}
		themeArchives = append(themeArchives, archive)
	}

	dataModel["Archives"] = themeArchives
	dataModel["Title"] = i18n.GetMessage(locale, "archives") + " - " + dataModel["Title"].(string)

	c.HTML(http.StatusOK, getTheme(c)+"/archives.html", dataModel)
}

// 显示某一个归档下面的所有文章
func showArchiveArticlesAction(c *gin.Context) {
	dataModel := getDataModel(c)
	blogID := getBlogID(c)
	locale := getLocale(c)
	session := util.GetSession(c)
	date := strings.SplitAfter(c.Request.URL.Path, util.PathArchives+"/")[1]
	year := strings.Split(date, "/")[0]
	month := strings.Split(date, "/")[1]
	archiveModel := service.Archive.GetArchive(year, month, blogID)
	if nil == archiveModel {
		notFound(c)

		return
	}
	articleListStyleSetting := service.Setting.GetSetting(model.SettingCategoryPreference, model.SettingNamePreferenceArticleListStyle, blogID)
	articleModels, pagination := service.Article.GetArchiveArticles(archiveModel.ID, util.GetPage(c), blogID)
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
			Title:          articleModel.Title,
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
	dataModel["Archive"] = &model.ThemeArchive{
		Title:        i18n.GetMessagef(locale, "archiveYearMonth", archiveModel.Year, archiveModel.Month),
		URL:          getBlogURL(c) + util.PathArchives + "/" + archiveModel.Year + "/" + archiveModel.Month,
		ArticleCount: archiveModel.ArticleCount,
	}
	dataModel["Title"] = i18n.GetMessagef(locale, "archiveYearMonth", archiveModel.Year, archiveModel.Month) +
		" - " + i18n.GetMessage(locale, "archives") + " - " + dataModel["Title"].(string)

	c.HTML(http.StatusOK, getTheme(c)+"/archive-articles.html", dataModel)
}
