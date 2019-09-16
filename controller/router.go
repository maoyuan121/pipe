package controller

import (
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/controller/console"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/theme"
	"github.com/b3log/pipe/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// 返回一个 gin engine， 将请求 URL 绑定到控制器
func MapRoutes() *gin.Engine {
	ret := gin.New()
	// 为模板添加函数
	ret.SetFuncMap(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("len(values) is " + strconv.Itoa(len(values)%2))
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"minus":    func(a, b int) int { return a - b },
		"mod":      func(a, b int) int { return a % b },
		"noescape": func(s string) template.HTML { return template.HTML(s) },
	})

	// 如果是开发模式使用 Logger
	if "dev" == model.Conf.RuntimeMode {
		ret.Use(gin.Logger())
	}
	ret.Use(gin.Recovery())

	// 定义 session cookie 存储
	store := cookie.NewStore([]byte(model.Conf.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   model.Conf.SessionMaxAge,
		Secure:   strings.HasPrefix(model.Conf.Server, "https"),
		HttpOnly: true,
	})
	ret.Use(sessions.Sessions("pipe", store))


	ret.GET(util.PathPlatInfo, showPlatInfoAction)
	ret.GET(util.PathSitemap, outputSitemapAction)

	// api
	api := ret.Group(util.PathAPI)
	api.POST("/logout", logoutAction)
	api.Any("/hp/*apis", util.HacPaiAPI())
	api.GET("/status", getStatusAction)
	api.GET("/check-version", console.CheckVersionAction)
	api.GET("/blogs/top", showTopBlogsAction)
	api.GET("/oauth/github/redirect", redirectGitHubLoginAction)
	api.GET("/oauth/github/callback", githubCallbackAction)

	// api/管理页面
	consoleGroup := api.Group("/console")
	// api/管理页面都需要登录
	consoleGroup.Use(console.LoginCheck)
	// 如果是开发模式定义生成示例文章的路由
	if "dev" == model.Conf.RuntimeMode {
		consoleGroup.GET("/dev/articles/gen", console.GenArticlesAction)
	}
	consoleGroup.GET("/themes", console.GetThemesAction)
	consoleGroup.PUT("/themes/:id", console.UpdateThemeAction)
	consoleGroup.GET("/tags", console.GetTagsAction)
	consoleGroup.GET("/taglist", console.GetTagsPageAction)
	consoleGroup.DELETE("/tags/:id", console.RemoveTagsAction)
	consoleGroup.POST("/articles", console.AddArticleAction)
	consoleGroup.GET("/upload/token", console.UploadTokenAction)
	consoleGroup.POST("/articles/batch-delete", console.RemoveArticlesAction)
	consoleGroup.GET("/articles", console.GetArticlesAction)
	consoleGroup.GET("/articles/:id", console.GetArticleAction)
	consoleGroup.GET("/articles/:id/push", console.PushArticle2RhyAction)
	consoleGroup.DELETE("/articles/:id", console.RemoveArticleAction)
	consoleGroup.PUT("/articles/:id", console.UpdateArticleAction)
	consoleGroup.GET("/comments", console.GetCommentsAction)
	consoleGroup.POST("/comments/batch-delete", console.RemoveCommentsAction)
	consoleGroup.DELETE("/comments/:id", console.RemoveCommentAction)
	consoleGroup.GET("/categories", console.GetCategoriesAction)
	consoleGroup.POST("/categories", console.AddCategoryAction)
	consoleGroup.DELETE("/categories/:id", console.RemoveCategoryAction)
	consoleGroup.GET("/categories/:id", console.GetCategoryAction)
	consoleGroup.PUT("/categories/:id", console.UpdateCategoryAction)
	consoleGroup.GET("/navigations", console.GetNavigationsAction)
	consoleGroup.GET("/navigations/:id", console.GetNavigationAction)
	consoleGroup.PUT("/navigations/:id", console.UpdateNavigationAction)
	consoleGroup.POST("/navigations", console.AddNavigationAction)
	consoleGroup.DELETE("/navigations/:id", console.RemoveNavigationAction)
	consoleGroup.GET("/users", console.GetUsersAction)
	consoleGroup.POST("/users", console.AddUserAction)
	consoleGroup.GET("/thumbs", console.GetArticleThumbsAction)
	consoleGroup.POST("/markdown", console.MarkdownAction)
	consoleGroup.POST("/import/md", console.ImportMarkdownAction)
	consoleGroup.GET("/export/md", console.ExportMarkdownAction)
	// consoleGroup.POST("/blogs/switch/:id", console.BlogSwitchAction)

	// api/管理页面/设置
	consoleSettingsGroup := consoleGroup.Group("/settings")
	consoleSettingsGroup.GET("/basic", console.GetBasicSettingsAction)
	consoleSettingsGroup.PUT("/basic", console.UpdateBasicSettingsAction)
	consoleSettingsGroup.GET("/preference", console.GetPreferenceSettingsAction)
	consoleSettingsGroup.PUT("/preference", console.UpdatePreferenceSettingsAction)
	consoleSettingsGroup.GET("/sign", console.GetSignSettingsAction)
	consoleSettingsGroup.PUT("/sign", console.UpdateSignSettingsAction)
	consoleSettingsGroup.GET("/i18n", console.GetI18nSettingsAction)
	consoleSettingsGroup.PUT("/i18n", console.UpdateI18nSettingsAction)
	consoleSettingsGroup.GET("/feed", console.GetFeedSettingsAction)
	consoleSettingsGroup.PUT("/feed", console.UpdateFeedSettingsAction)
	consoleSettingsGroup.GET("/third-stat", console.GetThirdStatisticSettingsAction)
	consoleSettingsGroup.PUT("/third-stat", console.UpdateThirdStatisticSettingsAction)
	consoleSettingsGroup.GET("/ad", console.GetAdSettingsAction)
	consoleSettingsGroup.PUT("/ad", console.UpdateAdSettingsAction)
	consoleSettingsGroup.GET("/account", console.GetAccountAction)
	consoleSettingsGroup.PUT("/account", console.UpdateAccountAction)

	// 静态文件路由
	ret.StaticFile(util.PathFavicon, "console/static/favicon.ico")
	ret.StaticFile(util.PathManifest, "console/static/manifest.json")

	// 资源文件目录路由
	ret.Static(util.PathTheme+"/scss", "theme/scss")
	ret.Static(util.PathTheme+"/js", "theme/js")
	ret.Static(util.PathTheme+"/images", "theme/images")
	// 静态资源文件路由
	ret.StaticFile("/sw.min.js", "theme/sw.min.js")
	ret.StaticFile("/halt.html", "theme/halt.html")

	// 主题资源文件
	for _, theme := range theme.Themes {
		themePath := "theme/x/" + theme
		ret.Static("/theme/x/"+theme+"/css", themePath+"/css")
		ret.Static("/theme/x/"+theme+"/js", themePath+"/js")
		ret.Static("/theme/x/"+theme+"/images", themePath+"/images")
		ret.StaticFile("/theme/x/"+theme+"/thumbnail.jpg", themePath+"/thumbnail.jpg")
	}

	// 主题模板
	themeTemplates, err := filepath.Glob("theme/x/*/*.html")
	if nil != err {
		logger.Fatal("load theme templates failed: " + err.Error())
	}
	themeTemplates = append(themeTemplates, "theme/search/index.html")
	commentTemplates, err := filepath.Glob("theme/comment/*.html")
	if nil != err {
		logger.Fatal("load comment templates failed: " + err.Error())
	}
	headTemplates, err := filepath.Glob("theme/head/*.html")
	if nil != err {
		logger.Fatal("load head templates failed: " + err.Error())
	}
	templates := append(themeTemplates, commentTemplates...)
	templates = append(templates, headTemplates...)
	ret.LoadHTMLFiles(templates...)

	themeGroup := ret.Group(util.PathBlogs + "/:username")
	themeGroup.Use(fillUser, pjax, resolveBlog)
	themeGroup.GET("", showArticlesAction)
	themeGroup.Any("/*path", routePath)

	// admin 前缀
	adminPagesGroup := ret.Group(util.PathAdmin)
	adminPagesGroup.Use(fillUser)
	adminPagesGroup.GET("", console.ShowAdminPagesAction)
	adminPagesGroup.GET("/*path", console.ShowAdminPagesAction)

	// / 前缀
	indexGroup := ret.Group("")
	indexGroup.Use(fillUser)
	indexGroup.GET("", showIndexAction)

	initGroup := ret.Group(util.PathInit)
	initGroup.Use(fillUser)
	initGroup.GET("", showStartPageAction)

	ret.Static(util.PathConsoleDist, "console/dist")
	ret.StaticFile(util.PathChangelogs, "changelogs.html")
	ret.StaticFile(util.PathRobots, "theme/robots.txt")

	// 没有匹配到路由的话
	ret.NoRoute(func(c *gin.Context) {
		notFound(c)
	})

	return ret
}

