package model

// 配置
type Setting struct {
	Model
	Category string `sql:"index" gorm:"size:32" json:"category"` // 分类 eg: system, theme, basic, preference, sign, i18n ...
	Name     string `sql:"index" gorm:"size:64" json:"name"`     // 键
	Value    string `gorm:"type:text" json:"value"`              // 值
	BlogID   uint64 `sql:"index" json:"blogID"`                  // 所属博客 ID
}

// Setting names of category "system".
const (
	SettingCategorySystem = "system"

	SettingNameSystemVer = "systemVersion"
)

// Setting names of category "theme".
const (
	SettingCategoryTheme = "theme"

	SettingNameThemeName = "themeName"
)

// Setting names of category "basic".
const (
	SettingCategoryBasic = "basic"

	SettingNameBasicBlogURL         = "basicBlogURL"
	SettingNameBasicBlogSubtitle    = "basicBlogSubtitle"
	SettingNameBasicBlogTitle       = "basicBlogTitle"
	SettingNameBasicCommentable     = "basicCommentable"
	SettingNameBasicFooter          = "basicFooter"
	SettingNameBasicHeader          = "basicHeader" // Removed from UI since v1.1.0 caused by issue 54 (https://github.com/b3log/pipe/issues/54)
	SettingNameBasicNoticeBoard     = "basicNoticeBoard"
	SettingNameBasicMetaDescription = "basicMetaDescription"
	SettingNameBasicMetaKeywords    = "basicMetaKeywords"
	SettingNameBasicFaviconURL      = "basicFaviconURL"
	SettingNameBasicLogoURL         = "basicLogoURL"
)

// Setting values of category "basic".
const (
	SettingBasicFooterDefault           = "<!-- 这里可用于放置备案信息等，支持 Markdown -->"
	SettingBasicHeaderDefault           = "<!-- https://github.com/b3log/pipe/issues/54 -->"
	SettingBasicBasicNoticeBoardDefault = "<!-- 支持 Markdown -->本博客由 [Pipe](https://github.com/b3log/pipe) 强力驱动"
)

// Setting names of category "preference".
const (
	SettingCategoryPreference = "preference"

	SettingNamePreferenceArticleListPageSize        = "preferenceArticleListPageSize"
	SettingNamePreferenceArticleListWindowSize      = "preferenceArticleListWindowSize"
	SettingNamePreferenceArticleListStyle           = "preferenceArticleListStyle"
	SettingNamePreferenceMostCommentArticleListSize = "preferenceMostCommentArticleListSize"
	SettingNamePreferenceMostUseTagListSize         = "preferenceMostUseTagListSize"
	SettingNamePreferenceMostViewArticleListSize    = "preferenceMostViewArticleListSize"
	SettingNamePreferenceRecentCommentListSize      = "preferenceRecentCommentListSize"
	SettingNamePreferenceRecommendArticleListSize   = "preferenceRecommendArticleListSize"
)

// Setting values of category "preference".
const (
	SettingPreferenceArticleListStyleValueTitle         = 0
	SettingPreferenceArticleListStyleValueTitleAbstract = 1
	SettingPreferenceArticleListStyleValueTitleContent  = 2

	SettingPreferenceArticleListPageSizeDefault        = 20
	SettingPreferenceArticleListWindowSizeDefault      = 7
	SettingPreferenceArticleListStyleDefault           = SettingPreferenceArticleListStyleValueTitleAbstract
	SettingPreferenceMostCommentArticleListSizeDefault = 7
	SettingPreferenceMostUseTagListSizeDefault         = 15
	SettingPreferenceMostViewArticleListSizeDefault    = 15
	SettingPreferenceRecentCommentListSizeDefault      = 7
	SettingPreferenceRecommendArticleListSizeDefault   = 1
)

// Setting names of category "sign".
const (
	SettingCategorySign = "sign"

	SettingNameArticleSign = "signArticle"
)

// Setting values of category "sign".
const (
	SettingArticleSignDefault = "<!-- 支持 Markdown；可用变量 {title}, {author}, {url} -->"
)

// Setting names of category "i18n".
const (
	SettingCategoryI18n = "i18n"

	SettingNameI18nLocale   = "i18nLocale"
	SettingNameI18nTimezone = "i18nTimezone"
)

// Setting names of category "feed".
const (
	SettingCategoryFeed = "feed"

	SettingNameFeedOutputMode = "feedOutputMode"
)

// Setting values of category "feed".
const (
	SettingFeedOutputModeValueAbstract = 0
	SettingFeedOutputModeValueFull     = 1
)

// Setting names of category "thirdStatistic".
const (
	SettingCategoryThirdStatistic = "thirdStatistic"

	SettingNameThirdStatisticBaidu = "thirdStatisticBaidu"
)

// Setting names of category "statistic".
const (
	SettingCategoryStatistic = "statistic"

	SettingNameStatisticArticleCount = "statisticArticleCount"
	SettingNameStatisticCommentCount = "statisticCommentCount"
	SettingNameStatisticViewCount    = "statisticViewCount"
)

// Setting names of category "ad".
const (
	SettingCategoryAd = "ad"

	SettingNameAdGoogleAdSenseArticleEmbed = "adGoogleAdSenseArticleEmbed"
)
