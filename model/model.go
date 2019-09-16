package model

import "time"

// gorm  模型基类
type Model struct {
	ID        uint64     `gorm:"primary_key" json:"id" structs:"id"`        // 主键
	CreatedAt time.Time  `json:"createdAt" structs:"createdAt"`             // 创建时间
	UpdatedAt time.Time  `json:"updatedAt" structs:"updatedAt"`             // 最后更新时间
	DeletedAt *time.Time `sql:"index" json:"deletedAt" structs:"deletedAt"` // 删除时间（可为空）
}
