package cache

import (
	"fmt"

	"github.com/b3log/pipe/model"
	"github.com/bluele/gcache"
)

// 实例化一个设置缓存服务对象
var Setting = &settingCache{
	categoryNameHolder: gcache.New(1024 * 10).LRU().Build(),
}

// 定义设置缓存服务
type settingCache struct {
	categoryNameHolder gcache.Cache
}

// 设置设置缓存
func (cache *settingCache) Put(setting *model.Setting) {
	if err := cache.categoryNameHolder.Set(fmt.Sprintf("%s-%s-%d", setting.Category, setting.Name, setting.BlogID), setting); nil != err {
		logger.Errorf("put setting [id=%d] into cache failed: %s", setting.ID, err)
	}
}

// 获取设置缓存
func (cache *settingCache) Get(category, name string, blogID uint64) *model.Setting {
	ret, err := cache.categoryNameHolder.Get(fmt.Sprintf("%s-%s-%d", category, name, blogID))
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get setting [name=%s, category=%s, blogID=%d] from cache failed: %s", category, name, blogID, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.Setting)
}
