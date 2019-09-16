package cache

import (
	"github.com/b3log/pipe/model"
	"github.com/bluele/gcache"
)

// 实例化一个评论缓存服务对象
var Comment = &commentCache{
	idHolder: gcache.New(1024 * 10 * 10).LRU().Build(),
}

// 定义评论缓存服务
type commentCache struct {
	idHolder gcache.Cache
}

// 设置评论缓存
func (cache *commentCache) Put(comment *model.Comment) {
	if err := cache.idHolder.Set(comment.ID, comment); nil != err {
		logger.Errorf("put comment [id=%d] into cache failed: %s", comment.ID, err)
	}
}

// 获取评论缓存
func (cache *commentCache) Get(id uint) *model.Comment {
	ret, err := cache.idHolder.Get(id)
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get comment [id=%d] from cache failed: %s", id, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.Comment)
}
