package controllers

import (
	"benew/s7"
	"fmt"
	"github.com/astaxie/beego"
)

type RecvController struct {
	beego.Controller
}

func (rc *RecvController) Post() {
	mess := rc.GetString("ms")
	val, _ := rc.GetBool("val")
	switch mess {
	case "on":
		s7.NumCommand = 1
	case "off":
		s7.NumCommand = 2
	case "start":
		s7.NumCommand = 3
	case "stop":
		s7.NumCommand = 4
	case "plate_on":
		s7.NumCommand = 5
	case "mark_contact":
		s7.NumCommand = 7
	case "mark_up_pos":
		s7.NumCommand = 9
	case "punch_contact":
		s7.NumCommand = 11
	case "punch_up_pos":
		s7.NumCommand = 13
	}
	if val {
		s7.NumCommand = s7.NumCommand + 1
	}
	s7.FlagRecvCommand = true
	fmt.Println(mess, val)
	rc.Data["json"] = true
	rc.ServeJSON()
}
