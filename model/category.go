package model

// 分类
type Category struct {
	Model

	Title           string `gorm:"size:128" json:"title"`            // 分类名
	Path            string `gorm:"size:255" json:"path"`             // todo
	Description     string `gorm:"size:255" json:"description"`      // 描述
	MetaKeywords    string `gorm:"size:255" json:"metaKeywords"`     // meta keyword
	MetaDescription string `gorm:"type:text" json:"metaDescription"` // meta description
	Tags            string `gorm:"type:text" json:"tags"`            // 标签
	Number          int    `json:"number"`                           // 排序号
	BlogID          uint64 `sql:"index" json:"blogID"`               // 所属博客 ID
}
