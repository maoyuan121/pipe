package console

import (
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 获取基本设置
func GetBasicSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	settings := service.Setting.GetCategorySettings(model.SettingCategoryBasic, session.BID)
	data := map[string]interface{}{}
	for _, setting := range settings {
		if model.SettingNameBasicCommentable == setting.Name {
			v, err := strconv.ParseBool(setting.Value)
			if nil != err {
				logger.Errorf("value of basic setting [name=%s] must be \"true\" or \"false\"", setting.Name)
				data[setting.Name] = true
			} else {
				data[setting.Name] = v
			}
		} else {
			data[setting.Name] = setting.Value
		}
	}

	result.Data = data
}

// 更新基本设置
func UpdateBasicSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update basic settings request failed"

		return
	}

	session := util.GetSession(c)
	var basics []*model.Setting
	for k, v := range args {
		var value interface{}
		switch v.(type) {
		case bool:
			value = strconv.FormatBool(v.(bool))
		default:
			value = strings.TrimSpace(v.(string))
		}

		if model.SettingNameBasicBlogURL == k {
			blogURL := value.(string)
			if !strings.Contains(blogURL, "://") {
				blogURL = "http://" + blogURL
			}

			url, err := url.Parse(blogURL)
			if nil != err {
				result.Code = util.CodeErr
				result.Msg = "invalid URL format"

				return
			}

			blogURL = url.Scheme + "://" + url.Host
			if "" != url.Path {
				blogURL += path.Clean(url.Path)
			}
			value = blogURL
		}

		basic := &model.Setting{
			Category: model.SettingCategoryBasic,
			BlogID:   session.BID,
			Name:     k,
			Value:    value.(string),
		}
		basics = append(basics, basic)
	}

	if err := service.Setting.UpdateSettings(model.SettingCategoryBasic, basics, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取偏好设置
func GetPreferenceSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	settings := service.Setting.GetCategorySettings(model.SettingCategoryPreference, session.BID)
	data := map[string]interface{}{}
	for _, setting := range settings {
		if model.SettingNamePreferenceArticleListStyle != setting.Name {
			v, err := strconv.ParseInt(setting.Value, 10, 64)
			if nil != err {
				logger.Errorf("value of preference setting [name=%s] must be an integer", setting.Name)
				data[setting.Name] = 10
			} else {
				data[setting.Name] = v
			}
		} else {
			data[setting.Name] = setting.Value
		}
	}

	result.Data = data
}

// 更新偏好设置
func UpdatePreferenceSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update preference settings request failed"

		return
	}

	session := util.GetSession(c)
	var prefs []*model.Setting
	for k, v := range args {
		var value interface{}
		switch v.(type) {
		case float64:
			value = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		default:
			value = v.(string)
		}

		pref := &model.Setting{
			Category: model.SettingCategoryPreference,
			BlogID:   session.BID,
			Name:     k,
			Value:    value.(string),
		}
		prefs = append(prefs, pref)
	}

	if err := service.Setting.UpdateSettings(model.SettingCategoryPreference, prefs, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// GetSignSettingsAction gets sign settings.
func GetSignSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	signSetting := service.Setting.GetSetting(model.SettingCategorySign, model.SettingNameArticleSign, session.BID)
	result.Data = signSetting.Value
}

// UpdateSignSettingsAction updates sign settings.
func UpdateSignSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update sign settings request failed"

		return
	}

	session := util.GetSession(c)
	var signs []*model.Setting
	sign := &model.Setting{
		Category: model.SettingCategorySign,
		BlogID:   session.BID,
		Name:     model.SettingNameArticleSign,
		Value:    args["sign"].(string),
	}
	signs = append(signs, sign)

	if err := service.Setting.UpdateSettings(model.SettingCategorySign, signs, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取本地化设置
func GetI18nSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	settings := service.Setting.GetCategorySettings(model.SettingCategoryI18n, session.BID)
	data := map[string]interface{}{}
	for _, setting := range settings {
		data[setting.Name] = setting.Value
	}
	result.Data = data
}

