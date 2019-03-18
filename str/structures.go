package str

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

/*
Структура информации по листу, учитывая NULL.
*/
type PlateDataNull struct {
	Heat, Grade                                                                                                 sql.NullString
	Height                                                                                                      sql.NullFloat64
	Width, Length                                                                                               sql.NullInt64
	FabNumber, FabPos, OrderNumber, LotNumber, OrderPos, IdFabPos, RollFabId, OrderId, IdOrderPos, SlabId, Krat sql.NullInt64
	CNRS                                                                                                        sql.NullString
	PartyNum                                                                                                    sql.NullInt64
	NumInParty                                                                                                  sql.NullInt64
	SeqNum                                                                                                      sql.NullInt64
	M_pl                                                                                                        sql.NullString
	NewMark                                                                                                     sql.NullString
	Usl1                                                                                                        sql.NullString
	Usl2                                                                                                        sql.NullString
}

/*
Преобразование в стрктуру без учета NULL.
*/
func (pdn *PlateDataNull) Convert() PlateData {
	out := PlateData{}
	out.Heat = getValString(pdn.Heat, "")
	out.Grade = getValString(pdn.Grade, "")
	out.Height = getValFloat(pdn.Height, 0.0)
	out.Width = getValInt(pdn.Width, 0)
	out.Length = getValInt(pdn.Length, 0)
	out.FabNumber = getValInt(pdn.FabNumber, 0)
	out.FabPos = getValInt(pdn.FabPos, 0)
	out.OrderNumber = getValInt(pdn.OrderNumber, 0)
	out.LotNumber = getValInt(pdn.LotNumber, 0)
	out.OrderPos = getValInt(pdn.OrderPos, 0)
	out.IdFabPos = getValInt(pdn.IdFabPos, 0)
	out.RollFabId = getValInt(pdn.RollFabId, 0)
	out.OrderId = getValInt(pdn.OrderId, 0)
	out.IdOrderPos = getValInt(pdn.IdOrderPos, 0)
	out.SlabId = getValInt(pdn.SlabId, 0)
	out.Krat = getValInt(pdn.Krat, 0)
	out.CNRS = getValString(pdn.CNRS, "")
	out.PartyNum = getValInt(pdn.PartyNum, 0)
	out.NumInParty = getValInt(pdn.NumInParty, 0)
	out.SeqNum = getValInt(pdn.SeqNum, 0)
	out.M_pl = getValString(pdn.M_pl, "")
	out.NewMark = getValString(pdn.NewMark, "")
	out.Usl1 = getValString(pdn.Usl1, "")
	out.Usl2 = getValString(pdn.Usl2, "")
	return out
}
func getValString(str sql.NullString, def string) string {
	if !str.Valid {
		return def
	}
	return str.String
}
func getValInt(i sql.NullInt64, def int) int {
	if !i.Valid {
		return def
	}
	return int(i.Int64)
}
func getValFloat(f sql.NullFloat64, def float64) float64 {
	if !f.Valid {
		return def
	}
	return f.Float64
}

/*
Данные листа.
*/
type PlateData struct {
	Heat, Grade                                                                                                 string
	Height                                                                                                      float64
	Width, Length                                                                                               int
	FabNumber, FabPos, OrderNumber, LotNumber, OrderPos, IdFabPos, RollFabId, OrderId, IdOrderPos, SlabId, Krat int
	CNRS                                                                                                        string
	PartyNum                                                                                                    int
	NumInParty                                                                                                  int
	SeqNum                                                                                                      int
	M_pl                                                                                                        string
	NewMark                                                                                                     string
	Usl1                                                                                                        string
	Usl2                                                                                                        string
}

/*
Структура ID файла
*/
type Plate struct {
	SlabId, Krat, List int
}

/*
Шаблон клеймовки
*/
type TemplatePunch struct {
	IDTemplate             int
	TemplateName           string
	Line10, Line11, Line12 string
	Potok                  int
}

/**
Шаблон маркировки
*/
type TemplateMark struct {
	IDTemplate                                                    int
	TemplateName                                                  string
	Line1, Line2, Line3, Line4, Line5, Line6, Line7, Line8, Line9 string
	Potok                                                         int
}

/*
Шаблон возвращаемых слябов.
*/

type RetPlate struct {
	Id, Krat, List int
}

