package db

import (
	"github.com/tanwenzan/blog-api/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync/atomic"
)

var dataSource *gorm.DB

// 初始化标志，0为未初始化，1为已经初始化
var initialized uint32 = 0

func InitDatasource() error {
	if !atomic.CompareAndSwapUint32(&initialized, 0, 1) {
		// 代表初始化完成
		return nil
	}
	dsn := "user:password@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:" + err.Error())
		return err
	}
	log.Print("数据库连接成功:ip=localhost,port=3306,db=blog")
	// 自动迁移
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err == nil {
		dataSource = db
	} else {
		log.Fatal("数据库自动迁移失败:" + err.Error())
	}
	return err
}

func GetDatasource() *gorm.DB {
	if initialized == 0 {
		_ = InitDatasource()
	}
	return dataSource
}
