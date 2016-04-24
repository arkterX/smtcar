package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"smtcar/models"
)

type RoleController struct {
	beego.Controller
}

const (
	DefRolePageRow   = 10
	DefRolePageOrder = "Name"
)

func (this *RoleController) Get() {
	coord := models.RoleListCoord{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &coord); err != nil {
		// 没有坐标信息则用默认值
		coord.Page = 1
		coord.Row = DefRolePageRow
		coord.Sort = DefRolePageOrder
	}

	roles, count := models.GetRolelist(&coord)
	this.Data["json"] = &map[string]interface{}{"Total": count, "rows": &roles}
	this.ServeJSON()
}