func CreateMap(plate PlateData) map[string]string {
	mp := make(map[string]string)
	mp["<ID сляба>"] = fmt.Sprintf("%d", plate.SlabId)
	mp["<Крат>"] = fmt.Sprintf("%d", plate.Krat)
	mp["<Заказ>"] = fmt.Sprintf("%d", plate.OrderNumber)
	mp["<Плавка>"] = fmt.Sprintf("%s", plate.Heat)
	mp["<Марка>"] = fmt.Sprintf("%s", plate.Grade)
	mp["<№ Партии>"] = fmt.Sprintf("%d", plate.PartyNum)
	mp["<Лот>"] = fmt.Sprintf("%d", plate.LotNumber)
	mp["<Размеры>"] = fmt.Sprintf("%.2fx%dx%d", plate.Height, plate.Width, plate.Length)
	return mp
}

/* Данные для маркировки*/
type DataForMark struct {
	Num                    uint16
	Plate                  PlateData
	Tm                     TemplateMark
	Tp                     TemplatePunch
	DrawAround, RotateText bool
	Deep                   int
}

/* Получение картинки маркировки*/
func (dfm *DataForMark) GetPaintStrings() []string {
	mp := CreateMap(dfm.Plate)
	mlines := []string{dfm.Tm.Line1, dfm.Tm.Line2, dfm.Tm.Line3, dfm.Tm.Line4, dfm.Tm.Line5, dfm.Tm.Line6, dfm.Tm.Line7, dfm.Tm.Line8, dfm.Tm.Line9}
	lines := make([]string, 9)
	for i := 1; i <= 9; i++ {
		line := mlines[i-1]
		for key, value := range mp {
			line = strings.Replace(line, key, value, 5)
		}
		lines[i-1] = line

	}
	return lines
}

/* Получение картинки клеймовки*/
func (dfm *DataForMark) GetPunchStrings() []string {
	mp := CreateMap(dfm.Plate)
	plines := []string{dfm.Tp.Line10, dfm.Tp.Line11, dfm.Tp.Line12}
	lines := make([]string, 3)
	for i := 1; i <= 3; i++ {
		line := plines[i-1]
		for key, value := range mp {
			line = strings.Replace(line, key, value, 5)
		}
		lines[i-1] = line

	}
	return lines
}

/*Телеграмма маркировки*/
type MF struct {
	Id                      int
	SlabId                  int
	Krat                    int
	PunchResult, MarkResult int
	T                       time.Time
	Num                     int
}

/* МashineState - cocтoяние KMM*/
type MS struct {
	ActPlate                                               PlateData
	MachineMode, DataStatus, State, PaintState, PunchState int
	Alarms                                                 [23]AlarmType
	Time                                                   time.Time
	DBConnect, ControllerConnect                           bool
	Error                                                  bool
	ActAlarms                                              []Alarm
	ConcreteAlarms                                         []ConcreteAlarm // Список конкретных ошибок
	Inputs                                                 map[string]interface{}
}

type LogType int

const (
	DB LogType = iota
	S7
)

/*
Авария.
*/
type Alarm struct {
	UniqId  int
	Time    time.Time
	Num     int
	Atype   AlarmType
	TimeExp time.Time
	Active  bool
}

/*
Тип аварии.
*/
type AlarmType int

const (
	NoAlarm       AlarmType = 0 // Нет аварии
	Warning       AlarmType = 1 // Предупреждение
	Fault         AlarmType = 2 // Авария
	CriticalAlarm AlarmType = 3 // Критическая авария
)

/*Конкретная авария*/
type ConcreteAlarm struct {
	Zone    int    // Зона ошибки. Напряжение, передача данных и т.д.
	Place   int    // место ошибки. Общее, клеймовка, маркировка
	Num     int    // номер ошибки
	TextNum string // текстовый номер
	Text    string // текст ошибки

}

/* Дискретный вход*/
type Diskrete struct {
	Name string
	Val  bool
}

/* Аналоговый вход*/
type Analog struct {
	Name  string
	Tp    string
	Coeff float64
	Val   int
}

/* Целое значение аналогового входа*/
func (an *Analog) GetInt() int {
	return an.Val
}

/* Вещественное значение аналогового входа*/
func (an *Analog) GetFloat() float64 {
	return float64(an.Val) * an.Coeff
}

/* реализация FIFO для массива байт*/
type FIFO struct {
	buff [][]byte
}

func NewFIFO(size int) FIFO {
	fifo := FIFO{}
	fifo.buff = make([][]byte, size, size+2)
	return fifo
}

func (ff *FIFO) Push(val []byte) {
	buff := ff.buff
	copy(buff[1:], buff[:len(buff)-2])
	buff[0] = val
}
