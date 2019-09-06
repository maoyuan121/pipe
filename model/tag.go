package model

// 标签
type Tag struct {
	Model

	Title        string `gorm:"size:128" json:"title"` // 标签名
	ArticleCount int    `json:"articleCount"` // 文章数
	BlogID uint64 `sql:"index" json:"blogID"` // 所属博客 ID
}
