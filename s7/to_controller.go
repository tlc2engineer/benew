package s7

import (
	"benew/models"
	"benew/str"
	"bytes"
	"errors"
	"fmt"

	"net"
	"strings"
	"time"
)

var (
	FlagRecvCommand bool
	NumCommand      byte
	Connected       = false
	ChanMark        = make(chan str.DataForMark)
	CommandChan     = make(chan str.CommandAction, 2)
	RSData          = make(chan str.MS)
	MFData          = make(chan str.MF)
	ActiveConcrete  []str.ConcreteAlarm
	MTNumber        uint16
	MTelBuff        []str.DataForMark = make([]str.DataForMark, 10)
	ActualMarkData  str.DataForMark   = str.DataForMark{}
	Potok                             = 1
	LAdress                           = "192.168.51.48:2000"
	RAdress                           = "192.168.51.66:2000"
	Reverse                           = false
)

const SendLength = 2698
const RecvLength = 60

func main() {

	status, err := GetMashineStatus()
	if err != nil {

	}
	fmt.Println("Статус:", *status)

}

/*
Получение соединения.
*/
func getConnection() (*net.TCPConn, error) {
	laddr, err := net.ResolveTCPAddr("tcp", LAdress)
	if err != nil {
		return nil, err
	}
	raddr, err := net.ResolveTCPAddr("tcp", RAdress)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", laddr, raddr)
	return conn, err
}

type Int int // 4 байта

func (intv Int) getBytes() []byte {
	ret := make([]byte, 4)
	ret[3] = byte(int(intv) & 0xff)
	ret[2] = byte((int(intv) & 0xff00) >> 8)
	ret[1] = byte((int(intv) & 0xff0000) >> 16)
	ret[0] = byte((int64(intv) & 0xff000000) >> 24) // int64 для совместимости с 32 разрядной версией.
	return ret
}

type Short int // 2 байта

func (shortv Short) getBytes() []byte {
	ret := make([]byte, 2)
	ret[1] = byte(int(shortv) & 0xff)
	ret[0] = byte((int(shortv) & 0xff00) >> 8)
	return ret
}

type S7Str string

func (str S7Str) getBytes() []byte {
	ret := make([]byte, 28)
	ret[0] = 28
	ret[1] = byte(len(str))
	buff := bytes.NewBuffer(ret)
	buff.Truncate(2)
	buff.Write([]byte(str))
	return ret
}

/*	Статус машины. Новый объект. */
func NewMS(rsdata []byte) str.MS {
	ms := str.MS{}
	ms.MachineMode = int(rsdata[5])
	ms.DataStatus = int(rsdata[7])
	ms.State = int(rsdata[9])
	ms.PunchState = int(rsdata[11])
	ms.PaintState = int(rsdata[13])
	count := 0
	for i := 15; i <= 51; i += 2 {
		ms.Alarms[count] = str.AlarmType(int(rsdata[i] & 3))
		count++
	}
	ms.DBConnect = models.DBConnected
	ms.ControllerConnect = Connected
	ms.Time = time.Now()
	ActiveConcrete = handleConcreteAlarms(rsdata)
	ms.ConcreteAlarms = ActiveConcrete
	handleDInputs(rsdata) // Обработка дискретных входов

	handleAnalogInput(rsdata) // Обработка аналоговых входов
	ms.Inputs = make(map[string]interface{})

	for _, di := range DiskrInputs {
		ms.Inputs[di.Name] = di.Val
	}
	for _, ai := range AnalogInputs {
		ms.Inputs[ai.Name] = ai.Val
	}
	// актуальные маркирововочные данные
	telNum := ms.Inputs["val_tel_num"].(int)
	for _, mdata := range MTelBuff {
		if mdata.Num == uint16(telNum) {
			ActualMarkData = mdata
		}
	}
	ms.ActPlate = ActualMarkData.Plate
	return ms
}

func NewMF(mfdata []byte) str.MF {
	num := uint16(mfdata[3]) + uint16(mfdata[2])*256

	Id := 0
	Krat := 0
	for _, val := range MTelBuff {
		if val.Num == num {
			Id = val.Plate.SlabId
			Krat = val.Plate.Krat
		}
	}
	return str.MF{0, Id, Krat, int(mfdata[5]), int(mfdata[7]), time.Now(), int(num)}
}

