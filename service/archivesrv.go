package service

import (
	"sync"

	"github.com/b3log/pipe/model"
	"github.com/jinzhu/gorm"
)

// 实例化一个归档服务对象
var Archive = &archiveService{
	mutex: &sync.Mutex{},
}

// 定义归档服务
type archiveService struct {
	mutex *sync.Mutex
}

// 获取一个博客的所有归档信息
func (srv *archiveService) GetArchives(blogID uint64) []*model.Archive {
	var ret []*model.Archive
	if err := db.Where("`blog_id` = ? AND `article_count` > 0", blogID).Order("`year` DESC, `month` DESC").Find(&ret).Error; nil != err {
		logger.Error("get archives failed: " + err.Error())
	}

	return ret
}

func (srv *archiveService) UnArchiveArticleWithoutTx(tx *gorm.DB, article *model.Article) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	year := article.CreatedAt.Format("2006")
	month := article.CreatedAt.Format("01")
	archive := &model.Archive{Year: year, Month: month, BlogID: article.BlogID}
	if err := tx.Where("`year` = ? AND `month` = ? AND `blog_id` = ?",
		year, month, article.BlogID).First(archive).Error; nil != err {
		return err
	}
	archive.ArticleCount--
	if archive.ArticleCount < 0 {
		logger.Error("impossible: article count < 0")
		archive.ArticleCount = 0
	}
	if err := tx.Save(archive).Error; nil != err {
		return err
	}
	if err := tx.Where("`id1` = ? AND `id2` = ? AND `type` = ? AND `blog_id` = ?",
		article.ID, archive.ID, model.CorrelationArticleArchive, article.BlogID).
		Delete(&model.Correlation{}).Error; nil != err {
		return err
	}

	return nil
}

// 归档这篇文章
// 更新归档统计数据，保存归档和文章的关系
func (srv *archiveService) ArchiveArticleWithoutTx(tx *gorm.DB, article *model.Article) error {
	// 加锁，防并发引起的问题
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	year := article.CreatedAt.Format("2006")
	month := article.CreatedAt.Format("01")

	archive := &model.Archive{Year: year, Month: month, BlogID: article.BlogID}
	if err := tx.Where("`year` = ? AND `month` = ? AND `blog_id` = ?",
		year, month, article.BlogID).First(archive).Error; nil != err {
		if gorm.ErrRecordNotFound != err {
			return err
		}
	}
	archive.ArticleCount++
	// save == upsert
	if err := tx.Save(archive).Error; nil != err {
		return err
	}

	articleArchiveRel := &model.Correlation{
		ID1:    article.ID,
		ID2:    archive.ID,
		Type:   model.CorrelationArticleArchive,
		BlogID: article.BlogID,
	}
	if err := tx.Create(articleArchiveRel).Error; nil != err {
		return err
	}

	return nil
}

// 获取一个博客的指定年月的归档信息
func (srv *archiveService) GetArchive(year, month string, blogID uint64) *model.Archive {
	ret := &model.Archive{}
	if err := db.Where("`year` = ? AND `month` = ? AND `blog_id` = ?",
		year, month, blogID).First(ret).Error; nil != err {
		return nil
	}

	return ret
}
