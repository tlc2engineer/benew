package models

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

const fname = "./files/csv/temp.csv"

/*SaveTemp -
Сохранение временной переменной в csv файле.
*/
func SaveTemp(name, val string) {
	if file, err := os.Open(fname); err == nil {
		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			return
		}
		if file.Close() != nil {
			fmt.Println(err)
			return
		}
		for _, line := range lines {
			if line[0] == name {
				line[1] = val
			}
		}
		file, err := os.Create(fname)
		if err != nil {
			fmt.Println(err)
			return
		}

		writer := csv.NewWriter(file)
		for _, line := range lines {
			writer.Write(line)
		}
		writer.Flush()
		file.Close()
		if file.Close() != nil {
			fmt.Println(err)
			return
		}

	} else {
		fmt.Println(err)
	}
}

/*ReadString - чтение строки*/
func ReadString(name string) string {
	if file, err := os.Open(fname); err == nil {
		reader := csv.NewReader(file)
		lines, _ := reader.ReadAll()
		for _, line := range lines {
			if line[0] == name {
				return line[1]
			}
		}
	}
	return "0"
}

/*ReadInt - чтение int из string*/
func ReadInt(name string) int {
	val, _ := strconv.Atoi(ReadString(name))
	return val
}

/*ReadFloat - чтение float из string*/
func ReadFloat(name string) float64 {
	val, _ := strconv.ParseFloat(ReadString(name), 64)
	return val
}
