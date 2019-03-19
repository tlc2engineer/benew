package main

import (
	"benew/controllers"
	"benew/models"
	_ "benew/routers"
	"benew/s7"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

/* Номера потоков
1 - ЛКММ
2 - ПКММ
3 - TO

Левый поток
Адресс контроллера 192.168.211.70
Адрес компьютера 192.168.211.69

*/
func init() {

	cnf, err := config.NewConfig("ini", "./conf/contr.conf")
	if err != nil {
		fmt.Println("Ошибка чтения конфигурации")
	} else {
		s7.LAdress = cnf.String("laddr")
		s7.RAdress = cnf.String("raddr")
		s7.Reverse, _ = cnf.Bool("reverse")
		side := cnf.String("number")
		if err != nil {
			fmt.Println("Ошибка чтения конфигурации")
		} else {
			// чтение потоков
			switch side {
			case "left":
				s7.Potok = 1
			case "right":
				s7.Potok = 2
			case "termo":
				s7.Potok = 3
			}
		}
		controllers.Potok = s7.Potok
		models.Potok = s7.Potok

	}

}

func main() {
	drawSide()
	go s7.S7Manager()
	go models.DBManager()
	beego.Run()
}

func drawSide() {
	switch s7.Potok {
	case 1:
		fmt.Println("ЛКММ", s7.LAdress, s7.RAdress)
	case 2:
		fmt.Println("ПКММ", s7.LAdress, s7.RAdress)
	case 3:
		fmt.Println("КММТО", s7.LAdress, s7.RAdress)
	}
}
