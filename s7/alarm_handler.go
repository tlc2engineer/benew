package s7

import (
	"benew/models"
	"benew/str"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const fname = "./files/csv/err_descr.csv"

var AlarmMap map[int]str.ConcreteAlarm = make(map[int]str.ConcreteAlarm)

/*Получение конкретных аварий. На вход подается состояние машины MS присланное из контроллера.
На выходе - список конкретных аварий.
*/
func handleConcreteAlarms(msdata []byte) []str.ConcreteAlarm {
	alarms := make([]str.ConcreteAlarm, 0)
	for i := 14; i <= 59; i++ {
		if msdata[i] > 0 {
			for step := 0; step < 8; step++ {
				if (msdata[i] & (1 << uint(step))) > 0 {
					alarm, found := AlarmMap[i*8+step]
					if found {
						alarms = append(alarms, alarm)
					}
				}
			}
		}
	}
	return alarms
}

/* Получение зоны аварии по номеру байта*/
func getZone(i int) int {
	return (i-14)/2 + 1
}

/* Функция инициализации модуля*/
func init() {
	// Карта аварий.
	viewErr(makeAlarmMap())
	// Составление карты дискретных входов.
	viewErr(makeDInpMap())
	// Составление карты аналоговых входов.
	viewErr(makeAnalogMap())
	// Проверка пересечения
	viewErr(checkMapInputConsist())
	// Сохраненный  номер телеграммы
	MTNumber = uint16(models.ReadInt("MarkTelCount"))
	fmt.Println("mt number", MTNumber)

}

func viewErr(err error) {
	if err != nil {
		panic(err)
	}
}

func makeAlarmMap() error {
	if file, err := os.Open(fname); err != nil {
		return err
	} else {
		defer file.Close()
		reader := csv.NewReader(file)
		reader.Comma = ';'
		if records, err := reader.ReadAll(); err != nil {
			return err
		} else {
			for _, record := range records {
				if num, err := strconv.Atoi(record[0]); err != nil {
					return err
				} else {
					alarm := str.ConcreteAlarm{}
					alarm.TextNum = record[1]
					descr := strings.Split(record[1], "-")
					// Зона ошибки
					zone, err := strconv.Atoi(descr[0])
					if err != nil {
						return err
					}
					alarm.Zone = zone
					// Место
					place, err := strconv.Atoi(descr[1])
					if err != nil {
						return err
					}
					alarm.Place = place
					// Номер
					anum, err := strconv.Atoi(descr[2])
					if num > 479 || num < 8 {
						return errors.New(fmt.Sprintf("Номер бита %d больше 479 или меньше 8", num))
					}
					if err != nil {
						return err
					}
					alarm.Num = anum
					alarm.Text = record[2]
					AlarmMap[int(num)] = alarm
				}

			}
		}

	}
	return nil
}
