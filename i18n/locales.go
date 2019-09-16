// 包含多语言和相关的操作
package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/b3log/gulu"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// 本地化结构
type locale struct {
	Name     string                 // 名称
	Langs    map[string]interface{} // 语言包
	TimeZone string                 // 时区
}

// 多语言
// key 为 语种，value 为语言包和时区
// 多语言的数据都存在这个结构中
// 程序都是从这获取语言
var locales = map[string]locale{}

// 加载多语言
// 多语言多外的初始化方法
// 遍历 i18n 文件夹下的文件，将其中的 xxx.json 语言文件的内容加载进来
func Load() {
	f, _ := os.Open("i18n")
	names, _ := f.Readdirnames(-1)
	f.Close()

	for _, name := range names {
		// 如果文件名的第一个字符不是字母，或者文件不是 .json continue
		if !gulu.Rune.IsLetter(rune(name[0])) || !strings.HasSuffix(name, ".json") {
			continue
		}

		// loc 为语种
		loc := name[:strings.LastIndex(name, ".")]
		load(loc)
	}

	logger.Tracef("loaded [%d] language configuration files", len(locales))
}

// 根据指定的语种从对应的 xxx.json 中加载语言到 locales 中去
func load(localeStr string) {
	bytes, err := ioutil.ReadFile("i18n/" + localeStr + ".json")
	if nil != err {
		logger.Fatal("reads i18n configurations fialed: " + err.Error())
	}

	l := locale{Name: localeStr}

	err = json.Unmarshal(bytes, &l.Langs)
	if nil != err {
		logger.Fatal("parses i18n configurations failed: " + err.Error())
	}

	locales[localeStr] = l
}

// 根据给定的语种，获取给定 key 和 args 对应的文本
func GetMessagef(locale, key string, a ...interface{}) string {
	msg := GetMessage(locale, key)

	return fmt.Sprintf(msg, a...)
}

// 根据给定的语种，获取给定 key 对应的文本
func GetMessage(locale, key string) string {
	return locales[locale].Langs[key].(string)
}

// 获取指定语种的语言包
func GetMessages(locale string) map[string]interface{} {
	return locales[locale].Langs
}

// 获取多语言支持的语种
// 例如返回 ["zh_CN", "en_US"]
func GetLocalesNames() []string {
	var ret []string

	for name := range locales {
		ret = append(ret, name)
	}

	sort.Strings(ret)

	return ret
}
