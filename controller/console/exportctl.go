package console

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// 导出文章为 md 的压缩包
func ExportMarkdownAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	if 0 == session.UID {
		result.Code = util.CodeErr
		result.Msg = "please login before export"

		return
	}

	tempDir := os.TempDir()
	logger.Trace("temp dir path is [" + tempDir + "]")
	zipFilePath := filepath.Join(tempDir, session.UName+"-export-md.zip")
	zipFile, err := gulu.Zip.Create(zipFilePath)
	if nil != err {
		logger.Errorf("create zip file [" + zipFilePath + "] failed: " + err.Error())
		result.Code = util.CodeErr
		result.Msg = "create zip file failed"

		return
	}

	c.Header("Content-Disposition", "attachment; filename="+session.UName+"-export-md.zip")
	c.Header("Content-Type", "application/zip")

	mdFiles := service.Export.ExportMarkdowns(session.BID)
	if 1 > len(mdFiles) {
		zipFile.Close()
		file, err := os.Open(zipFilePath)
		if nil != err {
			logger.Errorf("open zip file [" + zipFilePath + " failed: " + err.Error())
			result.Code = util.CodeErr
			result.Msg = "open zip file failed"

			return
		}
		defer file.Close()

		io.Copy(c.Writer, file)

		return
	}

	zipPath := filepath.Join(tempDir, session.UName+"-export-md")
	if err = os.RemoveAll(zipPath); nil != err {
		logger.Errorf("remove temp dir [" + zipPath + "] failed: " + err.Error())
		result.Code = util.CodeErr
		result.Msg = "remove temp dir failed"

		return
	}
	if err = os.Mkdir(zipPath, 0755); nil != err {
		logger.Errorf("make temp dir [" + zipPath + "] failed: " + err.Error())
		result.Code = util.CodeErr
		result.Msg = "make temp dir failed"

		return
	}
	for _, mdFile := range mdFiles {
		filename := filepath.Join(zipPath, mdFile.Name+".md")
		if err := ioutil.WriteFile(filename, []byte(mdFile.Content), 0644); nil != err {
			logger.Errorf("write file [" + filename + "] failed: " + err.Error())
		}
	}

	zipFile.AddDirectory(session.UName+"-export-md", zipPath)
	if err := zipFile.Close(); nil != err {
		logger.Errorf("zip failed: " + err.Error())
		result.Code = util.CodeErr
		result.Msg = "zip failed"

		return
	}
	file, err := os.Open(zipFilePath)
	if nil != err {
		logger.Errorf("open zip file [" + zipFilePath + " failed: " + err.Error())
		result.Code = util.CodeErr
		result.Msg = "open zip file failed"

		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}
