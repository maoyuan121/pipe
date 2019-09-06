package model

import (
	"math"
	"time"
)

// 评论
type Comment struct {
	Model

	ArticleID       uint64    `json:"articleID"`                       // 所属文章 ID
	AuthorID        uint64    `json:"authorID"`                        // 作者 ID
	Content         string    `gorm:"type:text" json:"content"`        // 评论内容
	ParentCommentID uint64    `json:"parentCommentID"`                 // 父评论 ID
	IP              string    `gorm:"size:128" json:"ip"`              // 发表于哪个 IP 地址
	UserAgent       string    `gorm:"size:255" json:"userAgent"`       // 发表者的浏览器信息
	PushedAt        time.Time `json:"pushedAt"`                        // 评论时间
	AuthorName      string    `gorm:"size:32" json:"authorName"`       // 作者名 exist if this comment sync from Sym, https://github.com/b3log/pipe/issues/98
	AuthorAvatarURL string    `gorm:"size:255" json:"authorAvatarURL"` // 作者头像 URL exist if this comment sync from Sym, https://github.com/b3log/pipe/issues/98
	AuthorURL       string    `gorm:"size:255" json:"authorURL"`       // 作者 URL exist if this comment sync from Sym, https://github.com/b3log/pipe/issues/98
	BlogID          uint64    `sql:"index" json:"blogID"`              // 所属博客 ID
}

// SyncCommentAuthorID is the id of sync comment bot.
const SyncCommentAuthorID = math.MaxInt32
