package main

import (
	"github.com/astaxie/beego"
	"smtcar/models"
	_ "smtcar/routers"
)

func main() {
	// 创建一个别名为default的数据库
	models.Syncdb(false, true)
	beego.Run()
}
