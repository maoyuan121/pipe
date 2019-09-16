package console

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// ShowAdminPagesAction shows admin pages.
func ShowAdminPagesAction(c *gin.Context) {
	session := util.GetSession(c)
	if 0 == session.UID {
		c.Redirect(http.StatusSeeOther, model.Conf.Server+"/start")

		return
	}

	t, err := template.ParseFiles(filepath.Join("console/dist/admin" + c.Param("path") + "/index.html"))
	if nil != err {
		logger.Errorf("load console page [" + c.Param("path") + "] failed: " + err.Error())
		c.String(http.StatusNotFound, "load console page failed")

		return
	}

	t.Execute(c.Writer, nil)
}
