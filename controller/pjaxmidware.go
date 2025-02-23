package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func pjax(c *gin.Context) {
	isPJAX := isPJAX(c)
	dataModelVal, _ := c.Get("dataModel")
	dataModel := dataModelVal.(*DataModel)
	(*dataModel)["pjax"] = isPJAX
	c.Set("dataModel", dataModel)

	if !isPJAX {
		c.Next()

		return
	}

	c.Writer = &pjaxHTMLWriter{c.Writer, &strings.Builder{}, c}
	c.Next()
}

type pjaxHTMLWriter struct {
	gin.ResponseWriter
	bodyBuilder *strings.Builder
	c           *gin.Context
}

func (p *pjaxHTMLWriter) Write(data []byte) (int, error) {
	p.bodyBuilder.Write(data)
	if !strings.HasSuffix(string(data), "</html>\r\n") && !strings.HasSuffix(string(data), "</html>\n") {
		return 0, nil
	}

	pjaxContainer := p.c.Request.Header.Get("X-PJAX-Container")
	body := p.bodyBuilder.String()
	startTag := "<!---- pjax {" + pjaxContainer + "} start ---->"
	endTag := "<!---- pjax {" + pjaxContainer + "} end ---->"
	var containers []string
	count := 0
	part := body
	for {
		start := strings.Index(part, startTag)
		if 0 > start {
			break
		}
		start = start + len(startTag)
		end := strings.Index(part, endTag)
		containers = append(containers, part[start:end])
		count++
		if 10 <= count {
			break
		}
		part = part[end+len(endTag):]
	}

	if 0 != len(containers) {
		body = strings.Join(containers, "")
	}

	i, e := p.ResponseWriter.WriteString(body)
	p.ResponseWriter.Flush()

	return i, e
}

// 判断是不是 pjax 请求
func isPJAX(c *gin.Context) bool {
	return "true" == c.Request.Header.Get("X-PJAX") && "" != c.Request.Header.Get("X-PJAX-Container")
}
