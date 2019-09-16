package console

import (
	"html/template"

	"github.com/b3log/pipe/util"
)

// 文章
type ConsoleArticle struct {
	ID           uint64         `json:"id"`           // 主键
	Author       *ConsoleAuthor `json:"author"`       // 作者
	CreatedAt    string         `json:"createdAt"`    // 创建时间
	Title        string         `json:"title"`        // 标题
	Tags         []*ConsoleTag  `json:"tags"`         // 标签
	URL          string         `json:"url"`          // url
	Topped       bool           `json:"topped"`       // 是否指定
	ViewCount    int            `json:"viewCount"`    // 浏览数
	CommentCount int            `json:"commentCount"` // 评论数
}

// 标签
type ConsoleTag struct {
	ID    uint64 `json:"id"`            // 主键
	Title string `json:"title"`         // 标签名
	URL   string `json:"url,omitempty"` // url
}

// 作者
type ConsoleAuthor struct {
	URL       string `json:"url"`       // url
	Name      string `json:"name"`      // 名字
	AvatarURL string `json:"avatarURL"` // 头像 url
}

// 返回指定大小的头像 url
func (u *ConsoleAuthor) AvatarURLWithSize(size int) string {
	return util.ImageSize(u.AvatarURL, size, size)
}

// 分类
type ConsoleCategory struct {
	ID          uint64 `json:"id"`          // 主键
	Title       string `json:"title"`       // 分类名
	URL         string `json:"url"`         // url
	Description string `json:"description"` // 描述
	Number      int    `json:"number"`      // 排序号
	Tags        string `json:"tags"`        // 标签
}

// 评论
type ConsoleComment struct {
	ID            uint64         `json:"id"`            // 主键
	Author        *ConsoleAuthor `json:"author"`        // 作者
	ArticleAuthor *ConsoleAuthor `json:"articleAuthor"` // 文章作者
	CreatedAt     string         `json:"createdAt"`     // 创建时间
	Title         string         `json:"title"`         // 标题
	Content       template.HTML  `json:"content"`       // 内容
	URL           string         `json:"url"`           // url
}

// 导航
type ConsoleNavigation struct {
	ID         uint64 `json:"id"`         // 主键
	Title      string `json:"title"`      // 导航名
	URL        string `json:"url"`        // url
	IconURL    string `json:"iconURL"`    // 图标 url
	OpenMethod string `json:"openMethod"` // 打开方式 （eg: _self、_blank）
	Number     int    `json:"number"`     // 排序号
}

// 主题
type ConsoleTheme struct {
	Name         string `json:"name"`         // 主题名
	ThumbnailURL string `json:"thumbnailURL"` // 主题缩略图 url
}

// 用户
type ConsoleUser struct {
	ID           uint64 `json:"id"`           // 主键
	Name         string `json:"name"`         // 用户名
	Nickname     string `json:"nickname"`     // 昵称
	Role         int    `json:"role"`         // 角色
	URL          string `json:"url"`          // url
	AvatarURL    string `json:"avatarURL"`    // 头像 url
	ArticleCount int    `json:"articleCount"` // 发表的文章数
}
