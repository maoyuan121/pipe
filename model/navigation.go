package model

// 导航
type Navigation struct {
	Model

	Title      string `gorm:"size:128" json:"title"`     // 标题
	URL        string `gorm:"size:255" json:"url"`       // URL
	IconURL    string `gorm:"size:255" json:"iconURL"`   // 图标 URL
	OpenMethod string `gorm:"size:32" json:"openMethod"` // 点击的打开方式 eg:_blank, _self
	Number     int    `json:"number"`                    // 排序号
	BlogID     uint64 `sql:"index" json:"blogID"`        // 所属博客 ID
}

// Navigation open methods.
const (
	NavigationOpenMethodBlank = "_blank" // 新开窗口
	NavigationOpenMethodSelf  = "_self"  // 就地跳转
)
