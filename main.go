package main

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/controller"
	"github.com/b3log/pipe/cron"
	"github.com/b3log/pipe/i18n"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/theme"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
)

// Logger
var logger *gulu.Logger

// The only one init function in pipe.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	gulu.Log.SetLevel("warn")
	logger = gulu.Log.NewLogger(os.Stdout)

	// 加载配置，命令行的参数会覆盖配置文件的配置
	model.LoadConf()
	util.LoadMarkdown()

	// 加载多语言信息
	i18n.Load()
	// 加载所有的主题
	theme.Load()
	replaceServerConf()

	if "dev" == model.Conf.RuntimeMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
}

// Entry point.
func main() {
	// 连接数据库
	service.ConnectDB()
	service.Upgrade.Perform()
	// 开始所有的定时任务
	cron.Start()

	// 返回一个 gin engine， 将请求 URL 绑定到控制器
	router := controller.MapRoutes()
	server := &http.Server{
		Addr:    "0.0.0.0:" + model.Conf.Port,
		Handler: router,
	}

	handleSignal(server)

	logger.Infof("Pipe (v%s) is running [%s]", model.Version, model.Conf.Server)
	if err := server.ListenAndServe(); nil != err {
		logger.Fatalf("listen and serve failed: " + err.Error())
	}
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Infof("got signal [%s], exiting pipe now", s)
		if err := server.Close(); nil != err {
			logger.Errorf("server close failed: " + err.Error())
		}

		service.DisconnectDB()

		logger.Infof("Pipe exited")
		os.Exit(0)
	}()
}

func replaceServerConf() {
	path := "theme/sw.min.js.tpl"
	if gulu.File.IsExist(path) {
		data, err := ioutil.ReadFile(path)
		if nil != err {
			logger.Fatal("read file [" + path + "] failed: " + err.Error())
		}
		content := string(data)
		content = strings.Replace(content, "http://server.tpl.json", model.Conf.Server, -1)
		content = strings.Replace(content, "http://staticserver.tpl.json", model.Conf.StaticServer, -1)
		content = strings.Replace(content, "${StaticResourceVersion}", model.Conf.StaticResourceVersion, -1)
		writePath := strings.TrimSuffix(path, ".tpl")
		if err = ioutil.WriteFile(writePath, []byte(content), 0644); nil != err {
			logger.Fatal("replace sw.min.js in [" + path + "] failed: " + err.Error())
		}
	}

	if gulu.File.IsExist("console/dist/") {
		err := filepath.Walk("console/dist/", func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".tpl") {
				data, err := ioutil.ReadFile(path)
				if nil != err {
					logger.Fatal("read file [" + path + "] failed: " + err.Error())
				}
				content := string(data)
				content = strings.Replace(content, "http://server.tpl.json", model.Conf.Server, -1)
				content = strings.Replace(content, "http://staticserver.tpl.json", model.Conf.StaticServer, -1)
				writePath := strings.TrimSuffix(path, ".tpl")
				if err = ioutil.WriteFile(writePath, []byte(content), 0644); nil != err {
					logger.Fatal("replace server conf in [" + writePath + "] failed: " + err.Error())
				}
			}

			return err
		})
		if nil != err {
			logger.Fatal("replace server conf in [theme] failed: " + err.Error())
		}
	}
}
