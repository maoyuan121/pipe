package controller

import (
	"net/http"
	"strconv"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// addSymCommentAction adds a comment come from Sym. Sees https://hacpai.com/article/1457158841475 for more details.
func addSymCommentAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add comment request failed"

		return
	}

	logger.Infof("Add a comment from Sym: %+v", arg)

	client := arg["client"].(map[string]interface{})
	b3Key := client["userB3Key"].(string)
	articleAuthorName := client["userName"].(string)
	articleAuthor := service.User.GetUserByName(articleAuthorName)
	if articleAuthor.B3Key != b3Key {
		result.Code = util.CodeErr
		result.Msg = "wrong B3 Key"

		return
	}

	requestCmt := arg["comment"].(map[string]interface{})
	articleId, err := strconv.ParseUint(requestCmt["articleId"].(string), 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add comment request failed"

		return
	}

	article := service.Article.ConsoleGetArticle(articleId)
	if nil == article {
		result.Code = util.CodeErr
		result.Msg = "not found the specified article"

		return
	}

	blogID := getBlogID(c)
	commentableSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicCommentable, blogID)
	if "true" != commentableSetting.Value || !article.Commentable {
		result.Code = util.CodeErr
		result.Msg = "not allow comment"

		return
	}

	comment := &model.Comment{
		BlogID:          blogID,
		ArticleID:       articleId,
		AuthorID:        model.SyncCommentAuthorID,
		Content:         requestCmt["content"].(string),
		IP:              "",
		UserAgent:       "",
		AuthorName:      requestCmt["authorName"].(string),
		AuthorURL:       requestCmt["authorURL"].(string),
		AuthorAvatarURL: requestCmt["authorAvatarURL"].(string),
	}

	if err := service.Comment.AddComment(comment); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}

	if err := service.Comment.UpdatePushedAt(comment); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// addSymArticleAction adds an article come from Sym. Sees https://hacpai.com/article/1457158841475 for more details.
func addSymArticleAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add article request failed"

		return
	}

	logger.Infof("Add an article from Sym: %+v", arg)

	client := arg["client"].(map[string]interface{})
	b3Key := client["userB3Key"].(string)
	articleAuthorName := client["userName"].(string)
	articleAuthor := service.User.GetUserByName(articleAuthorName)
	if articleAuthor.B3Key != b3Key {
		result.Code = util.CodeErr
		result.Msg = "wrong B3 Key"

		return
	}

	requestArticle := arg["article"].(map[string]interface{})
	articleId, err := strconv.ParseUint(requestArticle["id"].(string), 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses add article request failed"

		return
	}

	article := service.Article.ConsoleGetArticle(articleId)
	if nil == article {
		blogID := getBlogID(c)
		article = &model.Article{
			BlogID:      blogID,
			AuthorID:    articleAuthor.ID,
			Title:       requestArticle["title"].(string),
			Tags:        requestArticle["tags"].(string),
			Content:     requestArticle["content"].(string),
			Commentable: true,
		}
		article.ID = articleId

		if err := service.Article.AddArticle(article); nil != err {
			result.Code = util.CodeErr
			result.Msg = err.Error()

			return
		}
	} else {
		article.Title = requestArticle["title"].(string)
		article.Tags = requestArticle["tags"].(string)
		article.Content = requestArticle["content"].(string)

		if err := service.Article.UpdateArticle(article); nil != err {
			result.Code = util.CodeErr
			result.Msg = err.Error()

			return
		}
	}

	if err := service.Article.UpdatePushedAt(article); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}