func GetMashineStatus() (*str.MS, error) {
	sendBuff := make([]byte, SendLength)
	buff := bytes.NewBuffer(sendBuff)
	buff.Reset()
	buff.WriteString("RS")
	recvBuff := make([]byte, RecvLength)
	conn, err := getConnection()
	defer conn.Close()
	if err != nil {
		Connected = false
		return nil, err
	}
	_, err = conn.Write(sendBuff)
	if err != nil {
		Connected = false
		return nil, err
	}
	_, err = conn.Read(recvBuff)
	if err != nil {
		Connected = false
		return nil, err
	}
	ms := NewMS(recvBuff)
	return &ms, nil
}

/**
Создание маркировочной телеграммы.
*/

func createMarkTelegramm(mdata str.DataForMark) []byte {
	plate := mdata.Plate // данные по листу
	mt := mdata.Tm       // шаблон маркировки
	pt := mdata.Tp       // шаблон клеймовки
	sendBuff := make([]byte, SendLength)
	buff := bytes.NewBuffer(sendBuff)
	buff.Reset()
	buff.WriteString("MD") // MD
	buff.WriteByte(byte(mdata.Num & 0xff00 >> 8))
	buff.WriteByte(byte(mdata.Num & 0xff))
	//----------------------------------------
	var length Int = Int(plate.Length) // Длина
	buff.Write(length.getBytes())
	var width Short = Short(plate.Width) // Ширина
	buff.Write(width.getBytes())
	var t Short = Short(40) // Температура
	buff.Write(t.getBytes())
	var child Short = 1 // Дочерние
	buff.Write(child.getBytes())
	var punchSel Short = 1 // Выбор клеймовки
	buff.Write(punchSel.getBytes())
	var markSel Short = 1 // Выбор маркировки
	buff.Write(markSel.getBytes())
	//----------------------------------------
	var drawAround Short = 1 // Рисовать вокруг клейма
	if !mdata.DrawAround {
		drawAround = 0
	}
	buff.Write(drawAround.getBytes())
	//----------------------------------------
	var reverseText Short = 1 // Перевернуть текст
	if !mdata.RotateText {
		reverseText = 0
	}
	buff.Write(reverseText.getBytes())
	//----------------------------------------
	var spare1 Short = 1
	buff.Write(spare1.getBytes())
	var spare2 Int = 1
	buff.Write(spare2.getBytes())
	//---------------------------------------
	mp := str.CreateMap(plate)
	var markPos Short = 20
	buff.Write(markPos.getBytes())
	// if !Reverse {
	// 	buff.Write(handleS7(S7Str(mt.Line1), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line2), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line3), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line4), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line5), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line6), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line7), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line8), mp))
	// 	buff.Write(handleS7(S7Str(mt.Line9), mp))
	// } else {
	buff.Write(handleS7(S7Str(mt.Line9), mp))
	buff.Write(handleS7(S7Str(mt.Line8), mp))
	buff.Write(handleS7(S7Str(mt.Line7), mp))
	buff.Write(handleS7(S7Str(mt.Line6), mp))
	buff.Write(handleS7(S7Str(mt.Line5), mp))
	buff.Write(handleS7(S7Str(mt.Line4), mp))
	buff.Write(handleS7(S7Str(mt.Line3), mp))
	buff.Write(handleS7(S7Str(mt.Line2), mp))
	buff.Write(handleS7(S7Str(mt.Line1), mp))
	//}
	empty := make([]byte, 2326)
	buff.Write(empty)
	buff.Write(handleS7(S7Str(pt.Line10), mp))
	buff.Write(handleS7(S7Str(pt.Line11), mp))
	buff.Write(handleS7(S7Str(pt.Line12), mp))
	empty = make([]byte, 4)
	buff.Write(empty)
	var deep Short = Short(mdata.Deep) // Глубина
	buff.Write(deep.getBytes())
	return sendBuff
}

/*
Обработка строки для передачи в контроллер.
*/
func handleS7(line S7Str, mp map[string]string) []byte {
	for key, value := range mp {
		line = S7Str(strings.Replace(string(line), key, value, 5))
	}
	sdata := models.HandleSymb(models.ToCp1251(string(line)))
	out := make([]byte, 28)
	out[0] = 28
	out[1] = byte(len(sdata))
	copy(out[2:], sdata)
	return out
}

//-------------------------------------------------------------------------------
//-------------------------------------------------------------------------------

/* Посылка данных*/
func send(conn net.TCPConn, data []byte) error {
	done := make(chan error)
	go func() {
		if _, err := conn.Write(data); err != nil {
			done <- err
		}
		done <- nil
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(time.Millisecond * 100):
		return errors.New("Timeout send")
	}
	return nil
}

