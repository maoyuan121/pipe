package model

import (
	"github.com/b3log/pipe/util"
)

// 用户
type User struct {
	Model

	Name              string `gorm:"size:32" json:"name"`       // 用户名
	Nickname          string `gorm:"size:32" json:"nickname"`   // 昵称
	AvatarURL         string `gorm:"size:255" json:"avatarURL"` // 头像 url 地址
	B3Key             string `gorm:"size:32" json:"b3Key"`      // todo
	Locale            string `gorm:"size:32" json:"locale"`     // 哪个语言
	TotalArticleCount int    `json:"totalArticleCount"`         // 总文章数
	GithubId          string `gorm:"255" json:"githubId"`       // todo
}

// 用户角色
const (
	UserRoleNoLogin       = iota // 匿名用户
	UserRolePlatformAdmin        // 平台管理员
	UserRoleBlogAdmin            // 博客管理员
	UserRoleBlogUser             // 博客用户
)

// 返回指定图片大小的 avatar url
func (u *User) AvatarURLWithSize(size int) string {
	return util.ImageSize(u.AvatarURL, size, size)
}
