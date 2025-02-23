package cron

import (
	"time"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/dustin/go-humanize"
)

// RecommendArticles saves all recommend articles.
var RecommendArticles []*model.ThemeArticle

// 每 30 分钟执行一次， 取阅读数最多的 7 篇文章作为推荐文件
func refreshRecommendArticlesPeriodically() {
	go refreshRecommendArticles()

	go func() {
		for range time.Tick(time.Minute * 30) {
			refreshRecommendArticles()
		}
	}()
}

func refreshRecommendArticles() {
	defer gulu.Panic.Recover()

	size := 7
	articles := service.Article.GetPlatMostViewArticles(size)
	size = len(articles)
	indics := gulu.Rand.Ints(0, size, size)
	images := util.RandImages(size)
	indics = indics[:len(images)]
	var recommendations []*model.ThemeArticle
	for i, index := range indics {
		article := articles[index]
		authorModel := service.User.GetUser(article.AuthorID)
		if nil == authorModel {
			logger.Errorf("not found author of article [id=%d, authorID=%d]", article.ID, article.AuthorID)

			continue
		}

		blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, article.BlogID)
		blogURL := blogURLSetting.Value
		author := &model.ThemeAuthor{
			Name:      authorModel.Name,
			URL:       blogURL + util.PathAuthors + "/" + authorModel.Name,
			AvatarURL: authorModel.AvatarURL,
		}
		themeArticle := &model.ThemeArticle{
			Title:        article.Title,
			URL:          blogURL + article.Path,
			CreatedAt:    humanize.Time(article.CreatedAt),
			Author:       author,
			CommentCount: article.CommentCount,
			ViewCount:    article.ViewCount,
			ThumbnailURL: util.ImageSize(images[i], 280, 90),
		}
		recommendations = append(recommendations, themeArticle)
	}

	RecommendArticles = recommendations
}
