package service

import (
	"os"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"  // mysql
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// 数据库
var db *gorm.DB

// 是不是用的 SQLite
var useSQLite bool

// 连接数据库、迁移数据库、创建索引、相关数据库的设置
func ConnectDB() {
	var err error
	useSQLite = false
	if "" != model.Conf.SQLite {
		db, err = gorm.Open("sqlite3", model.Conf.SQLite)
		useSQLite = true
	} else if "" != model.Conf.MySQL {
		db, err = gorm.Open("mysql", model.Conf.MySQL)
	} else {
		logger.Fatal("please specify database")
	}
	if nil != err {
		logger.Fatalf("opens database failed: " + err.Error())
	}
	if useSQLite {
		logger.Debug("used [SQLite] as underlying database")
	} else {
		logger.Debug("used [MySQL] as underlying database")
	}

	// 迁移数据库
	if err = db.AutoMigrate(model.Models...).Error; nil != err {
		logger.Fatal("auto migrate tables failed: " + err.Error())
	}

	// 添加索引
	// 为文章的 created_at 创建索引
	if err = db.Model(&model.Article{}).AddIndex("idx_b3_pipe_articles_created_at", "created_at").Error; nil != err {
		logger.Fatal("adds index failed: " + err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.LogMode(model.Conf.ShowSQL)
}

// 关闭数据库连接
func DisconnectDB() {
	if err := db.Close(); nil != err {
		logger.Errorf("Disconnect from database failed: " + err.Error())
	}
}

// 返回使用的是什么数据库 SQLite 或者 MySQL
func Database() string {
	if useSQLite {
		return "SQLite"
	}

	return "MySQL"
}
