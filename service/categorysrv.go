package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/util"
	"github.com/jinzhu/gorm"
)

// 实例化一个分类服务对象
var Category = &categoryService{
	mutex: &sync.Mutex{},
}

// 定义分类服务
type categoryService struct {
	mutex *sync.Mutex
}

// 后台的分类分页信息
const (
	adminConsoleCategoryListPageSize   = 15 // 管理界面一页显示多少分类
	adminConsoleCategoryListWindowSize = 20 // 管理界面最多显示多少个页码按钮
)

// 根据 path 在指定博客下获取分类
func (srv *categoryService) GetCategoryByPath(path string, blogID uint64) *model.Category {
	path = strings.TrimSpace(path)
	if "" == path || util.IsReservedPath(path) {
		return nil
	}
	path, _ = url.PathUnescape(path)

	ret := &model.Category{}
	if err := db.Where("`path` = ? AND `blog_ID` = ?", path, blogID).First(ret).Error; nil != err {
		return nil
	}

	return ret
}

// 更新分类
// 更新分类和标签的关系
func (srv *categoryService) UpdateCategory(category *model.Category) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	count := 0
	if db.Model(&model.Category{}).Where("`id` = ? AND `blog_id` = ?", category.ID, category.BlogID).
		Count(&count); 1 > count {
		return fmt.Errorf("not found category [id=%d] to update", category.ID)
	}

	tagStr := normalizeTagStr(category.Tags)
	category.Tags = tagStr

	if err := normalizeCategoryPath(category); nil != err {
		return err
	}

	tx := db.Begin()
	if err := tx.Model(category).Updates(category).Error; nil != err {
		tx.Rollback()

		return err
	}
	if err := tx.Where("`id1` = ? AND `type` = ? AND `blog_id` = ?",
		category.ID, model.CorrelationCategoryTag, category.BlogID).Delete(model.Correlation{}).Error; nil != err {
		tx.Rollback()

		return err
	}
	if err := tagCategory(tx, category); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

// 获取分类 for console
func (srv *categoryService) ConsoleGetCategory(id uint64) *model.Category {
	ret := &model.Category{}
	if err := db.First(ret, id).Error; nil != err {
		return nil
	}

	return ret
}

// 新建分类
// 新增分类和标签的关系
func (srv *categoryService) AddCategory(category *model.Category) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	tagStr := normalizeTagStr(category.Tags)
	category.Tags = tagStr

	if err := normalizeCategoryPath(category); nil != err {
		return err
	}

	tx := db.Begin()
	if err := tx.Create(category).Error; nil != err {
		tx.Rollback()

		return err
	}
	if err := tagCategory(tx, category); nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

// 通过 tag 获取分类
func (srv *categoryService) GetCategoriesByTag(tagTitle string, blogID uint64) (ret []*model.Category) {
	if err := db.Where("`blog_id` = ? AND `tags` LIKE ?", blogID, tagTitle).Find(&ret).Error; nil != err {
		logger.Errorf("get categories failed: " + err.Error())
	}

	return
}

// 获取指定博客下的分类列列表 for console
func (srv *categoryService) ConsoleGetCategories(page int, blogID uint64) (ret []*model.Category, pagination *util.Pagination) {
	offset := (page - 1) * adminConsoleCategoryListPageSize
	count := 0
	if err := db.Model(&model.Category{}).Order("`number` ASC, `id` DESC").
		Where("`blog_id` = ?", blogID).
		Count(&count).Offset(offset).Limit(adminConsoleCategoryListPageSize).Find(&ret).Error; nil != err {
		logger.Errorf("get categories failed: " + err.Error())
	}

	pagination = util.NewPagination(page, adminConsoleCategoryListPageSize, adminConsoleCategoryListWindowSize, count)

	return
}

// 获取指定博客下的 N 个分类（按 number 升序）
func (srv *categoryService) GetCategories(size int, blogID uint64) (ret []*model.Category) {
	if err := db.Where("`blog_id` = ?", blogID).Order("`number` asc").Limit(size).Find(&ret).Error; nil != err {
		logger.Errorf("get categories failed: " + err.Error())
	}

	return
}

// 删除分类
// 删除分类和标签的关系
func (srv *categoryService) RemoveCategory(id, blogID uint64) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	category := &model.Category{}

	tx := db.Begin()
	if err := tx.Where("`id` = ? AND `blog_id` = ?", id, blogID).Find(category).Error; nil != err {
		return err
	}

	if err := tx.Where("`id1` = ? AND `type` = ? AND `blog_id` = ?",
		category.ID, model.CorrelationCategoryTag, category.BlogID).Delete(model.Correlation{}).Error; nil != err {
		tx.Rollback()

		return err
	}
	if err := tx.Delete(category).Error; nil != err {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil
}

func normalizeCategoryPath(category *model.Category) error {
	path := strings.TrimSpace(category.Path)
	if "" == path {
		path = "/" + category.Title
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	count := 0
	if db.Model(&model.Category{}).Where("`path` = ? AND `id` != ? AND `blog_id` = ?", path, category.ID, category.BlogID).Count(&count); 0 < count {
		return errors.New("path is reduplicated")
	}

	category.Path = path

	return nil
}

// 新建分类和标签的关系
func tagCategory(tx *gorm.DB, category *model.Category) error {
	tags := strings.Split(category.Tags, ",")
	for _, tagTitle := range tags {
		tag := &model.Tag{BlogID: category.BlogID}
		tx.Where("`title` = ? AND `blog_id` = ?", tagTitle, category.BlogID).First(tag)
		if "" == tag.Title {
			tag.Title = tagTitle
			if err := tx.Create(tag).Error; nil != err {
				return err
			}
		}

		rel := &model.Correlation{
			ID1:    category.ID,
			ID2:    tag.ID,
			Type:   model.CorrelationCategoryTag,
			BlogID: category.BlogID,
		}
		if err := tx.Create(rel).Error; nil != err {
			return err
		}
	}

	return nil
}
