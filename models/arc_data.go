package models

import (
	"benew/str"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

/*MSArc - список телеграмм статуса KMM*/
var MSArc []str.MS = make([]str.MS, 200)

/*MFArc - список телеграмм результатов маркировки*/
var MFArc []str.MF
var file *os.File

const maxCountMS = 100
const maxCountMF = 100

/*DbFileName - путь к файлу хранящему логи базы данных*/
const DbFileName = "./files/logs/db.log"

/*S7FileName - путь к файлу хранящему логи связи с контроллером S7*/
const S7FileName = "./files/logs/s7.log"

/*S7Logger - Запись ошибок контроллера*/
var S7Logger log.Logger

/*DBLogger - Запись ошибок базы данных*/
var DBLogger log.Logger

/*Alarms - список текущих аварий*/
var Alarms []str.Alarm

const alarmFilename = "./files/saveData/alarms.dat" // имя файла хранящего аварии
const markFilename = "./files/saveData/mark.dat"

var mid = 0  // максимальный номер сохраненного аварийного сообщения
var mmid = 0 //максимальный номер сохраненной телеграммы маркировки

/*Test - тестирование*/
func Test() {
	fmt.Println("Good news")
}

/*AddMS - добавить телеграмму статуса КММ*/
func AddMS(ms str.MS) {
	AlarmHandler(ms)
	MSArc = append(MSArc, ms)
	if len(MSArc) > 100 {
		MSArc = MSArc[1:]
	}
}

/*AddMF - добавить маркировочную телеграмму*/
func AddMF(mf str.MF) {
	mf.Id = mmid
	mmid++
	MFArc = append(MFArc, mf)
	if len(MFArc) > 100 {
		MFArc = MFArc[1:]
	}
	saveMarkResult()
}

func init() {
	fmt.Println("init")

	s7file, err := os.OpenFile(S7FileName, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	S7Logger = *log.New(s7file, "", 3)
	dbfile, err := os.OpenFile(DbFileName, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	DBLogger = *log.New(dbfile, "", 3)
	Alarms = make([]str.Alarm, 0)
	restoreAlarmsFromDisk()                                   // Восстановление аварий с диска.
	if mmid, MFArc, err = restoreMarkFromDisk(); err != nil { // Восстановление телеграмм с диска
		MFArc = make([]str.MF, 0)
		mmid = 0
	}

}

/*AddLog -
Добавление сообщения
*/
func AddLog(logType str.LogType, text string) {
	switch logType {
	case str.DB:
		DBLogger.Println(text)
	case str.S7:
		S7Logger.Println(text)

	}
}

/*AlarmHandler -
Обработка ошибок в сообщении.
*/
func AlarmHandler(ms str.MS) {

	for i, alarm := range ms.Alarms { // Просмотр списка ошибок в сообщениии
		if alarm > 0 { // активные ошибки
			fmt.Println("Active", i)
			exists := false // Флаг существует в false
			for j, eAlarm := range Alarms {
				if eAlarm.Num == i && eAlarm.Active { // активную ошибку с таким номером.
					if alarm == eAlarm.Atype { // Есть такая ошибка с таким же приоритетом
						exists = true // есть такая ошибка, висит дальше
						break
					}
					if alarm != eAlarm.Atype { // приоритет изменился
						Alarms[j].Active = false       // старая ошибка неактивна
						Alarms[j].TimeExp = time.Now() // фиксируем время исчезновения старой ошибки
					}
				}
			}
			if !exists { // Если нет такой ошибки. Новая ошибка.
				newAlarm := NewAlarm(i, alarm)    // Создаем аварию
				Alarms = append(Alarms, newAlarm) // Добавление в список.
				if len(Alarms) > 1000 {
					Alarms = Alarms[1:] // Отбрысываем последнюю сохраненную
				}
				fmt.Println(Alarms)
				saveAlarmOnDisk() // Сохранение на диск!

			}
		} else { // если ошибка исчезла
			for j, eAlarm := range Alarms { //перебор списка ошибок
				if eAlarm.Num == i && eAlarm.Active { // есть активная ошибка с тем же номером
					Alarms[j].Active = false       // ошибка неактивна
					Alarms[j].TimeExp = time.Now() // фиксируем время исчезновения
					saveAlarmOnDisk()              // Сохранение на диск!
				}
			}
		}
	}
}

/*NewAlarm Новая авария*/
func NewAlarm(num int, atype str.AlarmType) str.Alarm {
	newAlarm := str.Alarm{} // Создаем
	newAlarm.Num = num      // Номер
	newAlarm.Atype = atype  // Тип
	newAlarm.Active = true  // ошибка активна
	newAlarm.Time = time.Now()
	mid++ // увеличение номера на 1
	newAlarm.UniqId = mid
	return newAlarm
}

/*
Тип данных для сохранения на диск
*/
type saveAlarm struct {
	Mid    int         // максимальный номер сохраненный
	Alarms []str.Alarm // список ошибок
}

/*
Сохранение аварий на диск
*/
func saveAlarmOnDisk() {
	file, err := os.Create(alarmFilename)
	if err != nil {
		fmt.Println(err)
	}
	encoder := gob.NewEncoder(file)
	saveAlarm := saveAlarm{}
	saveAlarm.Mid = mid
	saveAlarm.Alarms = Alarms
	err = encoder.Encode(saveAlarm)
	if err != nil {
		fmt.Println(err)
	}
}

/*
Восстановление аварий с диска
*/
func restoreAlarmsFromDisk() {
	file, err := os.Open(alarmFilename)
	if err != nil {
		fmt.Println(err)
	}
	decoder := gob.NewDecoder(file)
	var saveAlarm saveAlarm
	err = decoder.Decode(&saveAlarm)
	if err != nil {
		fmt.Println(err)
	}
	mid = saveAlarm.Mid

	Alarms = append(Alarms, saveAlarm.Alarms...)
}

/*SaveMarkDat Структура телеграмм маркировки*/
type SaveMarkDat struct {
	MaxId int
	Data  []str.MF
}

/* Сохранение телеграмм маркировки на диск*/
func saveMarkResult() {
	if file, err := os.Create(markFilename); err != nil {
		fmt.Println(err)
	} else {
		encoder := gob.NewEncoder(file)
		saveMarkDat := SaveMarkDat{}
		saveMarkDat.MaxId = mmid
		saveMarkDat.Data = MFArc
		if err = encoder.Encode(saveMarkDat); err != nil {
			fmt.Println(err)
		}
	}

}

/*
Восстановление телеграмм маркировки с диска
*/
func restoreMarkFromDisk() (int, []str.MF, error) {
	if file, err := os.Open(markFilename); err != nil {
		return -1, nil, err
	} else {
		decoder := gob.NewDecoder(file)
		var saveMarkDat SaveMarkDat
		if err = decoder.Decode(&saveMarkDat); err != nil {
			return -1, nil, err
		} else {
			return saveMarkDat.MaxId, saveMarkDat.Data, nil
		}
	}

}
