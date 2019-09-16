package controller

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/gin-gonic/gin"
)

func showIndexAction(c *gin.Context) {
	t, err := template.ParseFiles("console/dist/index.html")
	if nil != err {
		logger.Errorf("load index page failed: " + err.Error())
		c.String(http.StatusNotFound, "load index page failed")

		return
	}

	t.Execute(c.Writer, nil)
}

func showStartPageAction(c *gin.Context) {
	t, err := template.ParseFiles("console/dist/start/index.html")
	if nil != err {
		logger.Errorf("load start page failed: " + err.Error())
		c.String(http.StatusNotFound, "load start page failed")

		return
	}

	t.Execute(c.Writer, nil)
}

// 系统信息
func showPlatInfoAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	data := map[string]interface{}{}
	data["version"] = model.Version
	data["database"] = service.Database()
	data["mode"] = model.Conf.RuntimeMode
	data["server"] = model.Conf.Server
	data["staticServer"] = model.Conf.StaticServer
	data["staticResourceVer"] = model.Conf.StaticResourceVersion

	result.Data = data
}

// 展示最近的 10 个博客
func showTopBlogsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	blogs := service.User.GetTopBlogs(10)
	for _, blog := range blogs {
		blog.ID = 0
		blog.UserID = 0
		blog.UserRole = 0
	}

	result.Data = blogs
}

// 展示 manifest 信息
func showManifestAction(c *gin.Context) {
	data, err := ioutil.ReadFile(filepath.FromSlash("theme/js/manifest.json"))
	if nil != err {
		notFound(c)

		return
	}

	manifest := string(data)
	manifest = strings.ReplaceAll(manifest, "{server}", model.Conf.Server)

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.Write([]byte(manifest))
}
