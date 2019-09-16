// 包含主题相关的操作
package theme

import (
	"os"

	"github.com/b3log/gulu"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// 默认的主题
const DefaultTheme = "Littlewin"

// 所有的主题名
var Themes []string

// 扫描 theme/x 目录，将所有的主题名赋给 Themes
func Load() {
	f, _ := os.Open("theme/x")
	names, _ := f.Readdirnames(-1)
	f.Close()

	for _, name := range names {
		// 如果第一个字符不是字母 continue
		if !gulu.Rune.IsNumOrLetter(rune(name[0])) {
			continue
		}

		Themes = append(Themes, name)
	}

	logger.Debugf("loaded [%d] themes", len(Themes))
}
