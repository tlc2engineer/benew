package controllers

import (
	"github.com/astaxie/beego"

	"benew/models"
	"benew/s7"
	"benew/str"
	"fmt"
	"image/png"
	"log"
	"strconv"
	"strings"
)

type PaintController struct {
	beego.Controller
}

func (pc *PaintController) Get() {
	/* Картинка маркировки для акуального листа*/
	actual := pc.GetString("actual", "none")
	if actual != "none" {
		img := models.GetPaintImage(s7.ActualMarkData.GetPaintStrings())
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
	lines := make([]string, 9)
	for i := 1; i <= 9; i++ {
		line := pc.GetString(fmt.Sprintf("line%d", i))
		for key, value := range mp {
			line = strings.Replace(line, key, value, 5)
		}
		lines[i-1] = line

	}
	img := models.GetPaintImage(lines)
	rw := pc.Ctx.ResponseWriter
	rw.Header().Set("Content-Type", "image/png")
	if err := png.Encode(rw, img); err != nil {
		log.Fatal(err)
	}

}
