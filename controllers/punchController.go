package controllers

import (
	"benew/models"
	"benew/s7"
	"benew/str"
	"fmt"
	"github.com/astaxie/beego"
	"image/png"
	"log"
	"strconv"
	"strings"
)

type PunchController struct {
	beego.Controller
}

func (pc *PunchController) Get() {
	/* Картинка клеймовки для актуального листа*/
	actual := pc.GetString("actual", "none")
	/* Получение картинки*/
	if actual != "none" {
		img := models.GetPaintImage(s7.ActualMarkData.GetPunchStrings())
		rw := pc.Ctx.ResponseWriter
		rw.Header().Set("Content-Type", "image/png")
		if err := png.Encode(rw, img); err != nil {
			log.Fatal(err)
		}
	}
	id, _ := strconv.Atoi(pc.GetString("id"))
	krat, _ := strconv.Atoi(pc.GetString("krat"))
	plate := models.GetPlate(id, krat)
	mp := str.CreateMap(plate)
	lines := make([]string, 3)
	for i := 10; i <= 12; i++ {
		line := pc.GetString(fmt.Sprintf("line%d", i))
		for key, value := range mp {
			line = strings.Replace(line, key, value, 5)
		}
		lines[i-10] = line

	}
	img := models.GetPaintImage(lines)
	rw := pc.Ctx.ResponseWriter
	rw.Header().Set("Content-Type", "image/png")
	if err := png.Encode(rw, img); err != nil {
		log.Fatal(err)
	}

}
