package controller

import (
	"github.com/b3log/pipe/service"
	"github.com/gin-gonic/gin"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

// sitemap 显示最近的 10 个文章最多的博客
func outputSitemapAction(c *gin.Context) {
	sm := stm.NewSitemap(1)
	sm.Create()

	blogs := service.User.GetTopBlogs(10)
	for _, blog := range blogs {
		sm.Add(stm.URL{{"loc", blog.URL}})
	}

	c.Writer.Write(sm.XMLContent())
}
