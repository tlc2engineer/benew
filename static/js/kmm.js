
$(document).ready(function(){
    var markTemplates; // Шаблоны маркировки
    var punchTemplates; // Шаблоны rktqvjdrb
    var dict; // Словарь
    var currMTemplate; // Текущий шаблон маркировки
    var Id; // Текущий номер
    var Krat; // Текущий крат
    var tc=0; // Счетчик
    var tableData;// Данные таблицы
    var deep=2;// Глубина клеймовки
    var connectController=true;
    var connectDB=true;
    function wrap(str,wr){
        return "<"+wr+">"+str+"</"+wr+">";
    }
    //displayMess("Тестовое сообщение: Чо там делает Корефан?")
    //------------------------Автодополнение----------------------------------------
    $(".inbtn").bind("click",onTextChange);
    function onTextChange() {
        //--Координаты области ввода------------------
        clip=$(this).parent().parent().find(".line").get(0).getBoundingClientRect();
        //--Область ввода-----------------------------
        txtinput=$(this).parent().parent().find(".line").get(0);
        $("#add_text").css("left",clip.left).css("top",clip.top).css("display","block");
        var parts=dict.split("\n");
        out="<table><tbody>";
        for(i=0;i<parts.length;i++){
            out+="<tr><td>"+parts[i]+"</td></tr>";
        }
        out+="</tbody></table>";
        // Установка списка слов
        $("#add_text_data").html(out);
        //-----Закрытие окна--------------------------------
        $("#add_text_close_btn").bind("click",function(){
            $("#add_text").css("display","none");
        });
        insert=function(text){
            old=txtinput.value
            var startPosition = txtinput.selectionStart;
            var endPosition = txtinput.selectionEnd;
            part1=old.slice(0,startPosition)
            part2=old.slice(endPosition,old.length)
            txtinput.value=part1+text+part2
            var event = new Event('change');
            txtinput.dispatchEvent(event);

        }
        //----------------------Двойной щелчок----------------
        $("#add_text_data>table>tbody>tr").bind("dblclick",function (){
            text=$(this).find("td").text();
            $("#add_text").css("display","none");
            insert(text)
        });
        //--------------Выделение элемента-------------------
        $("#add_text_data>table>tbody>tr").bind("click",function () {
            $("#add_text_data>table>tbody>tr").removeClass("selected");
            $(this).addClass("selected");
            text=$(this).find("td").html();
            $("#add_text_close_btn").bind("click",function(){
                $("#add_text").css("display","none");
                insert(text)
            });

        });

    }
    //------------------------Кнопка словарь----------------------------------------
    updateDict();
    function updateDict(){
        $.ajax({
            url: "/dict",
        }).done(function(data) {
            darea=$("#dict textarea").get(0)
            out=""
            for(i=0;i<data.length;i++){
                out+=data[i]+"\n";
            }
            $(darea).text(out)
            dict=out;

        });
    }
    $("#dict_btn").bind("click",function () {
        updateDict();
        $("#dict").css("display","block");
        $("#close_dict_btn").bind("click",function () {
            $("#dict").css("display","none");
        })
        $("#ok_dict_btn").bind("click",function () {
            text=$("#dict textarea").get(0).value;
            // console.log(this.getBoundingClientRect());
            $.ajax({
                url: "/dict",
                type: "POST",
                data: ({text:text})
            }).done(function(data) {
                $("#dict").css("display","none");


            });
        });

    })
    //-----------------------Функция обновления--------------------------------------
    function update(){
        slab_date=$("#slab_date").get(0).value
        allSlab=$("#all_slabs_cbx").get(0).value
        idCbx=$("#id_cbx").get(0).value
        slab_num=$("#slab_num").get(0).value
        $.ajax({
            url: "/table",
            data: ({date:slab_date,all:allSlab,idcbx:idCbx, slab_num:slab_num}),
            type: "GET"
        }).done(redrawMainTable);
    }
    //---------------Кнопка обновить---------------------------------------------------------------
    $("#upd_btn").bind("click",function(){
        update();
    });
    //---------------Возврат слябов---------------------------------------------------------------

    $("#upd_ret_w_btn").bind("click",function(){
        getRetList();
    });
    $("#ret_btn").bind("click",function(){
        $("#ret_slab").css("display","block");
        getRetList();
        // $("#ret_slab_btn>button")

    });

    function getRetList(){
        $.ajax({
            url: "/ret_plates",
        }).done(function (data) {
            out=""
            for(i=0;i<data.length;i++){
                out+="<tr>"
                out+="<td >"+data[i].Id+"</td>";
                out+="<td>"+data[i].Krat+"</td>";
                out+="<td class='listn'>"+data[i].List+"</td>";
                out+="</tr>"
            }
            $("#ret_slab_data>table>tbody").html(out);
            $("#ret_slab_data>table>tbody>tr").bind("click",function(){
                $("#ret_slab_data>table>tbody>tr").removeClass("selected");
                $(this).addClass("selected");
                id=$(this).find(".listn").html();
                console.log(id);
                $("#ret_ret_w_btn").bind("click",function(){
                    console.log("cl")
                    $.ajax({
                        url: "/ret_plate",
                        data: ({id:id}),
                    }).done(function(data){
                        if(data)
                            displayMess("Лист возвращен");
                        else
                            displayMess("Ошибка");
                    });
                });
            });
        })
    }
    $("#close_ret_w_btn").bind("click",function(){
        $("#ret_slab").css("display","none");
    });
    //---------------Маркировка-------------------------------------------------------------------
    $("#mark_templ_btn").bind("click",viewMTWindow);
    getMarkTemplList();
    getPunchTemplList();
    // --------Сохранение шаблона маркировки-------------------------------------------------------
    $("#MarkSaveBtn").bind("click",function(){
        tn=document.getElementById("templ_name").value
        lines=[]
        for(i=1;i<=9;i++){
            lines[i]=document.getElementById("mline"+i).value
        }

        $.ajax({
            url: "/save_m_template",
            data: ({name : tn,line1:lines[1],line2:lines[2],line3:lines[3],line4:lines[4],line5:lines[5],line6:lines[6],line7:lines[7],line8:lines[8],line9:lines[9]}),
            type: "POST"

        }).done(function (data) {
            if (data)
                displayMess("Шаблон маркировки"+tn+" добавлен.");
            else
                displayMess("Шаблон маркировки"+tn+" не добавлен,ошибка");
            getMarkTemplList();
        })
    });
    // Показ окна выбора шаблонов
    function viewMTWindow(){
        document.getElementById('win').style.display='block'
        setMTemplList();
        $("#templ_list").html(out);
        tname=$("#templ_list>tbody>tr:first").find("td").html();
        $("#templ_list>tbody>tr").bind("click",function () {
            $("#templ_list>tbody>tr").removeClass("selected");
            $(this).addClass("selected");
            tname=$(this).find("td").html();
            for(i=0;i<markTemplates.length;i++){
                if(markTemplates[i].TemplateName==tname){
                    out="<span>"+markTemplates[i].Line1+"</span><br>"
                    out+="<span>"+markTemplates[i].Line2+"</span><br>"
                    out+="<span>"+markTemplates[i].Line3+"</span><br>"
                    out+="<span>"+markTemplates[i].Line4+"</span><br>"
                    out+="<span>"+markTemplates[i].Line5+"</span><br>"
                    out+="<span>"+markTemplates[i].Line6+"</span><br>"
                    out+="<span>"+markTemplates[i].Line7+"</span><br>"
                    out+="<span>"+markTemplates[i].Line8+"</span><br>"
                    out+="<span>"+markTemplates[i].Line9+"</span><br>"
                    $("#right_content").html(out)
                    break;
                }
            }

        });
        $("#okMTmplBtn").bind("click",function () {
            console.log(tname);
            data=findMTemplateByName(tname);
            //console.log(data);
            setMarkTemplate(data,false);
            document.getElementById('win').style.display='none';
        })
        $("#delMTmplBtn").bind("click",function () {
            tmpl=findMTemplateByName(tname);
            $.ajax({
                url: "/delete_m_template",
                data: ({num : tmpl.IDTemplate}),
                type: "POST"

            }).done(function (data) {
                displayMess("Шаблон "+tname+" удален.");
                getMarkTemplList();
            })
        })
        $("#newMarkBtn").bind("click",function(){
            data=findMTemplateByName(tname);

            setMarkTemplate(data,true);
            document.getElementById('win').style.display='none';
        })

    }

    // --------------Клеймовка--------------------------------------------------------------------
    $("#punch_templ_btn").bind("click",viewPTWindow);
    function viewPTWindow(){
        document.getElementById('p_win').removeAttribute('style');
        setPTemplList();

    }
    $("#PunchSaveBtn").bind("click",function(){
        tn=document.getElementById("p_templ_name").value
        lines=[]
        for(i=10;i<=12;i++){
            lines[i]=document.getElementById("pline"+i).value
        }

        $.ajax({
            url: "/save_p_template",
            data: ({name : tn,line10:lines[10],line11:lines[11],line12:lines[12]}),
            type: "POST"

        }).done(function (data) {
            if(data)
                displayMess("Шаблон клеймовки "+tn+" добавлен ");
            else
                displayMess("Шаблон клеймовки "+tn+" не добавлен ошибка ");
            getPunchTemplList();
        })
    });
    //-----------------------------------------Обновление--------------------------------------------
    update();
    //------------------------------------------Перерисовка главной таблицы--------------------------
    function redrawMainTable(data) {
        tableData=data;
        ret=""
        dat=[];
        ghs=[]
        id_chked=$("#id_cbx").get(0).checked
        for(i=0;i<data.length;i++){
            dat[i]={id:data[i].SlabId,krat: data[i].Krat,idOrderPos: data[i].IdOrderPos}
            ghs[i]={heat:data[i].Heat,grade:data[i].Grade,height: data[i].Height,width: data[i].Width,length: data[i].Length}
            theSame=false
            if(!id_chked && ghs.length>1) {
                for (j = 0; j < ghs.length; j++) {
                    if (ghs[i] == ghs[j]) {
                        theSame = true
                    }
                }
            }
            if (theSame) continue;
            ret+="<tr>"
            ret+=wrap(data[i].Heat,"td");
            ret+=wrap(data[i].Grade,"td");
            ret+=wrap(data[i].Height,"td");
            ret+=wrap(data[i].Width,"td");
            ret+=wrap(data[i].Length,"td");
            ret+=wrap(data[i].FabNumber,"td");
            ret+=wrap(data[i].FabPos,"td");
            ret+=wrap(data[i].OrderNumber,"td");
            ret+=wrap(data[i].LotNumber,"td");
            ret+=wrap(data[i].OrderPos,"td");
            ret += wrap(data[i].SlabId, "td class=\"idcol\" ");
            ret += wrap(data[i].Krat, "td class=\"idkrat\"  ");
            ret += wrap(data[i].CNRS, "td class=\"idsx\"");
            ret += wrap(data[i].PartyNum, "td class=\"idsx\"");
            ret += wrap(data[i].SeqNum, "td class=\"idsx\"");
            ret+="</tr>"
        }


        if(dat.length>0){
            console.log("Add");
            Id=dat[0].id;
            Krat=dat[0].krat;
            addRemark(Id,Krat);
            idOrderPos=dat[0].idOrderPos;
            getMarkTemplate(idOrderPos);
            getPunchTemplate(idOrderPos);
        }
        $("#mbody").html(ret)
        if (id_chked){
            $(".idsx").removeClass("invis");
            $(".idcol").removeClass("invis");
            $(".idkrat").removeClass("invis");
        }
        else{
            $(".idsx").addClass("invis")
            $(".idcol").addClass("invis");
            $(".idkrat").addClass("invis");
        }
        //--------------------------------------------------
        $("#mbody>tr").bind("click",function (){

            $("#mbody>tr").removeClass("selected");
            $(this).addClass("selected")
            Id=$(this).find(".idcol").html();
            Krat=$(this).find(".idkrat").html();
            addRemark(Id,Krat);
            for(i=0;i<dat.length;i++){
                if(dat[i].id==Id && dat[i].krat==Krat){
                    getMarkTemplate(dat[i].idOrderPos);
                    getPunchTemplate(dat[i].idOrderPos);
                    break;
                }
            }

        });
        if (id_chked )
            $("#mbody>tr:first").addClass("selected");

    }
    //---------Привязка кнопки-------------------------------------------------------------------
    $("#send").bind("click",send);
    //-------------------------------------------------------------------------------------------
    $("#id_cbx").bind("change",function () {
        if (this.checked){
            $(".idsx").removeClass("invis");
            $(".idcol").removeClass("invis");
            $(".idkrat").removeClass("invis");
        }
        else{
            $(".idsx").addClass("invis")
            $(".idcol").addClass("invis");
            $(".idkrat").addClass("invis");
        }
        redrawMainTable(tableData)
    })
// --------------------------------------------Добавление пояснений----------------------------------
    function addRemark(id,krat){
        $.ajax({
            url: "/remark",
            data: ({SlabId : id,Krat:krat}),

        }).done(function (data) {
            $("#remarks").text(data.Remark)
            $("#add").text(data.AddCond)
            $("#marking").text(data.Marking)
            $("#stigma").text(data.Stigma)
        })
    }

    function getMarkTemplate(idOrderPos){
        $.ajax({
            url: "/mark_template",
            data: ({idOrderPos : idOrderPos}),

        }).done(function (data) {
            setMarkTemplate(data,false);
        })
    }
    function getPunchTemplate(idOrderPos){
        $.ajax({
            url: "/punch_template",
            data: ({idOrderPos : idOrderPos}),

        }).done(function (data) {
            setPunchTemplate(data,false);
        })
    }

    function getMarkTemplList(){
        $.ajax({
            url: "/mark_templates",
        }).done(function (data) {
            markTemplates=data;
            setMTemplList();
        })
    }
    function getPunchTemplList(){
        $.ajax({
            url: "/punch_templates",
        }).done(function (data) {
            punchTemplates=data;
            setPTemplList();
        })
    }
    /*
    Поиск шаблона маркировки по имени
    */
    function findMTemplateByName(name){
        for(i=0;i<markTemplates.length;i++){
            if(markTemplates[i].TemplateName==name){
                return markTemplates[i];
            }
        }
        return 0
    }
    /*
   Поиск шаблона клеймовки по имени
   */
    function findPTemplateByName(name){
        for(i=0;i<punchTemplates.length;i++){
            if(punchTemplates[i].TemplateName==name){
                return punchTemplates[i];
            }
        }
        return 0
    }
    /*
    Вывод шаблона маркировки
     */
    function setMarkTemplate(data,n){
        GetPaintImage(data);

        if(!n)
            $("#templ_name").get(0).value=data.TemplateName;
        else
            $("#templ_name").get(0).value="";
        $("#mline1").get(0).value=data.Line1
        $("#mline2").get(0).value=data.Line2
        $("#mline3").get(0).value=data.Line3
        $("#mline4").get(0).value=data.Line4
        $("#mline5").get(0).value=data.Line5
        $("#mline6").get(0).value=data.Line6
        $("#mline7").get(0).value=data.Line7
        $("#mline8").get(0).value=data.Line8
        $("#mline9").get(0).value=data.Line9
        $(".line").bind("change",function(){
            data.Line1=$("#mline1").get(0).value
            data.Line2=$("#mline2").get(0).value
            data.Line3=$("#mline3").get(0).value
            data.Line4=$("#mline4").get(0).value
            data.Line5=$("#mline5").get(0).value
            data.Line6=$("#mline6").get(0).value
            data.Line7=$("#mline7").get(0).value
            data.Line8=$("#mline8").get(0).value
            data.Line9=$("#mline9").get(0).value
            GetPaintImage(data);
        })

    }

    /*
Вывод шаблона клеймовки.
 */
    function setPunchTemplate(data,n){
        GetPanchImage(data);
        if (!n)
            $("#p_templ_name").get(0).value=data.TemplateName;
        else
            $("#p_templ_name").get(0).value="";
        $("#pline10").get(0).value=data.Line10
        $("#pline11").get(0).value=data.Line11
        $("#pline12").get(0).value=data.Line12
        $("#pline11,#pline10,#pline12").bind("change",function(){
            data.Line10=$("#pline10").get(0).value
            data.Line11=$("#pline11").get(0).value
            data.Line12=$("#pline12").get(0).value
            GetPanchImage(data)
        });
    }
    /*
     Установка списка шаблонов в окне выбора маркировки.
      */
    function setMTemplList(){
        out="<tbody>"
        for(i=0;i<markTemplates.length;i++){
            out+="<tr><td>"+markTemplates[i].TemplateName+"</td></tr>";
        }
        out+="</tbody>"
        $("#templ_list").html(out);
        $("#templ_list>tbody>tr").bind("click",function () {
            $("#templ_list>tbody>tr").removeClass("selected");
            $(this).addClass("selected");
            tname=$(this).find("td").html();
            for(i=0;i<markTemplates.length;i++){
                if(markTemplates[i].TemplateName==tname){
                    out="<span>"+markTemplates[i].Line1+"</span><br>"
                    out+="<span>"+markTemplates[i].Line2+"</span><br>"
                    out+="<span>"+markTemplates[i].Line3+"</span><br>"
                    out+="<span>"+markTemplates[i].Line4+"</span><br>"
                    out+="<span>"+markTemplates[i].Line5+"</span><br>"
                    out+="<span>"+markTemplates[i].Line6+"</span><br>"
                    out+="<span>"+markTemplates[i].Line7+"</span><br>"
                    out+="<span>"+markTemplates[i].Line8+"</span><br>"
                    out+="<span>"+markTemplates[i].Line9+"</span><br>"
                    $("#right_content").html(out)
                    break;
                }
            }

        });
    }
    /*
 Установка списка шаблонов в окне выбора клеймовки.
  */
    function setPTemplList() {
        out = "<tbody>"
        for (i = 0; i < punchTemplates.length; i++) {
            out += "<tr><td>" + punchTemplates[i].TemplateName + "</td></tr>";
        }
        out += "</tbody>"
        $("#p_templ_list").html(out);
        $("#p_templ_list>tbody>tr").bind("click", function () {
            $("#p_templ_list>tbody>tr").removeClass("selected");
            $(this).addClass("selected");
            tname = $(this).find("td").html();
            for (i = 0; i < punchTemplates.length; i++) {
                if (punchTemplates[i].TemplateName == tname) {
                    out = "<span>" + punchTemplates[i].Line10 + "</span><br>"
                    out += "<span>" + punchTemplates[i].Line11 + "</span><br>"
                    out += "<span>" + punchTemplates[i].Line12 + "</span><br>"
                    $("#p_right_content").html(out)
                    break;
                }
            }
            $("#p_okMTmplBtn").bind("click",function () {
                data=findPTemplateByName(tname);
                //console.log(data);
                setPunchTemplate(data,false);
                document.getElementById('p_win').style.display='none';
            })
            $("#delPTmplBtn").bind("click",function () {
                tmpl=findPTemplateByName(tname);
                $.ajax({
                    url: "/delete_p_template",
                    data: ({num : tmpl.IDTemplate}),
                    type: "POST"

                }).done(function (data) {
                    if(data)
                        displayMess("Шаблон "+tname+" удален.");
                    else
                        displayMess("Ошибка");
                    getPunchTemplList();
                })
            })
            $("#p_newMarkBtn").bind("click",function(){
                data=findPTemplateByName(tname);
                //console.log(data);
                setPunchTemplate(data,true);
                document.getElementById('p_win').style.display='none';
            })

        });

    }
//-------Обновление картинки маркировки
    function GetPaintImage(data) {
        tc+=1;
        src="/paint_image?line1="+data.Line1+"&line2="+data.Line2+"&line3="+data.Line3+"&line4="+data.Line4
            +"&line5="+data.Line5+"&line6="+data.Line6+"&line7="+data.Line7+"&line8="+data.Line8+"&line9="
            +data.Line9+"&id="+Id+"&krat="+Krat+"&tc="+tc;

        $("#markImage").html('<img src="'+src+'" '+'/>')
    }
//-------Обновление картинки клеймовки
    function GetPanchImage(data) {
        tc+=1;
        src="/panch_image?line10="+data.Line10+"&line11="+data.Line11+"&line12="+data.Line12+"&id="+Id+"&krat="+Krat+"&tc="+tc;
        $("#punchImage").html('<img src="'+src+'" '+'/>')
    }
    //--------Посылка данных по кнопке---------------------
    function send(){
        mtName=$("#templ_name").get(0).value
        ptName=$("#p_templ_name").get(0).value
        mt=findMTemplateByName(mtName)
        mt.Line1=$("#mline1").get(0).value
        mt.Line2=$("#mline2").get(0).value
        mt.Line3=$("#mline3").get(0).value
        mt.Line4=$("#mline4").get(0).value
        mt.Line5=$("#mline5").get(0).value
        mt.Line6=$("#mline6").get(0).value
        mt.Line7=$("#mline7").get(0).value
        mt.Line8=$("#mline8").get(0).value
        mt.Line9=$("#mline9").get(0).value
        console.log(mt.Line2);
        console.log(mt.Line3);
        pt=findPTemplateByName(ptName)
        pt.Line10=$("#pline10").get(0).value
        pt.Line11=$("#pline11").get(0).value
        pt.Line12=$("#pline12").get(0).value
        drawArownd=$("#drawArownd").get(0).checked
        rotateText=$("#rotateText").get(0).checked
        $.ajax({
            url: "/send_data",
            data: ({id : Id,krat: Krat,mt: mt,pt: pt,deep: deep,drawArownd: drawArownd,rotateText: rotateText}),
            type: "POST"

        }).done(function (data) {
            $("#send").css("background","yellow").css("color","white");
        })
        $("#send").css("background","white").css("color","black");
    }
    //----------Запрос статуса-------------------
    $("#conn_btn").bind("click",function(){
        $.ajax({
            url: "/rs",
        }).done(function (data) {
            console.log(data);
        })
    });
    //---------------------------------------------------------------------------------------
    $("#cbx_deep_light,#cbx_deep_middle,#cbx_deep_strong").bind("click",function(){
        lcbx=$("#cbx_deep_light").get(0)
        mcbx=$("#cbx_deep_middle").get(0)
        scbx=$("#cbx_deep_strong").get(0)
        id=$(this).attr("id")
        switch (id) {
            case  "cbx_deep_light":
                mcbx.checked=false;
                scbx.checked=false;
                deep=1;
                break;
            case   "cbx_deep_middle":
                lcbx.checked=false;
                scbx.checked=false;
                deep=2;
                break;
            case     "cbx_deep_strong":
                mcbx.checked=false;
                lcbx.checked=false;
                deep=3;
                break;
        }
    })
    //---------------------Соединение с контроллером и базой данных----------------------------------
    timer=setInterval(function () {
        $.ajax({
            url: "/rs",
        }).done(function (data) {
            if (!connectDB && data.DBConnect){
                update();
            }
            connectController=data.ControllerConnect
            connectDB=data.DBConnect;

            $("#conn_btn,#send").attr("disabled",!data.ControllerConnect)
            $("#ret_btn,#upd_btn,#mark_templ_btn,#punch_templ_btn,.inbtn,#PunchSaveBtn,#MarkSaveBtn").attr("disabled",!data.DBConnect)
            $("#contr_info").text(connectController ? "Контроллер:Есть":"Контроллер:Нет")
            if(connectController)
                $("#contr_info").css("color","blue");
            else
                $("#contr_info").css("color","red");
            $("#db_info").text(connectDB ? "База данныx:Есть":"База данныx:Нет")

            if(connectDB)
                $("#db_info").css("color","blue");
            else
                $("#db_info").css("color","red");
            dateTime=data.Time.split("T")
            date=dateTime[0]
            time=dateTime[1].split(".")[0]
            $("#last_mess_time_info").text("Последняя телеграмма: "+date+" "+time);
        })
    },2000);
    //----------Отображение сообщения------------------
    function displayMess(mess) {
        document.getElementById('alert_info').style.display='block';
        $("#alert_info span").text(mess)
    }
    $("#alert_close").bind("click",function(){
        document.getElementById('alert_info').style.display='none';
    });

});
