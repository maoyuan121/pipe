package model

// Correlation types.
const (
	CorrelationCategoryTag = iota
	CorrelationArticleTag
	CorrelationBlogUser
	CorrelationArticleArchive // 文章 ID，归档 ID，博客 ID
)

// 关联表
// id1(category_id) - id2(tag_id) 分类和标签的关联
// id1(article_id) - id2(tag_id) 文章和标签的关联
// id1(blog_id) - id2(user_id) - int1(role) - int2(article_count) 博客-用户-角色-文章数量的关联
// id1(article_id) - id2(archive_id) 文章和归档的关联
type Correlation struct {
	Model

	ID1  uint64 `json:"id1"`
	ID2  uint64 `json:"id2"`
	Str1 string `gorm:"size:255" json:"str1"`
	Str2 string `gorm:"size:255" json:"str2"`
	Str3 string `gorm:"size:255" json:"str3"`
	Str4 string `gorm:"size:255" json:"str4"`
	Int1 int    `json:"int1"`
	Int2 int    `json:"int2"`
	Int3 int    `json:"int3"`
	Int4 int    `json:"int4"`
	Type int    `json:"type"`
	BlogID uint64 `sql:"index" json:"blogID"`
}
