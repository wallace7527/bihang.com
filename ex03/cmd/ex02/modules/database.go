package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)

//var defaultDatabase *sql.DB

func init() {
	////defaultDatabase, _ = sql.Open("mysql", "ebkadmin:ebkadmin@tcp(192.168.1.77:3306)/ebk?charset=utf8")
	//defaultDatabase, _ = sql.Open("mysql", "wanglei:123123@tcp(192.168.1.200:3306)/east_database?charset=utf8")
	//defaultDatabase.SetMaxOpenConns(10)
	//defaultDatabase.SetMaxIdleConns(1)
	//defaultDatabase.Ping()

	gormDB, err := gorm.Open("mysql", "wanglei:123123@tcp(192.168.1.200:3306)/east_database?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	//defer gormDB.Close()

	gormDB.DB().SetMaxOpenConns(100)
	gormDB.DB().SetMaxIdleConns(10)
	gormDB.DB().Ping()
}
