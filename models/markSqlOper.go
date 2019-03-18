package models

import (
	"benew/str"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

const connStr = "server=192.168.213.130;user id=asutp_tlc2;database=TLC2;password=djctvm tlbybw;encrypt=disable;port=1433"

var (
	platesMem []str.PlateData
	/*DBConnected - флаг соединения с базой данных*/
	DBConnected = false
	/*Commands - канал поступления команд*/
	Commands = make(chan str.CommandAction)
	errchan  = make(chan error)
	/*Potok - номер потока*/
	Potok = 1
)

func main() {
	fmt.Println("Go")

	fmt.Println(GetMarkTemplateBId(1252616, 1))
}

/*
Получение списка шаблонов маркировки
*/
func GetMarkTemplate(potok int) []TemplateMark {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]TemplateMark, 0)
	}
	defer db.Close()
	rows, err := db.Query("SELECT  IDTemplate ,TemplateName,isnull(Line1,''),isnull(Line2,''),isnull(Line3,''),isnull(Line4,''),"+
		"isnull(Line5,''),isnull(Line6,''),isnull(Line7,''),isnull(Line8,''),isnull(Line9,''),[potok] FROM [TLC2].[Stamp].[tb_Template] where potok=? order by TemplateName", potok)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]TemplateMark, 0)
	}
	templates := make([]TemplateMark, 0, 200)
	for rows.Next() {
		templ := TemplateMark{}
		rows.Scan(&templ.IDTemplate, &templ.TemplateName, &templ.Line1, &templ.Line2, &templ.Line3, &templ.Line4, &templ.Line5, &templ.Line6, &templ.Line7, &templ.Line8, &templ.Line9, &templ.Potok)
		templates = append(templates, templ)
	}
	defer rows.Close()
	return templates
}

/*GetPunchTemplate -Получение списка шаблонов клеймовки
 */
func GetPunchTemplate(potok int) []TemplatePunch {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]TemplatePunch, 0)
	}
	defer db.Close()
	rows, err := db.Query("SELECT  IDTemplate1 ,TemplateName,isnull(Line10,''),isnull(Line11,''),isnull(Line12,''),[potok] FROM [TLC2].[Stamp].[tb_Template1] where potok=? order by TemplateName", potok)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]TemplatePunch, 0)
	}
	templates := make([]TemplatePunch, 0, 200)
	for rows.Next() {
		templ := TemplatePunch{}
		rows.Scan(&templ.IDTemplate, &templ.TemplateName, &templ.Line10, &templ.Line11, &templ.Line12, &templ.Potok)
		templates = append(templates, templ)
	}
	defer rows.Close()
	return templates
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
Шаблон клеймовки
*/
type TemplatePunch struct {
	IDTemplate             int
	TemplateName           string
	Line10, Line11, Line12 string
	Potok                  int
}

/*
Получение примечаний по фабрикационному номеру и номеру заказа
*/
func GetRemarks(rollFabId, orderId int) Remarks {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return Remarks{}
	}
	defer db.Close()
	rows, err := db.Query("exec Stamp.sp_GetRemarks ?,?", rollFabId, orderId)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return Remarks{}
	}
	rem := Remarks{}
	for rows.Next() {
		rows.Scan(&rem.Remark, &rem.AddCond, &rem.Marking, &rem.Stigma)
	}
	defer rows.Close()
	return rem
}

/* Дополнительные данные */
type Remarks struct {
	Remark, AddCond, Marking, Stigma string
}

func GetRemarksId(id, krat int) Remarks {
	for _, plate := range platesMem {
		if int(plate.SlabId) == id && int(plate.Krat) == krat {
			return GetRemarks(int(plate.RollFabId), int(plate.OrderId))
		}
	}
	return Remarks{}
}

/*
Получение номера шаблона маркировки
*/
func getIdTemplateMark(name string) int {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return -1
	}
	defer db.Close()
	rows, err := db.Query("select IDTemplate from Stamp.tb_Template where TemplateName like ?", name)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return -1
	}
	num := -1
	for rows.Next() {
		rows.Scan(&num)
	}
	defer rows.Close()
	return num
}

/*
Получение данных о штуке.
*/

