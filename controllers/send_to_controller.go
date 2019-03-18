package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	//"benew/s7"
	"benew/models"
	"benew/s7"
	"benew/str"
)

type SendController struct {
	beego.Controller
}
/* Прием маркировочных данных */
func (sc *SendController) Post() {
	id, err := sc.GetInt("id")
	if err != nil {
		fmt.Println(err)
	}
	krat, err := sc.GetInt("krat")
	mark := str.TemplateMark{}
	mark.Line1 = sc.GetString("mt[Line1]")
	mark.Line2 = sc.GetString("mt[Line2]")
	mark.Line3 = sc.GetString("mt[Line3]")
	mark.Line4 = sc.GetString("mt[Line4]")
	mark.Line5 = sc.GetString("mt[Line5]")
	mark.Line6 = sc.GetString("mt[Line6]")
	mark.Line7 = sc.GetString("mt[Line7]")
	mark.Line8 = sc.GetString("mt[Line8]")
	mark.Line9 = sc.GetString("mt[Line9]")
	mark.IDTemplate, err = sc.GetInt("mt[IDTemplate]")
	mark.TemplateName = sc.GetString("mt[TemplateName]")
	punch := str.TemplatePunch{}
	punch.TemplateName = sc.GetString("pt[TemplateName]")
	punch.IDTemplate, err = sc.GetInt("pt[IDTemplate]")
	punch.Line10 = sc.GetString("pt[Line10]")
	punch.Line11 = sc.GetString("pt[Line11]")
	punch.Line12 = sc.GetString("pt[Line12]")
	drawArownd, err := sc.GetBool("drawArownd")
	rotateText, err := sc.GetBool("rotateText")
	deep, err := sc.GetInt("deep")
	plate := models.GetPlate(id, krat)
	sendData := str.DataForMark{0, plate, mark, punch, drawArownd, rotateText, deep}
	s7.ChanMark <- sendData
	sc.Data["json"] = true
	sc.ServeJSON()
}
