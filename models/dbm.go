package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/beego/admin/src/lib"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var ormer orm.Ormer
var db_type = beego.AppConfig.String("db_type")
var db_host = beego.AppConfig.String("db_host")
var db_port = beego.AppConfig.String("db_port")
var db_user = beego.AppConfig.String("db_user")
var db_pass = beego.AppConfig.String("db_pass")
var db_name = beego.AppConfig.String("db_name")
var db_path = beego.AppConfig.String("db_path")
var db_sslmode = beego.AppConfig.String("db_sslmode")

//数据库连接
func connect() {
	var dns string

	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
		break
	case "postgres":
		orm.RegisterDriver("postgres", orm.DRPostgres)
		dns = fmt.Sprintf("dbname=%s host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_name, db_host, db_user, db_pass, db_port, db_sslmode)
	case "sqlite3":
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
		if db_path == "" {
			db_path = "./"
		}
		dns = fmt.Sprintf("%s%s.db", db_path, db_name)
		break
	default:
		beego.Critical("Database driver is not allowed:", db_type)
	}
	orm.RegisterDataBase("default", db_type, dns)
}

//创建数据库
func createdb() (err error) {
	var dns string
	var sqlstring string
	switch db_type {
	case "mysql":
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
		sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
		break
	case "postgres":
		dns = fmt.Sprintf("host=%s  user=%s  password=%s  port=%s  sslmode=%s", db_host, db_user, db_pass, db_port, db_sslmode)
		sqlstring = fmt.Sprintf("CREATE DATABASE %s", db_name)
		break
	case "sqlite3":
		if db_path == "" {
			db_path = "./"
		}
		dns = fmt.Sprintf("%s%s.db", db_path, db_name)
		os.Remove(dns)
		sqlstring = "create table init (n varchar(32));drop table init;"
		break
	default:
		beego.Critical("Database driver is not allowed:", db_type)
	}

	beego.Info("sql:" + dns)
	db, err := sql.Open(db_type, dns)
	if err != nil {
		panic(err.Error())
	}

	r, err := db.Exec(sqlstring)
	if err != nil {
		beego.Debug(err)
		beego.Debug(r)
	} else {
		log.Println("Database ", db_name, " created")
	}

	defer db.Close()
	return nil
}

func initAdminUser() {
	fmt.Println("insert user ...")
	u := new(User)
	u.Username = "admin"
	u.Nickname = "ClownFish"
	u.Password = Pwdhash("admin")
	u.Email = "osgochina@gmail.com"
	u.Remark = "I'm admin"
	u.Status = 2
	ormer = orm.NewOrm()
	if created, id, err := ormer.ReadOrCreate(u, "Username"); err == nil {
		if created {
			fmt.Println("New Insert an user. Id:", id)
		} else {
			fmt.Println("Admin user is exist. Id:", id)
		}
	}
	fmt.Println("insert user end")
}

func initAdminRole() {
	fmt.Println("insert role ...")
	r := new(Role)
	r.Name = "Admin"
	r.Remark = "I'm a admin role"
	r.Status = 2
	r.Title = "Admin role"
	if created, id, err := ormer.ReadOrCreate(r, "Name"); err == nil {
		if created {
			fmt.Println("New Insert an role. Id:", id)
		} else {
			fmt.Println("Admin role is exist. Id:", id)
		}
	}

	fmt.Println("insert role end")
}

// name means table's alias name. default is "default".
// force means run next sql if the current is error.
// verbose means show all info when running command or not.
func Syncdb(force bool, verbose bool) {
	createdb()
	connect()
	ormer = orm.NewOrm()

	// 遇到错误立即返回
	err := orm.RunSyncdb("default", force, verbose)
	if err != nil {
		fmt.Println(err)
		return
	}

	initAdminUser()
	initAdminRole()
}
