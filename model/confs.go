// Package util defines variety of utilities.
package model

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/b3log/gulu"
	"github.com/jinzhu/gorm"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// Version of Pipe.
const Version = "1.8.9"

// Conf of Pipe.
var Conf *Configuration

// HTTP Client 的 UserAgent
var UserAgent = "Pipe/" + Version + "; +https://github.com/b3log/pipe"

// Models 表示所有的模型
var Models = []interface{}{
	&User{}, &Article{}, &Comment{}, &Navigation{}, &Tag{},
	&Category{}, &Archive{}, &Setting{}, &Correlation{},
}

// 表前缀
const tablePrefix = "b3_pipe_"

// Push 到 blog 的时间的零值
var ZeroPushTime, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

// 配置信息 (对应 pipe.json).
type Configuration struct {
	Server                string // server scheme, host and port
	StaticServer          string // 静态资源的 host 和 port
	StaticResourceVersion string // version of static resources
	LogLevel              string // logging level: trace/debug/info/warn/error/fatal
	ShowSQL               bool   // 是否在 log 中打印 sql
	SessionSecret         string // HTTP session secret
	SessionMaxAge         int    // HTTP session max age (in second)
	RuntimeMode           string // 运行模式 (dev/prod)
	SQLite                string // SQLite 数据库我呢见地址
	MySQL                 string // MySQL 连接 URL
	Port                  string // 端口
	AxiosBaseURL          string // axio base URL
	MockServer            string // mock 服务器地址
}

// 从加载配置，命令行的参数会覆盖配置文件的配置
func LoadConf() {
	version := flag.Bool("version", false, "prints current pipe version")
	confPath := flag.String("conf", "pipe.json", "path of pipe.json")
	confServer := flag.String("server", "", "this will override Conf.Server if specified")
	confStaticServer := flag.String("static_server", "", "this will override Conf.StaticServer if specified")
	confStaticResourceVer := flag.String("static_resource_ver", "", "this will override Conf.StaticResourceVersion if specified")
	confLogLevel := flag.String("log_level", "", "this will override Conf.LogLevel if specified")
	confShowSQL := flag.Bool("show_sql", false, "this will override Conf.ShowSQL if specified")
	confSessionSecret := flag.String("session_secret", "", "this will override Conf.SessionSecret")
	confSessionMaxAge := flag.Int("session_max_age", 0, "this will override Conf.SessionMaxAge")
	confRuntimeMode := flag.String("runtime_mode", "", "this will override Conf.RuntimeMode if specified")
	confSQLite := flag.String("sqlite", "", "this will override Conf.SQLite if specified")
	confMySQL := flag.String("mysql", "", "this will override Conf.MySQL if specified")
	confPort := flag.String("port", "", "this will override Conf.Port if specified")
	s2m := flag.Bool("s2m", false, "dumps SQLite data to MySQL SQL script file")

	flag.Parse()

	if *version {
		fmt.Println(Version)

		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(*confPath)
	if nil != err {
		logger.Fatal("loads configuration file [" + *confPath + "] failed: " + err.Error())
	}

	Conf = &Configuration{}
	if err = json.Unmarshal(bytes, Conf); nil != err {
		logger.Fatal("parses [pipe.json] failed: ", err)
	}

	gulu.Log.SetLevel(Conf.LogLevel)
	if "" != *confLogLevel {
		Conf.LogLevel = *confLogLevel
		gulu.Log.SetLevel(*confLogLevel)
	}

	if *confShowSQL {
		Conf.ShowSQL = true
	}

	if "" != *confSessionSecret {
		Conf.SessionSecret = *confSessionSecret
	}

	if 0 < *confSessionMaxAge {
		Conf.SessionMaxAge = *confSessionMaxAge
	}

	if "" == Conf.SessionSecret {
		Conf.SessionSecret = gulu.Rand.String(32)
	}

	home, err := gulu.OS.Home()
	if nil != err {
		logger.Fatal("can't find user home directory: " + err.Error())
	}
	logger.Debugf("${home} [%s]", home)

	if "" != *confRuntimeMode {
		Conf.RuntimeMode = *confRuntimeMode
	}

	if "" != *confServer {
		Conf.Server = *confServer
	}

	if "" != *confStaticServer {
		Conf.StaticServer = *confStaticServer
	}
	if "" == Conf.StaticServer {
		Conf.StaticServer = Conf.Server
	}

	time := strconv.FormatInt(time.Now().UnixNano(), 10)
	logger.Debugf("${time} [%s]", time)
	Conf.StaticResourceVersion = strings.Replace(Conf.StaticResourceVersion, "${time}", time, 1)
	if "" != *confStaticResourceVer {
		Conf.StaticResourceVersion = *confStaticResourceVer
	}

	Conf.SQLite = strings.Replace(Conf.SQLite, "${home}", home, 1)
	if "" != *confSQLite {
		Conf.SQLite = *confSQLite
	}
	if "" != *confMySQL {
		Conf.MySQL = *confMySQL
		Conf.SQLite = ""
	}

	if "" != *confPort {
		Conf.Port = *confPort
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	if *s2m {
		if "" == Conf.SQLite {
			logger.Fatal("please specify -sqlite")
		}
		if "" == Conf.MySQL {
			logger.Fatal("please specify -mysql")
		}

		sqlite2MySQL(Conf.SQLite, Conf.MySQL)

		os.Exit(0)
	}

	logger.Debugf("configurations [%#v]", Conf)
}
