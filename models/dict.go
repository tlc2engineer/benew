package models

import (
	"io/ioutil"
	"os"
	"strings"
)

const FileName = "./files/dict.txt"

func ReadDict() []string {
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		panic(err)
	}
	s := string(data)
	return strings.Split(s, "\n")

}

func WriteDict(text string) bool {
	err := ioutil.WriteFile(FileName, []byte(text), os.ModeDevice)
	if err != nil {
		return false
	}
	return true
}
