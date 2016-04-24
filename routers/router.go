package routers

import (
	"github.com/astaxie/beego"
	"smtcar/controllers/rbac"
)

func init() {
	beego.RESTRouter("/user", &controllers.UserController{})
	beego.RESTRouter("/role", &controllers.RoleController{})
}
