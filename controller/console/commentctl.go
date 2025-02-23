package console

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 获取评论集合
func GetCommentsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	commentModels, pagination := service.Comment.ConsoleGetComments(c.Query("key"), util.GetPage(c), session.BID)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var comments []*ConsoleComment
	for _, commentModel := range commentModels {
		article := service.Article.ConsoleGetArticle(commentModel.ArticleID)
		articleAuthor := service.User.GetUser(article.AuthorID)
		consoleArticleAuthor := &ConsoleAuthor{
			URL:       blogURLSetting.Value + util.PathAuthors + "/" + articleAuthor.Name,
			Name:      articleAuthor.Name,
			AvatarURL: articleAuthor.AvatarURL,
		}

		author := &ConsoleAuthor{}
		if model.SyncCommentAuthorID == commentModel.AuthorID {
			author.URL = commentModel.AuthorURL
			author.Name = commentModel.AuthorName
			author.AvatarURL = commentModel.AuthorAvatarURL
		} else {
			commentAuthor := service.User.GetUser(commentModel.AuthorID)
			commentAuthorBlog := service.User.GetOwnBlog(commentModel.AuthorID)
			author.URL = service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, commentAuthorBlog.ID).Value + util.PathAuthors + "/" + commentAuthor.Name
			author.Name = commentAuthor.Name
			author.AvatarURL = commentAuthor.AvatarURL
		}

		page := service.Comment.GetCommentPage(commentModel.ArticleID, commentModel.ID, commentModel.BlogID)
		mdResult := util.Markdown(commentModel.Content)
		comment := &ConsoleComment{
			ID:            commentModel.ID,
			Author:        author,
			ArticleAuthor: consoleArticleAuthor,
			CreatedAt:     commentModel.CreatedAt.Format("2006-01-02"),
			Title:         article.Title,
			Content:       template.HTML(mdResult.ContentHTML),
			URL:           blogURLSetting.Value + article.Path + "?p=" + strconv.Itoa(page) + "#pipeComment" + strconv.Itoa(int(commentModel.ID)),
		}

		comments = append(comments, comment)
	}

	data := map[string]interface{}{}
	data["comments"] = comments
	data["pagination"] = pagination
	result.Data = data
}

// 删除指定的评论
func RemoveCommentAction(c *gin.Context) {
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

	if err := service.Comment.RemoveComment(id, blogID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 批量删除批量
func RemoveCommentsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses batch remove comments request failed"

		return
	}

	session := util.GetSession(c)
	blogID := session.BID
	ids := arg["ids"].([]interface{})
	for _, id := range ids {
		if err := service.Comment.RemoveComment(uint64(id.(float64)), blogID); nil != err {
			logger.Errorf("remove comment failed: " + err.Error())
		}
	}
}