// 更新本地化设置
func UpdateI18nSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update i18n settings request failed"

		return
	}

	session := util.GetSession(c)
	var i18ns []*model.Setting
	for k, v := range args {
		i18n := &model.Setting{
			Category: model.SettingCategoryI18n,
			BlogID:   session.BID,
			Name:     k,
			Value:    v.(string),
		}
		i18ns = append(i18ns, i18n)
	}

	if err := service.Setting.UpdateSettings(model.SettingCategoryI18n, i18ns, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取 feed 设置
func GetFeedSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	settings := service.Setting.GetCategorySettings(model.SettingCategoryFeed, session.BID)
	data := map[string]interface{}{}
	for _, setting := range settings {
		if model.SettingNameFeedOutputMode == setting.Name {
			v, err := strconv.ParseInt(setting.Value, 10, 64)
			if nil != err {
				logger.Errorf("value of feed setting [name=%s] must be an integer", setting.Name)
				data[setting.Name] = 20
			} else {
				data[setting.Name] = v
			}
		} else {
			data[setting.Name] = setting.Value
		}
	}
	result.Data = data
}

// 更新 feed 设置
func UpdateFeedSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update feed settings request failed"

		return
	}

	session := util.GetSession(c)
	var feeds []*model.Setting
	for k, v := range args {
		var value interface{}
		switch v.(type) {
		case float64:
			value = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		default:
			value = v.(string)
		}

		feed := &model.Setting{
			Category: model.SettingCategoryFeed,
			BlogID:   session.BID,
			Name:     k,
			Value:    value.(string),
		}
		feeds = append(feeds, feed)
	}

	if err := service.Setting.UpdateSettings(model.SettingCategoryFeed, feeds, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取第三方统计设置
func GetThirdStatisticSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	baiduStatisticSetting := service.Setting.GetSetting(model.SettingCategoryThirdStatistic, model.SettingNameThirdStatisticBaidu, session.BID)
	data := map[string]string{
		model.SettingNameThirdStatisticBaidu: baiduStatisticSetting.Value,
	}
	result.Data = data
}

// 更新第三方统计设置
func UpdateThirdStatisticSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update third statistic settings request failed"

		return
	}

	session := util.GetSession(c)
	var thridStatistics []*model.Setting
	baiduStatistic := &model.Setting{
		Category: model.SettingCategoryThirdStatistic,
		BlogID:   session.BID,
		Name:     model.SettingNameThirdStatisticBaidu,
		Value:    args["thirdStatisticBaidu"].(string),
	}
	thridStatistics = append(thridStatistics, baiduStatistic)

	if err := service.Setting.UpdateSettings(model.SettingCategoryThirdStatistic, thridStatistics, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// 获取广告设置
func GetAdSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	googleAdSenseArticleEmbedSetting := service.Setting.GetSetting(model.SettingCategoryAd, model.SettingNameAdGoogleAdSenseArticleEmbed, session.BID)
	data := map[string]string{
		model.SettingNameAdGoogleAdSenseArticleEmbed: googleAdSenseArticleEmbedSetting.Value,
	}
	result.Data = data
}

// 更新广告设置
func UpdateAdSettingsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	args := map[string]interface{}{}
	if err := c.BindJSON(&args); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses update ad settings request failed"

		return
	}

	googleAdSenseArticleEmbedVal := args["adGoogleAdSenseArticleEmbed"].(string)
	googleAdSenseArticleEmbedVal = strings.TrimSpace(googleAdSenseArticleEmbedVal)
	if !strings.HasPrefix(googleAdSenseArticleEmbedVal, "<ins ") || !strings.HasSuffix(googleAdSenseArticleEmbedVal, "</ins>") {
		result.Code = util.CodeErr
		result.Msg = "please just put <ins>....</ins> part"

		return
	}

	session := util.GetSession(c)
	var ads []*model.Setting
	googleAdSenseArticleEmbed := &model.Setting{
		Category: model.SettingCategoryAd,
		BlogID:   session.BID,
		Name:     model.SettingNameAdGoogleAdSenseArticleEmbed,
		Value:    googleAdSenseArticleEmbedVal,
	}
	ads = append(ads, googleAdSenseArticleEmbed)

	if err := service.Setting.UpdateSettings(model.SettingCategoryAd, ads, session.BID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}
