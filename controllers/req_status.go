package controllers

import (
	"benew/s7"
	"benew/str"
	"github.com/astaxie/beego"
	"time"

	"benew/models"
	//"fmt"
)

/*Запрос статуса*/
type RS struct {
	beego.Controller
}

func (rs *RS) Get() {
	s7.CommandChan <- str.GetMS
	select {
	case dat := <-s7.RSData:
		dat.ActAlarms = GetActiveAlarms(models.Alarms) // активные аварии
		rs.Data["json"] = dat
		rs.ServeJSON()
	case <-time.After(100 * time.Millisecond):
		dat := str.MS{}
		dat.Error = true
		rs.Data["json"] = dat
		rs.ServeJSON()
	}
}
