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
	"github.com/vinta/pangu"
)

// 一个博客的作者列表
func showAuthorsAction(c *gin.Context) {
	dataModel := getDataModel(c)
	blogID := getBlogID(c)
	locale := getLocale(c)
	var themeAuthors []*model.ThemeAuthor
	authorModels, _ := service.User.GetBlogUsers(1, blogID)
	for _, authorModel := range authorModels {
		userBlog := service.User.GetUserBlog(authorModel.ID, blogID)
		author := &model.ThemeAuthor{
			Name:         authorModel.Name,
			URL:          getBlogURL(c) + util.PathAuthors + "/" + authorModel.Name,
			ArticleCount: userBlog.UserArticleCount,
			AvatarURL:    authorModel.AvatarURL,
		}
		themeAuthors = append(themeAuthors, author)
	}

	dataModel["Authors"] = themeAuthors
	dataModel["Title"] = i18n.GetMessage(locale, "team") + " - " + dataModel["Title"].(string)

	c.HTML(http.StatusOK, getTheme(c)+"/authors.html", dataModel)
}

// 获取一个人在某个博客下写的所有文章
func showAuthorArticlesAction(c *gin.Context) {
	authorName := strings.SplitAfter(c.Request.URL.Path, util.PathAuthors+"/")[1]
	author := service.User.GetUserByName(authorName)
	if nil == author {
		notFound(c)

		return
	}

	page := util.GetPage(c)
	blogID := getBlogID(c)
	dataModel := getDataModel(c)
	session := util.GetSession(c)
	locale := getLocale(c)
	articleListStyleSetting := service.Setting.GetSetting(model.SettingCategoryPreference, model.SettingNamePreferenceArticleListStyle, blogID)
	articleModels, pagination := service.Article.GetAuthorArticles(author.ID, page, blogID)
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
	userBlog := service.User.GetUserBlog(author.ID, blogID)
	dataModel["Author"] = &model.ThemeAuthor{
		Name:         author.Name,
		ArticleCount: userBlog.UserArticleCount,
	}
	dataModel["Title"] = authorName + " - " + i18n.GetMessage(locale, "team") + " - " + dataModel["Title"].(string)

	c.HTML(http.StatusOK, getTheme(c)+"/author-articles.html", dataModel)
}
