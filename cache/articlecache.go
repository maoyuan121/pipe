package cache

import (
	"os"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/bluele/gcache"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// 实例化一个文章缓存服务对象
var Article = &articleCache{
	idHolder: gcache.New(1024 * 10).LRU().Build(),
}

// 定义文章缓存服务
type articleCache struct {
	idHolder gcache.Cache
}

// 设置文章缓存
func (cache *articleCache) Put(article *model.Article) {
	if err := cache.idHolder.Set(article.ID, article); nil != err {
		logger.Errorf("put article [id=%d] into cache failed: %s", article.ID, err)
	}
}

// 获取文章缓存
func (cache *articleCache) Get(id uint) *model.Article {
	ret, err := cache.idHolder.Get(id)
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get article [id=%d] from cache failed: %s", id, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.Article)
}
