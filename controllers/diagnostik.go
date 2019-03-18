package controllers

import (
	"benew/models"
	"benew/str"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"sort"
	"strconv"
	"strings"

	"benew/s7"
)

type Diagn struct {
	beego.Controller
}

func (d *Diagn) Get() {
	d.Data["Title"] = "КММ"
	d.TplName = "diagnostik.html"
}

type DataFileLog struct {
	beego.Controller
}

type Clog struct {
	Date string
	Time string
	Log  string
}

func (sl *DataFileLog) Get() {
	tp := sl.GetString("type")
	fileName := models.S7FileName
	if tp == "db" {
		fileName = models.DbFileName
	}
	file, _ := os.Open(fileName)
	buffer := bytes.Buffer{}
	buffer.ReadFrom(file)
	s := buffer.String()
	lines := strings.Split(s, "\n")
	logs := make([]Clog, len(lines))
	count := 0
	for i, line := range lines {
		tokens := strings.Split(line, " ")
		if len(tokens) > 2 {
			count++
			logs[i] = Clog{}
			logs[i].Date = tokens[0]
			logs[i].Time = tokens[1]
			logs[i].Log = strings.Join(tokens[2:], " ")
		}
	}
	sl.Data["json"] = logs[:count]
	sl.ServeJSON()

}

type Alarm struct {
	beego.Controller
}

func (al *Alarm) Get() {
	al.Data["json"] = models.Alarms
	al.ServeJSON()
}

type ErrorList struct {
	beego.Controller
}

func (el *ErrorList) Get() {
	file, _ := os.Open("./files/csv/errors_cat.csv")
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Println(err)
	}
	descr := make([]ErrDescr, len(records))
	for i, record := range records {
		Id, _ := strconv.Atoi(record[0])
		descr[i] = ErrDescr{Id, record[1]}
	}
	el.Data["json"] = descr
	el.ServeJSON()
}

type ErrDescr struct {
	Id   int
	Text string
}

type MarkRL struct {
	beego.Controller
}

func (mrl *MarkRL) Get() {

	mrl.Data["json"] = models.MFArc
	mrl.ServeJSON()
}

/*Активные ошибки*/
type ActiveAlarms struct {
	beego.Controller
}

/*Тип для сортировки*/
type ByImp []str.Alarm

/*Длина*/
func (a ByImp) Len() int {
	return len(a)
}

/*Замена*/
func (a ByImp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

/*Функция сортировки*/
func (a ByImp) Less(i, j int) bool {
	if a[i].Atype != a[j].Atype {
		return a[i].Atype > a[j].Atype
	} else {
		return a[i].Time.After(a[j].Time)
	}
	return false
}

/*Активные ошибки с сортировкой*/
func (aa *ActiveAlarms) Get() {
	alarms := GetActiveAlarms(models.Alarms)
	aa.Data["json"] = alarms
	aa.ServeJSON()
}

func GetActiveAlarms(inAlarms []str.Alarm) []str.Alarm {
	alarms := make([]str.Alarm, 0)
	for _, alarm := range inAlarms {
		if alarm.Active {
			alarms = append(alarms, alarm)
		}
	}

	sort.Sort(ByImp(alarms))
	return alarms
}

/*Конкретные ошибки */
type CAlarms struct {
	beego.Controller
}

func (ca *CAlarms) Get() {
	ca.Data["json"] = s7.ActiveConcrete
	ca.ServeJSON()
}
