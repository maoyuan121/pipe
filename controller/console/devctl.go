package console

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 为了测试或者 demo，生成一些文章
func GenArticlesAction(c *gin.Context) {
	session := util.GetSession(c)

	for i := 0; i < 100; i++ {
		article := &model.Article{
			AuthorID: session.UID,
			Title:    "title " + strconv.Itoa(i) + "_" + strconv.Itoa(rand.Int()),
			Tags:     "开发生成",
			Content:  "开发生成",
			BlogID:   session.BID,
		}
		if err := service.Article.AddArticle(article); nil != err {
			logger.Errorf("generate article failed: " + err.Error())
		}
	}

	c.Redirect(http.StatusTemporaryRedirect, model.Conf.Server)
}
