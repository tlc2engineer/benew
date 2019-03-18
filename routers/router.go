package routers

import (
	"benew/controllers"
	"github.com/astaxie/beego"
)

const Par = 1

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/kmm", &controllers.KmmController{})
	beego.Router("/table", &controllers.TableData{})
	beego.Router("/remark", &controllers.Remark{})
	beego.Router("/mark_template", &controllers.MarkTemplateById{})
	beego.Router("/punch_template", &controllers.PunchTemplateById{})
	beego.Router("/mark_templates", &controllers.ListMarkTamplates{})
	beego.Router("/punch_templates", &controllers.ListPunchTamplates{})
	beego.Router("/delete_m_template", &controllers.DeleteMTemplate{})
	beego.Router("/save_m_template", &controllers.SaveMTemplate{})
	beego.Router("/delete_p_template", &controllers.DeletePTemplate{})
	beego.Router("/save_p_template", &controllers.SavePTemplate{})
	beego.Router("/ret_plates", &controllers.RetListPlates{})
	beego.Router("/ret_plate", &controllers.RetPlate{})
	beego.Router("/dict", &controllers.DictController{})
	beego.Router("/paint_image", &controllers.PaintController{})
	beego.Router("/panch_image", &controllers.PunchController{})
	beego.Router("/send_data", &controllers.SendController{})
	beego.Router("/rs", &controllers.RS{})
	beego.Router("/fileLog", &controllers.DataFileLog{})
	beego.Router("/alarms", &controllers.Alarm{})
	beego.Router("/act_alarms", &controllers.ActiveAlarms{})
	beego.Router("/error_list", &controllers.ErrorList{})
	beego.Router("/mark_res_list", &controllers.MarkRL{})
	beego.Router("/concrete", &controllers.CAlarms{})
	beego.Router("/recv_ctrl", &controllers.RecvController{})
	beego.Router("/conf", &controllers.ConfController{})
}