func ReturnTableData(firstDT string, all bool, accessRole bool, show_id bool, slab_id string) []str.PlateData {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]str.PlateData, 0)
	}
	defer db.Close()
	allI := 1
	if !all {
		allI = 0
	}
	show_idI := 1
	if !show_id {
		show_idI = 0
	}
	var rows *sql.Rows
	//Stamp.sp_GetPost3Slabs_right -  left  левая right  правая
	switch Potok {
	case 1:
		rows, err = db.Query("exec Stamp.sp_GetPost3Slabs_right ?,?,false,?,?", firstDT, allI, show_idI, slab_id)
	case 2:
		rows, err = db.Query("exec Stamp.sp_GetPost3Slabs_left ?,?,false,?,?", firstDT, allI, show_idI, slab_id)
	case 3:
		rows, err = db.Query("exec Stamp.sp_GetPost3Slabs_temp ?,?,false,?,?", firstDT, allI, show_idI, slab_id)
	}
	//rows, err := db.Query("exec Stamp.sp_GetPost3Slabs_right ?,?,false,?,?", firstDT, allI, show_idI, slab_id)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]str.PlateData, 0)
	}
	plates := make([]str.PlateDataNull, 0, 200)
	for rows.Next() {
		plate := str.PlateDataNull{}
		err := rows.Scan(&plate.Heat, &plate.Grade, &plate.Height, &plate.Width, &plate.Length, &plate.FabNumber, &plate.FabPos, &plate.OrderNumber, &plate.LotNumber,
			&plate.OrderPos, &plate.IdFabPos, &plate.RollFabId, &plate.OrderId, &plate.IdOrderPos, &plate.SlabId, &plate.Krat, &plate.CNRS,
			&plate.PartyNum, &plate.NumInParty, &plate.SeqNum, &plate.M_pl, &plate.NewMark, &plate.Usl1, &plate.Usl2)
		plates = append(plates, plate)
		if err != nil {
			errchan <- err
		}
	}
	out := make([]str.PlateData, len(plates))
	for i, plate := range plates {
		out[i] = plate.Convert()
	}
	platesMem = out
	return out

}

/*
 Возврат сляба с заданным id
*/

func returnSlabs(id int) bool {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}
	defer db.Close()
	_, err = db.Exec("exec Stamp.sp_ReturnSlabs ?", id)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}

	return true
}

/*
Список замаркированных слябов.
*/

func getParts() []Plate {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]Plate, 0)
	}
	defer db.Close()
	rows, err := db.Query("exec Stamp.sp_GetSlabsMark")
	if err != nil {
		errchan <- err
		return make([]Plate, 0)
	}
	plates := make([]Plate, 0, 200)
	for rows.Next() {
		plate := Plate{}
		rows.Scan(&plate.SlabId, &plate.Krat, &plate.List)
		plates = append(plates, plate)
	}
	return plates
}

/*
Структура ID файла
*/
type Plate struct {
	SlabId, Krat, List int
}

/*
Получение шаблона маркировки по idOrderPos
*/

func GetMarkTemplateBId(orderId int, potok int) TemplateMark {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return TemplateMark{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT  IDTemplate ,TemplateName,isnull(Line1,''),isnull(Line2,''),isnull(Line3,''),isnull(Line4,''),"+
		"isnull(Line5,''),isnull(Line6,''),isnull(Line7,''),isnull(Line8,''),isnull(Line9,''),[potok] FROM [TLC2].[Stamp].[tb_Template] where potok=? "+
		"and  IDTemplate =(select top 1 IDTemplate from Stamp.tb_OrderTemplate where idOrderPos=?) order by TemplateName", potok, orderId)
	if err != nil {
		errchan <- err
	}
	templ := TemplateMark{}
	for rows.Next() {
		err := rows.Scan(&templ.IDTemplate, &templ.TemplateName, &templ.Line1, &templ.Line2, &templ.Line3, &templ.Line4,
			&templ.Line5, &templ.Line6, &templ.Line7, &templ.Line8, &templ.Line9, &templ.Potok)
		if err != nil {
			fmt.Println(err)
		}
	}
	return templ
}

/**
Получение шаблона клеймовки по idOrderPos
*/

func GetPunchTemplateById(orderId int, potok int) TemplatePunch {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return TemplatePunch{}
	}
	defer db.Close()
	rows, err := db.Query("select * from Stamp.tb_Template1 where potok=? "+
		"and  IDTemplate1 =(select top 1 IDTemplate1 from Stamp.tb_OrderTemplate where idOrderPos=?) order by TemplateName", potok, orderId)
	if err != nil {
		errchan <- err
		return TemplatePunch{}
	}
	templ := TemplatePunchNull{}
	for rows.Next() {
		err := rows.Scan(&templ.IDTemplate, &templ.TemplateName, &templ.Line10, &templ.Line11, &templ.Line12, &templ.Potok)
		if err != nil {
			errchan <- err
		}
	}
	return templ.convert()
}

type TemplatePunchNull struct {
	IDTemplate             sql.NullInt64
	TemplateName           sql.NullString
	Line10, Line11, Line12 sql.NullString
	Potok                  sql.NullInt64
}

func (tnp *TemplatePunchNull) convert() TemplatePunch {
	out := TemplatePunch{}
	if tnp.IDTemplate.Valid {
		out.IDTemplate = int(tnp.IDTemplate.Int64)
	}
	if tnp.TemplateName.Valid {
		out.TemplateName = tnp.TemplateName.String
	}
	if tnp.Line10.Valid {
		out.Line10 = tnp.Line10.String
	}
	if tnp.Line11.Valid {
		out.Line11 = tnp.Line11.String
	}
	if tnp.Line12.Valid {
		out.Line12 = tnp.Line12.String
	}
	if tnp.Potok.Valid {
		out.Potok = int(tnp.Potok.Int64)
	} else {
		out.Potok = 1
	}
	return out
}

