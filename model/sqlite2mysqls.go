// 把 sqllite 的数据导入到 mysql 中去

package model

import (
	"github.com/jinzhu/gorm"
)

func sqlite2MySQL(sqliteDataFilePath, mysqlConn string) {
	sqlite, err := gorm.Open("sqlite3", Conf.SQLite)
	if nil != err {
		logger.Fatalf("opens SQLite database failed: " + err.Error())
	}
	mysql, err := gorm.Open("mysql", Conf.MySQL)
	if nil != err {
		logger.Fatalf("opens MySQL database failed: " + err.Error())
	}
	if err = mysql.AutoMigrate(Models...).Error; nil != err {
		logger.Fatal("auto migrate tables failed: " + err.Error())
	}

	importArchives(sqlite, mysql, []*Archive{})
	importArticles(sqlite, mysql, []*Article{})
	importCategories(sqlite, mysql, []*Category{})
	importComments(sqlite, mysql, []*Comment{})
	importCorrelations(sqlite, mysql, []*Correlation{})
	importNavigations(sqlite, mysql, []*Navigation{})
	importSettings(sqlite, mysql, []*Setting{})
	importTags(sqlite, mysql, []*Tag{})
	importUsers(sqlite, mysql, []*User{})
}

func importArchives(sqlite, mysql *gorm.DB, models []*Archive) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] archives", len(models))
}

func importArticles(sqlite, mysql *gorm.DB, models []*Article) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if model.PushedAt.Before(ZeroPushTime) {
			model.PushedAt = ZeroPushTime
		}
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: "+err.Error()+": %+v", model)
		}
	}
	logger.Infof("imported [%d] articles", len(models))
}

func importCategories(sqlite, mysql *gorm.DB, models []*Category) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] categories", len(models))
}

func importComments(sqlite, mysql *gorm.DB, models []*Comment) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] comments", len(models))
}

func importCorrelations(sqlite, mysql *gorm.DB, models []*Correlation) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] correlations", len(models))
}

func importNavigations(sqlite, mysql *gorm.DB, models []*Navigation) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] navigations", len(models))
}

func importSettings(sqlite, mysql *gorm.DB, models []*Setting) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] settings", len(models))
}

func importTags(sqlite, mysql *gorm.DB, models []*Tag) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] tags", len(models))
}

func importUsers(sqlite, mysql *gorm.DB, models []*User) {
	if err := sqlite.Find(&models).Error; nil != err {
		logger.Fatalf("queries data failed: " + err.Error())
	}
	for _, model := range models {
		if err := mysql.Save(model).Error; nil != err {
			logger.Fatalf("saves data failed: " + err.Error())
		}
	}
	logger.Infof("imported [%d] users", len(models))
}
