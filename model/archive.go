package model

// 归档信息
// 用于统计一个博客某年某月发表了多少篇文章
type Archive struct {
	Model

	Year         string `gorm:"size:4" json:"year"`  // 年
	Month        string `gorm:"size:2" json:"month"` // 月
	ArticleCount int    `json:"articleCount"`        // 文章数
	BlogID       uint64 `sql:"index" json:"blogID"`  // 博客 ID
}
