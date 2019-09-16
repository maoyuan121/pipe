package cache

import (
	"github.com/b3log/pipe/model"
	"github.com/bluele/gcache"
)

// 实例化一个用户缓存服务对象
var User = &userCache{
	idHolder: gcache.New(1024 * 10).LRU().Build(),
}

// 定义用户缓存服务
type userCache struct {
	idHolder gcache.Cache
}

// 设置用户缓存
func (cache *userCache) Put(user *model.User) {
	if err := cache.idHolder.Set(user.ID, user); nil != err {
		logger.Errorf("put user [id=%d] into cache failed: %s", user.ID, err)
	}
}

// 获取用户缓存
func (cache *userCache) Get(id uint64) *model.User {
	ret, err := cache.idHolder.Get(id)
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get user [id=%d] from cache failed: %s", id, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.User)
}
