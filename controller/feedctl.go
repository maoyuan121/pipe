package controller

import (
	"strconv"

	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func outputAtomAction(c *gin.Context) {
	feed := generateFeed(c)

	feed.WriteAtom(c.Writer)
}

func outputRSSAction(c *gin.Context) {
	feed := generateFeed(c)

	feed.WriteRss(c.Writer)
}

func generateFeed(c *gin.Context) *feeds.Feed {
	blogID := getBlogID(c)

	feedOutputModeSetting := service.Setting.GetSetting(model.SettingCategoryFeed, model.SettingNameFeedOutputMode, blogID)
	blogTitleSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogTitle, blogID)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, blogID)
	blogSubtitleSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogSubtitle, blogID)
	ret := &feeds.Feed{
		Title:       blogTitleSetting.Value,
		Link:        &feeds.Link{Href: blogURLSetting.Value},
		Description: blogSubtitleSetting.Value,
	}

	var items []*feeds.Item
	articles, _ := service.Article.GetArticles("", 1, blogID)
	for _, article := range articles {
		mdResult := util.Markdown(article.Content)
		description := mdResult.AbstractText
		if strconv.Itoa(model.SettingFeedOutputModeValueFull) == feedOutputModeSetting.Value {
			description = mdResult.ContentHTML
		}
		user := service.User.GetUser(article.AuthorID)
		items = append(items, &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: blogURLSetting.Value + article.Path},
			Description: description,
			Author:      &feeds.Author{Name: user.Name},
			Created:     article.CreatedAt,
		})
	}
	ret.Items = items

	return ret
}