func routePath(c *gin.Context) {
	path := c.Param("path")

	switch path {
	case util.PathActivities:
		showActivitiesAction(c)

		return
	case util.PathArchives:
		showArchivesAction(c)

		return
	case util.PathAuthors:
		showAuthorsAction(c)

		return
	case util.PathCategories:
		showCategoriesAction(c)

		return
	case util.PathTags:
		showTagsAction(c)

		return
	case util.PathComments:
		addCommentAction(c)

		return
	case util.PathAtom:
		outputAtomAction(c)

		return
	case util.PathRSS:
		outputRSSAction(c)

		return
	case util.PathSearch:
		searchAction(c)

		return
	case util.PathOpensearch:
		showOpensearchAction(c)

		return
	case util.PathManifest:
		showManifestAction(c)

		return
	case util.PathAPIsSymComment:
		addSymCommentAction(c)

		return
	case util.PathAPIsSymArticle:
		addSymArticleAction(c)

		return
	case "/api/oauth/github/callback":
		githubCallbackAction(c)

		return
	case "/api/markdown":
		console.MarkdownAction(c)

		return
	}

	if strings.Contains(path, util.PathArchives+"/") {
		showArchiveArticlesAction(c)

		return
	}
	if strings.Contains(path, util.PathAuthors+"/") {
		showAuthorArticlesAction(c)

		return
	}
	if strings.Contains(path, util.PathCategories+"/") {
		showCategoryArticlesArticlesAction(c)

		return
	}
	if strings.Contains(path, util.PathTags+"/") {
		showTagArticlesAction(c)

		return
	}
	if strings.Contains(path, util.PathComments+"/") {
		if "DELETE" == c.Request.Method {
			commentID := strings.Split(c.Request.RequestURI, util.PathComments+"/")[1]
			c.Params = append(c.Params, gin.Param{Key: "id", Value: commentID})
			console.RemoveCommentAction(c)
		} else {
			getRepliesAction(c)
		}

		return
	}
	if "/" == path {
		showArticlesAction(c)

		return
	}

	logger.Tracef("can't handle path [" + path + "]")
	notFound(c)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}