package controllers

import (
	"benew/s7"

	"github.com/astaxie/beego"
)

type ConfController struct {
	beego.Controller
}

func (cc *ConfController) Get() {
	cc.Data["Title"] = "Конфигурация"
	cc.Data["ip_comp"] = s7.LAdress
	cc.Data["ip_contr"] = s7.RAdress
	switch s7.Potok {
	case 1:
		cc.Data["stream"] = "Левый поток"
	case 2:
		cc.Data["stream"] = "Правый поток"
	case 3:
		cc.Data["stream"] = "TO"
	}

	cc.TplName = "conf.html"
}
