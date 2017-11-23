package models

import (
	"fmt"
	"myblog/models/class"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)
func init() {
	orm.Debug = true
	// fmt.Printf(beego.AppConfig.String("DB::db"))
	
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8&loc=%s",
		beego.AppConfig.String("DB::user"),
		beego.AppConfig.String("DB::pass"),
		beego.AppConfig.String("DB::name"),
		`Asia%2FShanghai`,
	))
	orm.RegisterModel(new(class.User), new(class.Article), new(class.Tag), new(class.Reply))

	orm.RunSyncdb("default", false, true)
}
