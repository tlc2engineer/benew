package s7

import (
	"benew/str"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const diFilename = "./files/csv/inputs.csv"
const aiFilename = "./files/csv/analog_input.csv"

/* Карта дискретных входов*/
var DiskrInputs map[int]str.Diskrete = make(map[int]str.Diskrete)
var AnalogInputs map[int]str.Analog = make(map[int]str.Analog)

/*Обработка входов. Начиная со второго байта и до 59.*/
func handleDInputs(msdata []byte) {
	for i := 2; i < 59; i++ {
		for step := 0; step < 8; step++ {
			inp, found := DiskrInputs[i*8+step]
			if found {
				inp.Val = msdata[i]&(byte(1<<uint(step))) > 0
				DiskrInputs[i*8+step] = inp
			}
		}
	}
}

/* Отображение заполняется значениеми false*/
func makeDInpMap() error {
	if file, err := os.Open(diFilename); err != nil {
		return err
	} else {
		defer file.Close()
		reader := csv.NewReader(file)
		if records, err := reader.ReadAll(); err == nil {
			for _, record := range records {
				num, err := strconv.Atoi(record[0])
				if num > 479 || num < 8 {
					return errors.New(fmt.Sprintf("Номер бита %d больше 479 или меньше 8", num))
				}
				if err != nil {
					return err
				}
				name := record[1]
				input := str.Diskrete{name, false}
				DiskrInputs[num] = input
			}
		} else {
			return err
		}

	}
	return nil
}

/*Проверка пересечения дискретных входов и аварий*/
func checkMapInputConsist() error {
	for k, _ := range DiskrInputs {
		_, found := AlarmMap[k]
		if found {
			return errors.New("Область входов и дискретных ошибок пересекается.")
		} // Области пересекаются
	}
	for k, v := range AnalogInputs {
		switch v.Tp {
		case "byte":
			for i := k; i < k+8; i++ {
				if _, found := AlarmMap[k]; found {
					return errors.New(fmt.Sprint("Область ошибок %d и аналоговых входов %s пересекается", k, v.Name))
				}
				if _, found := DiskrInputs[k]; found {
					return errors.New(fmt.Sprint("Область дискретных ошибок %d и аналоговых входов %s пересекается", k, v.Name))
				}
			}
		case "short":
			fallthrough
		case "float":
			for i := k; i < k+16; i++ {
				if _, found := AlarmMap[k]; found {
					return errors.New(fmt.Sprint("Область ошибок %d и аналоговых входов %s пересекается", k, v.Name))
				}
				if _, found := DiskrInputs[k]; found {
					return errors.New(fmt.Sprint("Область дискретных ошибок %d и аналоговых входов %s пересекается", k, v.Name))
				}
			}

		}

	}

	return nil
}

/* Составление карты аналоговых входов*/
func makeAnalogMap() error {
	file, err := os.Open(aiFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		addr, err := strconv.Atoi(record[0]) // адрес
		if addr > 479 || addr < 8 {
			return errors.New(fmt.Sprintf("Номер бита %d больше 479 или меньше 8", addr))
		}
		if addr%8 > 0 {
			return errors.New("Начальный бит не кратен 8, не начало байта")
		}
		if err != nil {
			return err
		}
		name := record[1]
		tp := record[2]
		coeff, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return err
		}
		AnalogInputs[addr] = str.Analog{name, tp, coeff, 0}
	}
	return nil
}

/* Считывание аналоговых входов*/
func handleAnalogInput(msdata []byte) {
	for k, ai := range AnalogInputs {
		ba := k / 8
		dat := msdata[ba]
		switch ai.Tp {
		case "byte":
			ai.Val = int(dat)
			AnalogInputs[k] = ai
		case "short":
			fallthrough
		case "float":
			val := int(dat) * 256
			val += int(msdata[ba+1])
			ai.Val = val
			AnalogInputs[k] = ai
		}

	}
}