/* Подготовка телеграммы запроса статуса*/
func newMSReq() []byte {
	sendBuff := make([]byte, SendLength)
	buff := bytes.NewBuffer(sendBuff)
	buff.Reset()
	buff.WriteString("RS")
	if FlagRecvCommand {
		fmt.Println("Send to controller")
		buff.WriteByte(NumCommand)
		FlagRecvCommand = false
	}
	return sendBuff
}

/*S7Manager Управление соединением S7*/
func S7Manager() {
	var conn net.TCPConn // соединение по TCP
	//------Каналы------------------------
	inch := make(chan []byte)
	//--------------------------------------
	timer := time.Tick(1 * time.Second) // таймер 1 раз в сек
	var actualDataMS str.MS             // актуальная телеграмма статуса маркировочной машины
	var actualDataMF str.MF             // последняя телеграммы результата маркировки
	var flagReconnect bool
	/*Переподключение*/
	reconnect := func(ch chan str.CommandAction) {
		for !Connected {
			cn, err := getConnection()
			if err == nil {

				Connected = true
				conn = *cn
			} else {

				time.Sleep(500 * 1000)

			}
		}
		conn.SetKeepAlivePeriod(200 * time.Millisecond)
		//errors = make(chan error, 2)
		ch <- str.Connect
	}
	/*Обработка ошибок*/
	errHandler := func(err error) {
		fmt.Println("Error:", err)
		models.AddLog(str.S7, "Ошибка соединения: "+fmt.Sprintf("%s", err))
		Connected = false            // сброс флага соединения
		conn.Close()                 // закрываем соединение
		CommandChan <- str.Reconnect // переподключение
	}
	/*
		Поток получения данных из контролллера.
	*/
	getMessages := func(ch chan<- []byte) {
		run := true
		for run {
			buff := make([]byte, 60)
			_, err := conn.Read(buff)
			if err != nil {
				errHandler(err)
				run = false
				break
			} else {
				ch <- buff
			}
		}
		fmt.Println("End Get Messages")
	}

	/*-------Подключение----------*/
	go reconnect(CommandChan)
	for {
		select {
		case mdata := <-ChanMark: // посылка маркировочной телеграммы
			if Connected || true {
				mdata.Num = MTNumber
				tb := time.Now()
				if err := send(conn, createMarkTelegramm(mdata)); err != nil {
					errHandler(err)
				} else {
					MTNumber++
					models.SaveTemp("MarkTelCount", fmt.Sprintf("%d", MTNumber))
					if MTNumber > 60000 {
						MTNumber = 0
					}
					MTelBuff = append(MTelBuff, mdata)
					if len(MTelBuff) > 10 {
						MTelBuff = MTelBuff[1:]
					}
				}
				fmt.Println("Time Send", time.Since(tb))
			}
		case command := <-CommandChan: // Выполнение команд
			switch {
			case command == str.GetMS: // получение MS
				actualDataMS.ControllerConnect = Connected
				actualDataMS.DBConnect = models.DBConnected
				RSData <- actualDataMS
			case command == str.GetMF: // получение MF
				actualDataMS.ControllerConnect = Connected
				actualDataMS.DBConnect = models.DBConnected
				MFData <- actualDataMF
			case command == str.Reconnect: // Перезапуск
				if !flagReconnect {
					fmt.Println("Reconnect")
					flagReconnect = true
					go reconnect(CommandChan)
				}
			case command == str.Connect: // соединение
				fmt.Println("Connect")
				flagReconnect = false
				models.AddLog(str.S7, "Соединение с контроллером установлено")
				go getMessages(inch)
				Connected = true
			}
		case gDat := <-inch: // прием входных данных
			if len(gDat) == 60 && gDat[0] == 'M' && gDat[1] == 'S' { // телеграмма статуса MS
				actualDataMS = NewMS(gDat)
				models.AddMS(actualDataMS)
			}
			if len(gDat) == 60 && gDat[0] == 'M' && gDat[1] == 'F' { // телеграмма маркировки MF
				actualDataMF = NewMF(gDat)
				fmt.Println("MF", actualDataMF)
				models.AddMF(actualDataMF)
			}
		case <-timer: // посылка запроса статуса MS по таймеру.
			if Connected { // если есть соединение
				if err := send(conn, newMSReq()); err != nil {
					errHandler(err)
				}
			}
		}
	}

}
