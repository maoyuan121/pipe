package cron

import (
	"os"

	"github.com/b3log/gulu"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// 开始所有的定时任务
func Start() {
	//  每 30 分钟执行一次， 取阅读数最多的 7 篇文章作为推荐文件
	refreshRecommendArticlesPeriodically()

	// 每 30 分钟执行一次，把没有 push 到 b3log 的文章 push 到 b3log 上
	pushArticlesPeriodically()

	// 每 30 分钟执行一次，把没有 push 到 b3log 的评论 push 到 b3log 上
	pushCommentsPeriodically()
}
