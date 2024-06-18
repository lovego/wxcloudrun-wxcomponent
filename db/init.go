package db

import (
	"fmt"
	"time"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/config"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
)

var dbInstance *gorm.DB
var cacheInstance *cache.Cache

// Init 初始化数据库
func Init() error {
	var user, pwd, addr, dataBase string
	user = config.DBConf.Username
	pwd = config.DBConf.Password
	addr = config.DBConf.Address
	dataBase = config.DBConf.Database
	if dataBase == "" {
		dataBase = "wxcomponent"
	}
	source := "postgres://%s:%s@%s/%s?sslmode=disable"
	source = fmt.Sprintf(source, user, pwd, addr, dataBase)
	log.Debug("start inits postgresql with ::::: " + source)

	db, err := gorm.Open(postgres.Open(source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println("DB Open error,err=", err.Error())
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("DB Init error,err=", err.Error())
		return err
	}

	// 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxIdleConns(100)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(200)
	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbInstance = db

	fmt.Println("finish inits postgresql with ", source)

	checkTables()

	// 初始化cache
	cacheInstance = cache.New(5*time.Minute, 10*time.Minute)

	return nil
}

func checkTables() {
	// 创建表
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS wxcallback_component (id SERIAL PRIMARY KEY,receivetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,infotype VARCHAR(64) NOT NULL DEFAULT '',postbody TEXT NOT NULL);`)
	dbInstance.Exec("CREATE INDEX IF NOT EXISTS wxcallback_component_receivetime_idx ON wxcallback_component(receivetime);")
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS wxcallback_biz (id SERIAL PRIMARY KEY,receivetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,tousername VARCHAR(64) NOT NULL DEFAULT '',appid VARCHAR(64) NOT NULL DEFAULT '',msgtype VARCHAR(64) NOT NULL DEFAULT '',event VARCHAR(64) NOT NULL DEFAULT '',postbody TEXT NOT NULL);`)
	dbInstance.Exec("CREATE INDEX IF NOT EXISTS wxcallback_biz_receivetime_idx ON wxcallback_biz(receivetime);")
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS comm (key VARCHAR(64) NOT NULL PRIMARY KEY,value TEXT NOT NULL,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`)
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS "user" (id SERIAL PRIMARY KEY,username VARCHAR(32) NOT NULL UNIQUE,password VARCHAR(64) NOT NULL,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`)
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS authorizers (id SERIAL PRIMARY KEY,appid VARCHAR(32) NOT NULL UNIQUE,apptype INT NOT NULL DEFAULT 0,servicetype INT NOT NULL DEFAULT 0,nickname VARCHAR(32) NOT NULL DEFAULT '',username VARCHAR(32) NOT NULL DEFAULT '',headimg VARCHAR(256) NOT NULL DEFAULT '',qrcodeurl VARCHAR(256) NOT NULL DEFAULT '',principalname VARCHAR(64) NOT NULL DEFAULT '',refreshtoken VARCHAR(128) NOT NULL DEFAULT '',funcinfo VARCHAR(128) NOT NULL DEFAULT '',verifyinfo INT NOT NULL DEFAULT -1,authtime TIMESTAMP NOT NULL,updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`)
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS wxcallback_rules (id SERIAL PRIMARY KEY,name VARCHAR(64) NOT NULL DEFAULT '',infotype VARCHAR(64) NOT NULL DEFAULT '',msgtype VARCHAR(64) NOT NULL DEFAULT '',event VARCHAR(64) NOT NULL DEFAULT '',type INT NOT NULL DEFAULT 0,open INT NOT NULL DEFAULT 0,info TEXT NOT NULL,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,UNIQUE (infotype, msgtype, event));`)
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS wxtoken (id SERIAL PRIMARY KEY,type INT NOT NULL DEFAULT 0,appid VARCHAR(128) NOT NULL DEFAULT '' UNIQUE,token TEXT NOT NULL,expiretime TIMESTAMP NOT NULL,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`)
	dbInstance.Exec(`CREATE TABLE IF NOT EXISTS counter (id SERIAL PRIMARY KEY,key VARCHAR(64) NOT NULL UNIQUE,value INT NOT NULL,createtime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,updatetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`)

}

// Get
func Get() *gorm.DB {
	return dbInstance
}

// GetCache
func GetCache() *cache.Cache {
	return cacheInstance
}
