package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"smtcar/models"
)

type UserController struct {
	beego.Controller
}

const (
	DefUserPageRow   = 10
	DefUserPageOrder = "UserName"
)

func (this *UserController) Get() {
	coord := models.UserListCoord{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &coord); err != nil {
		// 没有坐标信息则用默认值
		coord.Page = 1
		coord.Row = DefUserPageRow
		coord.Sort = DefUserPageOrder
	}

	users, count := models.GetUserlist(&coord)
	this.Data["json"] = &map[string]interface{}{"Total": count, "rows": &users}
	this.ServeJSON()
}

func (this *UserController) Post() {
	user := models.User{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &user); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("user coordinate err"))
		return
	}
	_, err := models.AddUser(&user)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}
