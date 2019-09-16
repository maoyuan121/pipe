package console

import (
	"net/http"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/theme"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 切换主题
func UpdateThemeAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	theme := c.Param("id")
	session := util.GetSession(c)

	settings := []*model.Setting{
		{
			Category: model.SettingCategoryTheme,
			Name:     model.SettingNameThemeName,
			Value:    theme,
			BlogID:   session.BID,
		},
	}
	if err := service.Setting.UpdateSettings(model.SettingCategoryTheme, settings, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取所有主题
func GetThemesAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)

	currentID := theme.Themes[0]
	themeNameSetting := service.Setting.GetSetting(model.SettingCategoryTheme, model.SettingNameThemeName, session.BID)
	if nil == themeNameSetting {
		logger.Errorf("not found theme name setting")
	} else {
		currentID = themeNameSetting.Value
	}

	var themes []*ConsoleTheme
	for _, themeName := range theme.Themes {
		consoleTheme := &ConsoleTheme{
			Name:         themeName,
			ThumbnailURL: model.Conf.Server + "/theme/x/" + themeName + "/thumbnail.jpg",
		}

		themes = append(themes, consoleTheme)
	}

	result.Data = map[string]interface{}{
		"currentId": currentID,
		"themes":    themes,
	}
}
