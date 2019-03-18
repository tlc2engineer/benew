package controllers

import (
	models "benew/models"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"strings"
	"time"

	data "benew/str"

)

var Potok=1

var Offline = false                  // Флаг работы в оффлайн.
var offlineTable []data.PlateData    // Список листов оффлайн
var offlineMT []models.TemplateMark  // список шаблонов маркировки offline
var offlinePT []models.TemplatePunch // список шаблонов клеймовки offline



/*
Тестовое окно
*/
type MainController struct {
	beego.Controller
}

/*
Тестовый контроллер. Метод Get.
*/
func (c *MainController) Get() {
	c.Data["Title"] = "Главная страница"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

/*
КММ контроллер. Отображение главного окна.
*/
type KmmController struct {
	beego.Controller
}

/*
Метод Get.
*/
func (c *KmmController) Get() {
	dates := make([]string, 4)
	for i := 0; i < 4; i++ {
		dt := -int64(i) * int64(time.Hour) * 24 // Для совместимости с 32 разрядной версией. dt:=-i*int(time.Hour)*24
		dates[i] = formTime(time.Now().Add(time.Duration(dt)))
	}
	c.Data["Title"] = "Передача данных"
	c.Data["Dates"] = dates
	c.TplName = "kmm.html"

}

/* Форматирование времени*/
func formTime(t time.Time) string {
	return fmt.Sprintf("%d.%d.%d", t.Day(), t.Month(), t.Year())
}

/*Таблица листов для маркировки*/
type TableData struct {
	beego.Controller
}

func (c *TableData) Get() {
	slab_id := c.GetString("slab_num")
	date := c.GetString("date")
	if date == "Выберите дату" {
		date = "12/12/16"
	}
	all := c.GetString("all")
	var plates []data.PlateData
	if !Offline {
		fmt.Println(slab_id, date, all)
		plates = models.ReturnTableData(strings.Replace(date, ".", "/", 3), true, false, true, slab_id)
		// Запись в файл для наладки
		if file, err := os.Create("./files/offline/tdata.dat"); err != nil {
			fmt.Println(err)
		} else {
			encoder := gob.NewEncoder(file)
			if err = encoder.Encode(plates); err != nil {
				fmt.Println(err)
			}
		}
	} else {
		if file, err := os.Open("./files/offline/tdata.dat"); err != nil {
			fmt.Println(err)
		} else {
			decoder := gob.NewDecoder(file)
			decoder.Decode(&plates)
			offlineTable = plates
		}

	}
	c.Data["json"] = plates
	c.ServeJSON()
}

/* Пояснения к листу и дополнительная информация.*/

type Remark struct {
	beego.Controller
}

func (r *Remark) Get() {
	slabId, err := r.GetInt("SlabId")
	if err != nil {
		fmt.Println(err)
	}
	krat, err := r.GetInt("Krat")
	if err != nil {
		fmt.Println(err)
	}
	remarks := models.GetRemarksId(slabId, krat)
	r.Data["json"] = remarks
	r.ServeJSON()

}

type MarkTemplateById struct {
	beego.Controller
}

type PunchTemplateById struct {
	beego.Controller
}

func (mt *MarkTemplateById) Get() {
	var markT models.TemplateMark
	/* Возврат первого маркировки шаблона из памяти*/
	getOfflineFirstT := func() models.TemplateMark {
		if len(offlineMT) > 0 {
			return offlineMT[8]
		} else {
			var list []models.TemplateMark
			getMTLOff(&list)
			return offlineMT[8]
		}
	}
	if !Offline {
		idOrderPos, err := mt.GetInt("idOrderPos")
		if err != nil {
			fmt.Println(err)
		}
		markT = models.GetMarkTemplateBId(idOrderPos, Potok)
		if markT.Line1 == "" {
			markT = getOfflineFirstT()
		}
	} else {
		getOfflineFirstT()

	}
	mt.Data["json"] = markT
	mt.ServeJSON()
}

func (pt *PunchTemplateById) Get() {
	var punchT models.TemplatePunch
	getOfflineFirst := func() models.TemplatePunch {
		if len(offlinePT) > 0 {
			fmt.Println("offline PT", offlinePT)
			return offlinePT[1]
		} else {
			var list []models.TemplatePunch
			getPTLOff(&list)
			return list[1]
		}
	}
	if !Offline {
		idOrderPos, err := pt.GetInt("idOrderPos")
		if err != nil {
			fmt.Println(err)
		}
		punchT = models.GetPunchTemplateById(idOrderPos, Potok)
		if punchT.IDTemplate == 0 {
			punchT = getOfflineFirst()
		}
	} else {
		punchT = getOfflineFirst()

	}
	pt.Data["json"] = punchT
	pt.ServeJSON()
}

type ListMarkTamplates struct {
	beego.Controller
}

func (lmt *ListMarkTamplates) Get() {
	var list []models.TemplateMark // Список шаблонов маркировки
	if !Offline {
		list = models.GetMarkTemplate(Potok)
		// Запись в файл для наладки
		if file, err := os.Create("./files/offline/mt.dat"); err != nil {
			fmt.Println(err)
		} else {
			enc := gob.NewEncoder(file)
			enc.Encode(list)
		}
	} else {
		getMTLOff(&list)

	}
	lmt.Data["json"] = list
	lmt.ServeJSON()
}

/* Загрузка списка шаблонов маркировки из ранее сохраненного файла*/
func getMTLOff(list *[]models.TemplateMark) error {
	if file, err := os.Open("./files/offline/mt.dat"); err != nil {
		return err
	} else {
		decoder := gob.NewDecoder(file)
		if err = decoder.Decode(&list); err != nil {
			return err
		}
	}
	offlineMT = *list
	return nil
}

/* Список шаблонов клеймовки*/
type ListPunchTamplates struct {
	beego.Controller
}

func (lpt *ListPunchTamplates) Get() {
	var list []models.TemplatePunch // Список шаблонов клеймовки
	if !Offline {
		list = models.GetPunchTemplate(Potok)
		// Запись в файл для наладки
		file, err := os.Create("./files/offline/pt.dat")
		if err != nil {
			fmt.Println(err)
		} else {
			enc := gob.NewEncoder(file)
			enc.Encode(list)
		}
	} else {
		getPTLOff(&list)
	}
	lpt.Data["json"] = list
	lpt.ServeJSON()
}

/* Загрузка списка шаблонов клеймовки из ранее сохраненного файла*/
func getPTLOff(list *[]models.TemplatePunch) error {
	if file, err := os.Open("./files/offline/pt.dat"); err != nil {
		return err
	} else {
		decoder := gob.NewDecoder(file)
		if err = decoder.Decode(&list); err != nil {
			return err
		}
	}
	offlinePT = *list
	return nil
}

/*Удаление шаблона маркировки*/
type DeleteMTemplate struct {
	beego.Controller
}

func (dmt *DeleteMTemplate) Post() {
	num, err := dmt.GetInt("num")
	if err != nil {
		fmt.Println("parse err", err)
		dmt.Data["json"] = false
		dmt.ServeJSON()
	} else {
		dmt.Data["json"] = models.DeleteMarkTemplate(num)
		dmt.ServeJSON()
	}

}

type SaveMTemplate struct {
	beego.Controller
}

func (smt *SaveMTemplate) Post() {

	tmpl := models.TemplateMark{}
	tmpl.TemplateName = smt.GetString("name")
	fmt.Println("Save  ", tmpl.TemplateName)
	tmpl.Line1, tmpl.Line2, tmpl.Line3, tmpl.Line4, tmpl.Line5, tmpl.Line6, tmpl.Line7, tmpl.Line8, tmpl.Line9 = smt.GetString("line1"), smt.GetString("line2"),
		smt.GetString("line3"), smt.GetString("line4"), smt.GetString("line5"), smt.GetString("line6"),
		smt.GetString("line7"), smt.GetString("line8"), smt.GetString("line9")
	smt.Data["json"] = models.SaveTemplateMark(tmpl)
	smt.ServeJSON()
}

type DeletePTemplate struct {
	beego.Controller
}

func (dpt *DeletePTemplate) Post() {
	num, err := dpt.GetInt("num")
	if err != nil {
		fmt.Println("parse err", err)
		dpt.Data["json"] = false
		dpt.ServeJSON()
	} else {
		dpt.Data["json"] = models.DeletePunchTemplate(num)
		dpt.ServeJSON()
	}

}

type SavePTemplate struct {
	beego.Controller
}

func (smt *SavePTemplate) Post() {

	tmpl := models.TemplatePunch{}
	tmpl.TemplateName = smt.GetString("name")
	tmpl.Line10, tmpl.Line11, tmpl.Line12 = smt.GetString("line10"), smt.GetString("line11"), smt.GetString("line12")
	smt.Data["json"] = models.SaveTemplatePunch(tmpl)
	smt.ServeJSON()
}

/*
Список замаркированных слябов.
*/

type RetListPlates struct {
	beego.Controller
}

func (rp *RetListPlates) Get() {
	rp.Data["json"] = models.RetPlates()
	rp.ServeJSON()

}

/*
Возврат замаркированного сляба.
*/

type RetPlate struct {
	beego.Controller
}

func (rp *RetPlate) Get() {
	id, err := rp.GetInt("id")
	if err != nil {
		panic(err)
	}
	rp.Data["json"] = models.RetPlate(id)
	rp.ServeJSON()

}

/*
Словарь
*/

type DictController struct {
	beego.Controller
}

func (dc *DictController) Get() {
	dc.Data["json"] = models.ReadDict()
	dc.ServeJSON()
}

func (dc *DictController) Post() {
	text := dc.GetString("text")
	dc.Data["json"] = models.WriteDict(text)
	dc.ServeJSON()
}
