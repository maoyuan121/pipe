package service

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/util"
	"gopkg.in/yaml.v2"
)

// Import service.
var Import = &importService{}

type importService struct {
}

// MarkdownFile represents markdown file.
type MarkdownFile struct {
	Name    string
	Path    string
	Content string
}

type importArticles []*model.Article

func (a importArticles) Len() int {
	return len(a)
}
func (a importArticles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a importArticles) Less(i, j int) bool {
	return a[j].UpdatedAt.After(a[i].UpdatedAt)
}

func (srv *importService) ImportMarkdowns(mdFiles []*MarkdownFile, authorID, blogID uint64) {
	succCnt, failCnt := 0, 0
	var fails []string

	var articles importArticles
	for _, mdFile := range mdFiles {
		article := parseArticle(mdFile)
		article.AuthorID = authorID
		article.BlogID = blogID

		if strings.HasPrefix(article.Path, util.PathArticles) && len(util.PathArticles+"/") < len(article.Path) {
			article.Path = ""
		}

		articles = append(articles, article)
	}

	sort.Sort(articles)

	for _, article := range articles {
		if err := Article.AddArticle(article); nil != err {
			failCnt++
			fails = append(fails, article.Title)
			logger.Errorf("import article [" + article.Title + "] failed: " + err.Error())

			continue
		}

		succCnt++
	}

	if 0 == succCnt && 0 == failCnt {
		return
	}

	logBuilder := "[" + strconv.Itoa(succCnt) + "] imported, [" + strconv.Itoa(failCnt) + "] failed"
	if 0 < failCnt {
		logBuilder += ": \n"
		for _, fail := range fails {
			logBuilder += "    " + fail + "\n"
		}
	} else {
		logBuilder += " :p"
	}

	logger.Info(logBuilder)
}

func parseArticle(mdFile *MarkdownFile) *model.Article {
	gulu.Panic.Recover()

	content := strings.TrimSpace(mdFile.Content)
	frontMatter := strings.Split(content, "---")[0]
	if "" == frontMatter {
		content = strings.Split(content, "---")[1]
		frontMatter = strings.Split(content, "---")[0]
	}

	ret := &model.Article{}

	m := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(frontMatter), &m)
	if nil != err || 0 == len(m) {
		ext := filepath.Ext(mdFile.Name)
		ret.Title = strings.Split(mdFile.Name, ext)[0]
		ret.Content = mdFile.Content
		ret.Commentable = true
		ret.Tags = "笔记"

		return ret
	}

	ext := filepath.Ext(mdFile.Name)
	title := strings.Split(mdFile.Name, ext)[0]
	if t, ok := m["title"]; ok {
		title = strings.TrimSpace(t.(string))
	}
	ret.Title = title

	content = strings.TrimSpace(strings.Split(mdFile.Content, frontMatter)[1])
	if strings.HasPrefix(content, "---") {
		content = content[len("---"):]
		content = strings.TrimSpace(content)
	}
	ret.Content = content

	permalink := ""
	if p, ok := m["permalink"]; ok {
		permalink = strings.TrimSpace(p.(string))
	}
	ret.Path = permalink

	tags := parseTags(&m)
	ret.Tags = tags
	ret.CreatedAt = parseDate(&m)
	ret.UpdatedAt = ret.CreatedAt
	ret.PushedAt = ret.CreatedAt
	ret.Commentable = true

	return ret
}

func parseDate(m *map[string]interface{}) time.Time {
	frontMatter := *m
	date := frontMatter["date"]
	if nil == date {
		return time.Now()
	}
	dateStr := strings.TrimSpace(date.(string))
	if "" == dateStr {
		return time.Now()
	}

	ret, err := dateparse.ParseAny(dateStr)
	if nil != err {
		logger.Warn(err.Error())

		return time.Now()
	}

	return ret
}

func parseTags(m *map[string]interface{}) string {
	frontMatter := *m
	tags := frontMatter["tags"]
	if nil == tags {
		tags = frontMatter["category"]
	}
	if nil == tags {
		tags = frontMatter["categories"]
	}
	if nil == tags {
		tags = frontMatter["keyword"]
	}
	if nil == tags {
		tags = frontMatter["keywords"]
	}
	if nil == tags {
		return "笔记"
	}

	switch tags.(type) {
	case []interface{}:
		ts := tags.([]interface{})
		var tagStrs []string
		for _, t := range ts {
			tagStrs = append(tagStrs, t.(string))
		}

		return strings.Join(tagStrs, ",")
	case string:
		return tags.(string)
	default:
		logger.Warnf("unknown type of tags in front matter [%+v]", frontMatter)

		return "笔记"
	}
}
