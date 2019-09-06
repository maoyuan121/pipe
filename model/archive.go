// Package model is the "model" layer which defines entity structures with ORM and controller.
package model

// 归档信息
type Archive struct {
	Model

	Year         string `gorm:"size:4" json:"year"`  // 年
	Month        string `gorm:"size:2" json:"month"` // 月
	ArticleCount int    `json:"articleCount"`        // 文章数
	BlogID       uint64 `sql:"index" json:"blogID"`  // 博客 ID
}