/*
Удаление шаблона маркировки.
*/

func DeleteMarkTemplate(num int) bool {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}
	defer db.Close()
	res, err := db.Exec("delete  from Stamp.tb_Template where IDTemplate=?", num)
	if err != nil {
		errchan <- err
		return false
	}
	numR, err := res.RowsAffected()
	if err != nil {
		errchan <- err
		return false
	}
	if numR == 1 {
		return true
	}
	return false
}

/*
 Сохранение шаблона маркировки
*/

func SaveTemplateMark(tmpl TemplateMark) bool {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}
	defer db.Close()
	res, err := db.Exec("insert into Stamp.tb_Template values (?,?,?,?,?,?,?,?,?,?,?)", tmpl.TemplateName, tmpl.Line1, tmpl.Line2, tmpl.Line3, tmpl.Line4, tmpl.Line5,
		tmpl.Line6, tmpl.Line7, tmpl.Line8, tmpl.Line9, 1)
	if err != nil {
		errchan <- err
		return false
	}
	n, err := res.RowsAffected()
	if err != nil {
		errchan <- err
		return false
	}
	return n == 1

}

/*
Удаление шаблона клеймовки.
*/

func DeletePunchTemplate(num int) bool {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}
	defer db.Close()
	res, err := db.Exec("delete  from Stamp.tb_Template1 where IDTemplate1=?", num)
	if err != nil {
		errchan <- err
		return false
	}
	numR, err := res.RowsAffected()
	if err != nil {
		errchan <- err
		return false
	}
	if numR == 1 {
		return true
	}
	return false
}

/*
Сохранение шаблона клеймовки
*/

func SaveTemplatePunch(tmpl TemplatePunch) bool {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}
	defer db.Close()
	res, err := db.Exec("insert into Stamp.tb_Template1 values (?,?,?,?,?)", tmpl.TemplateName, tmpl.Line10, tmpl.Line11, tmpl.Line12, 1)
	if err != nil {
		errchan <- err
		return false
	}
	n, err := res.RowsAffected()
	if err != nil {
		errchan <- err
		return false
	}
	return n == 1

}

/*
Возврат списка замаркированных слябов
*/
func RetPlates() []str.RetPlate {
	plates := make([]str.RetPlate, 0, 20)
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return make([]str.RetPlate, 0)
	}
	defer db.Close()
	rows, err := db.Query("exec Stamp.sp_GetSlabsMark")
	if err != nil {
		errchan <- err
		return make([]str.RetPlate, 0)
	}
	for rows.Next() {
		plate := str.RetPlate{}
		err := rows.Scan(&plate.Id, &plate.Krat, &plate.List)
		if err != nil {
			errchan <- err
			return make([]str.RetPlate, 0)
		}
		plates = append(plates, plate)
	}
	return plates
}

/*
 Возврат замаркированного слябы по Id
*/

func RetPlate(id int) bool {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
		return false
	}
	defer db.Close()
	res, err := db.Exec("exec Stamp.sp_ReturnSlabs ?", id)
	if err != nil {
		errchan <- err
		return false
	}
	_, err = res.RowsAffected()
	if err != nil {
		errchan <- err
		return false
	}
	return true
}

/*GetPlate - Получение листа */
func GetPlate(id int, krat int) str.PlateData {
	fmt.Println(platesMem)
	for _, plate := range platesMem {
		if plate.SlabId == id && plate.Krat == krat {
			return plate
		}
	}
	return str.PlateData{}
}

/*TestConnection - Проверка соединения*/
func TestConnection(ch chan<- bool) {

	run := true
	for run {
		time.Sleep(time.Second)
		run = false
		db, err := sql.Open("mssql", connStr)
		if err != nil {
			run = true
			db.Close()
			continue
		}

		err = db.Ping()
		if err != nil {
			run = true
			db.Close()
			continue
		}
		db.Close()

	}
	ch <- true

}

/*DBManager - управление базой данных */
func DBManager() {
	chtest := make(chan bool)
	go TestConnection(chtest)
	timer := time.Tick(10 * time.Second)
	for {
		select {

		case <-chtest:
			DBConnected = true
			AddLog(str.DB, "Соединение с базой установлено.")
		case command := <-Commands:
			if command == str.Reconnect {
				DBConnected = false
				go TestConnection(chtest)
			}
		case error := <-errchan:
			AddLog(str.DB, fmt.Sprintf("Ошибка соединения: "+"%s", error))
			fmt.Println(error)
		case <-timer:
			if DBConnected {
				go isConnect()
			}
		}
	}
}

/* проверка соединения */
func isConnect() {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		errchan <- err
		Commands <- str.Reconnect
	}
}
