package model

import (
	"time"
)

// Article model.
type Article struct {
	Model

	AuthorID     uint64    `json:"authorID" structs:"authorID"`                        // 作者 ID
	Title        string    `gorm:"size:128" json:"title" structs:"title"`              // 标题
	Abstract     string    `gorm:"type:mediumtext" json:"abstract" structs:"abstract"` // 摘要
	Tags         string    `gorm:"type:text" json:"tags" structs:"tags"`               // 标签
	Content      string    `gorm:"type:mediumtext" json:"content" structs:"content"`   // 内容
	Path         string    `sql:"index" gorm:"size:255" json:"path" structs:"path"`    // todo
	Status       int       `sql:"index" json:"status" structs:"status"`                // 状态
	Topped       bool      `json:"topped" structs:"topped"`                            // 是否置顶
	Commentable  bool      `json:"commentable" structs:"commentable"`                  // 是否可评论
	ViewCount    int       `json:"viewCount" structs:"viewCount"`                      // 阅读数
	CommentCount int       `json:"commentCount" structs:"commentCount"`                // 评论数
	IP           string    `gorm:"size:128" json:"ip" structs:"ip"`                    // IP
	UserAgent    string    `gorm:"size:255" json:"userAgent" structs:"userAgent"`      // user agent
	PushedAt     time.Time `json:"pushedAt" structs:"pushedAt"`                        // todo
	BlogID       uint64    `sql:"index" json:"blogID" structs:"blogID"`                // 所属博客 ID
}

// 文章状态
const (
	ArticleStatusOK = iota
)
