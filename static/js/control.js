var s_comp= {
    props: ['name','state'],
    template: ' <div   class="cycle" >\n' +
    '            <span >{{ name }}</span>\n' +
    '            <div>{{state}}</div>\n' +
    '        </div>'
}
Vue.component('indicator', {
    props: ['state'],
    template: '<img v-bind:src="addr"/>',
    computed: {
        addr: function () {
            switch (this.state){
                case 0://off
                    return "/static/img/red_s.png"

                case 1: //on
                    return "/static/img/green_s.png"

                case -1: //undefined
                    return "/static/img/gray_s.png"


            }
            return "/static/img/gray_s.png"
        }
        }
})
Vue.component('info-box', {
    props: ['name',"p1_name","p2_name","state1","state2"],
    template: ' <div class="data_box">\n' +
    '                <span>{{name}}</span>\n' +
    '                <div>{{p1_name}} <indicator v-bind:state="state1" /></div>\n' +
    '                <div>{{p2_name}} <indicator v-bind:state="state2" /></div>\n' +
    '            </div>',
})
Vue.component('data-field', {
    props: ['name',"state"],
    template: ' <div >\n' +
    '                <span>{{name}}</span>\n' +
    '                <indicator v-bind:state="state" />\n' +
    '            </div>',
})

Vue.component('ticktack-btn', {
    props: ['name','on_data'],
    template: "<button v-on:click='change' v-bind:class='{btn_on: on_data }'>{{name}}</button >",
    methods: {
        change: function(){
            this.$emit('update:on_data', !this.on_data)
            this.$emit('send_mess')
        }
    }
})


new Vue({
    el: '#main',
    data: {
        test_mode: false,
        Id: 1340330,
        krat: 1,
        width: 2500,
        length: 6000,
        thick: 6.0,
        state:1,
        paint_state:1,
        punch_state:1,
        rs_data: [],
        alarms:[],
        act_alarms: [],
        connect_ctrl: 1,
        connect_db: 1,
        machine_state: 1,
        mark_data_status: 0,
        mark_data_descr:["","","","","Нет данных","","Редактировать данные","","","","Проверка данных","","Ошибка","","Преобразование данных","","Данные готовы","","Маркировка","","Маркировка / Редактирование"],

        machine_states: ["Выкл","Ноль","Ожидает старта","","","","","","","Под давлением","","Техническое обслуживание","Удаленное управление",
            "","","","","","","","Авто","200 Тревога","Возврат домой","100 Тревога"],

        alarms_category:["Напряжение","Обмен данными Host","Давление воздуха / Минимальный уровень краски","Аварийная остановка","Ошибка при автоматическом цикле","Маркировочные данные",
        "Ошибка датчика","Ошибка безопасности","","","","","","Ошибка контрольной платы контура","Ошибка прерывателя цепи",
            "Ошибка позиционирования головки","Ошибка коммуникации с драйвером","Ошибка драйвера","",""],

        alarms_type:["Нет аварий","Предупреждение","Авария","Критическая авария"],

        state_descr: {
            0: "Не активно",
            1: "Перемещение в исходное положение",
            10: "Ожидание данных",
            13: "Проверка",
            21: "Ожидание листа",
            22: "Ожидание листа в положении остановки",
            23: "Запуск последовательности маркировки",
            33: "Маркировка",
            39: "Выбор режима работы",
            40: "Ожидание запуска последовательности маркировки",
            43: "Ожидание конца последовательности маркировки",
            55: "Ожидание запуска последовательности клеймовки",
            59: "Ожидание конца последовательности клеймовки",
            70: "Ожидание запуска последовательности клеймовки/маркировки",
            75: "Ожидание конца последовательности клеймовки/маркировки",
            77: "Ожидание конца последовательности клеймовки/маркировки",
            90: "Маркировка окончена",
        },
        paint_state_descr: {
            0: "Не активно",
            1: "Перемещение в исходное положение",
            4: "Ожидание данных",
            6: "Проверка выбрана ли краска",
            13: "Начало передачи данных драйверу",
            16: "Ожидание конца передачи данных",
            20: "Ожидание листа в позиции маркировки",
            29: "Передвижение головки в нижнее положение",
            33: "Ожидание запуска последовательности маркировки",
            36: "Ожидание окончания маркировки",
            41: "Есть ли еще маркировки",
            44: "Подъем головки в верхнее положение",
            53: "Ожидание окончания маркировки",
            59: "Отправка сигнала о завершении маркировки в главную последовательность."
        },
        punch_state_descr:{
            0: "Не активно",
            1: "Перемещение в исходное положение",
            4: "Ожидание данных",
            6: "Проверка выбрана ли клеймовка",
            10: "Ожидание листа",
            13: "Начало передачи данных драйверу",
            16: "Ожидание конца передачи данных",
            21: "Ожидание сигнала о начале клеймовки",
            22: "Ожидание блокировки рольганга",
            28: "Перемещение клеймовочной головки в нижнее положение",
            33: "Ожидание окончания клеймовки",
            40: "Перемещение головки в верхнее положение",
            53: "Ожидание окончания клеймовки",
            55: "Отправка сигнала о завершении маркировки в главную последовательность."


        },
        plate_present_sensor_sim: false, // Симуляция датчика наличия листа
        plate_contact_sensor_sim: false, // Симуляция датчика контакта с листом
        up_pos_sensor_sim: true, // Симуляция датчика верхней позиции марк головки
        punch_cont_sens_sim: false,  // Симуляция датчика контакта клеймовки с листом
        up_punch_pos_sim: false,  // Симуляция датчика верхней позиции клейовочной головки
        roll_on: true, // Работа рольганга
        plate_present_sensor: false, // Датчик наличия листа
        mark_work_sensor: false, // Датчик контакта маркировочной головки с листом
        mark_up_sensor: true, // Верхнее положение маркировочной головки
        punch_cont_sensor: false, // Контакт с листом клеймовочной головки
        up_punch_sensor: true, // Верхнее положение клеймовочнй головки
        trav_pos1: true, // Положение траверсы
        trav_pos2: false,
        timer: null,
        mark_begin: 210.0,
        list_pos: 2400.0,
        mark_wheel_move: false, //движение колес
        svgDoc:null,
        punchDown: false,
        mark_down: false,
        safety_emerg: false,
        safety_door: false,
        paint_level_ok: false,
        solvent_level_ok: false,
        stop_operator: false,
        stop_maint: false,
        press_prim: false,
        press_sec: false,
        auto_switch: false,
        maint_switch: false,
        paint_driver_ok: false,
        punch_driver_ok: false,
        conveyor: false,
        outrigg_mark: false,
        outrigg_maint: false,
        paint_head_up_pos: false,
        paint_head_work_pos: false,
        punch_up_pos: false,
        punch_plate_contact: false,
        sym_change: false, // флаг изменения симуляции
        kmm_view: 0
    },created: function () {
        v=this;

        this.timer=setInterval(function (){
            // Таймер
            axios.get('/rs')
                .then(function (response) {
                    rs=response.data
                    v.rs_data=response.data;
                    v.alarms=rs.Alarms;
                    v.act_alarms=rs.ActAlarms;
                    v.connect_ctrl=rs.ControllerConnect ? 1:0;
                    v.connect_db=rs.DBConnect ? 1:0;
                    v.machine_state=rs.MachineMode
                    v.mark_data_status=rs.DataStatus
                    v.state=rs.State
                    v.paint_state=rs.PaintState
                    v.punch_state=rs.PunchState
                    v.Id=rs["ActPlate"]["SlabId"]
                    v.krat=rs["ActPlate"]["Krat"]
                    v.width=rs["ActPlate"]["Width"]
                    v.length=rs["ActPlate"]["Length"]
                   v.thick=rs["ActPlate"]["Height"]
                    while(v.act_alarms.length<3){
                        nalarm={}
                        nalarm.Atype=0
                        nalarm.Num==3
                        v.act_alarms[v.act_alarms.length]=nalarm
                    }
                    v.handleInputs(rs.Inputs)
                })
                .catch(function (error) {
                    console.log(error);
                });
        },1000);

    },
    components:{
        'state-component': s_comp
    },
    methods:{
        getDate: function(time){
            tdata=time.split("T")
            date=tdata[0]
            return date
        },
        getTime: function(time){
            tdata=time.split("T")
            time=tdata[1].split(".")[0]
            return time
        },
        getSensorClass: function(on){
            return on ? "sens_on": "sens_off";
        },
        setSensorColor: function(el_name,val){
            this.svgDoc.getElementById(el_name).setAttribute("style",val ?
                "fill:#00ff00;fill-opacity:1;stroke:#000000;stroke-width:0.5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1":
                "fill:#ff0000;fill-opacity:1;stroke:#000000;stroke-width:0.5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1");
        },
        setPresent: function (el_name,val) {
            this.svgDoc.getElementById(el_name).setAttribute("visibility",val ? "visible":"hidden");
        },
        move_p_down: function(){
            if(this.punchDown){
                console.log("punch down")
            }
        },
        animation: function(){
            this.mark_down=!this.mark_down
            if(this.mark_down){
                pa=this.svgDoc.getElementById("plate_anim")
                pa.beginElement()
                pna=this.svgDoc.getElementById("punch_anim")
                pna.beginElement()
                ma=this.svgDoc.getElementById("mark_anim")
                ma.beginElement()
            }
        },
        handleInputs: function(inputs){
            if (this.test_mode) return
            this.list_pos=inputs["roll_pos"]
            this.safety_emerg = inputs["SAFETY_MODULE_EMERGENCY"]
            this.safety_door= inputs["SAFETY_MODULE_DOOR"]
            this.paint_level_ok= inputs["Pressure_Switch.PAINT_LEVEL_OK_COLD"]
            this.solvent_level_ok= inputs["Pressure_Switch.SOLVENT_LEVEL_OK"]
            this.stop_operator=inputs["EMERGENCY_STOP_OPERATOR"]
            this.stop_maint=inputs["EMERGENCY_STOP_MAINT"]
            this.press_prim=inputs["Pressure_Switch.AIR_PRESSURE_PRIMARY"]
            this.press_sec=inputs["Pressure_Switch.AIR_PRESSURE_SECONDARY"]
            this.auto_switch=inputs["KEY_SWITCH_IN_AUTO"]
            this.maint_switch=inputs["KEY_SWITCH_IN_MAINTENANC"]
            this.paint_driver_ok=inputs["S02.Internal_Signals.DRIVER_OK"]
            this.punch_driver_ok=inputs["S04.Internal_Signals.DRIVER_OK"]
            // Рольганг
            this.conveyor=inputs["CONVEYOR_STOPPED"]
            this.setPresent("roll_run",!this.conveyor)
            // Концевые траверсы
            this.outrigg_mark=inputs["OUTRIGGER_IN_MARKING"]
            this.outrigg_maint=inputs["OUTRIGGER_IN_MAINTENANCE"]
            this.setSensorColor("tr_pos_1",this.outrigg_mark)
            this.setSensorColor("tr_pos_2",this.outrigg_maint)
            this.paint_head_up_pos=inputs["HEAD_IN_UPPER_POS"]
            this.setSensorColor("mark_up_pos",this.paint_head_up_pos)
            //------------------------------------------------------
            this.paint_head_work_pos=inputs["HEAD_IN_STANDBY_POS"]
            this.setSensorColor("mark_work_pos",this.paint_head_work_pos)
            this.punch_up_pos=inputs["P01.Position_Detect.CYLINDER_1_UPPER"]
            this.setSensorColor("punch_head_pos",this.punch_up_pos)
            this.punch_plate_contact=inputs["P01.Position_Detect.PLATE_CONTACT_1"]
            this.setSensorColor("punch_head_contact",this.punch_plate_contact)
            if(this.sym_change) // если флаг изменения симуляции
                this.sym_change=false; // пропускаем 1 цикл
            else {
                this.plate_present_sensor_sim = inputs["Plate_Present_Sim"]
                this.plate_contact_sensor_sim = inputs["Mark_Contact_Sim"]
                this.up_pos_sensor_sim = inputs["Mark_Up_Pos_Sim"]
                this.punch_cont_sens_sim = inputs["Punch_Contact_Sim"]
            }
        },
        input_change: function(mess){
            val=false
            switch(mess){
                case 'on':
                    console.log("On")
                case 'off':
                case 'start':
                case 'stop':
                    break;
                case 'plate_on':
                    val=this.plate_present_sensor_sim;
                    break;
                case 'mark_contact':
                    val=this.plate_contact_sensor_sim;
                    break;
                case 'mark_up_pos':
                    val=this.up_pos_sensor_sim;
                    break;
                case 'punch_contact':
                    val=this.punch_cont_sens_sim;
                    break;
                case 'punch_up_pos':
                    val=this.up_punch_pos_sim;
                    break;
            }
            axios.post('/recv_ctrl?ms='+mess+'&val='+val)
                .then(function (response) {
                    console.log(response);
                })
                .catch(function (error) {
                    console.log(error);
                });
            this.sym_change=true
        }

    },

    mounted:function(){
        var a = document.getElementById("svgObject");
        v=this;
        a.addEventListener("load",function(){
            // Добавление анимаций svg
            v.svgDoc = a.contentDocument;
            var mh=v.svgDoc.getElementById("mark_head")
            var ph=v.svgDoc.getElementById("punch_head")
            var plate=v.svgDoc.getElementById("plate")
            // Функция создания анимации
            function createAnim(id,val,dur){
                var ani = document.createElementNS("http://www.w3.org/2000/svg","animateTransform");
                ani.setAttribute("type", "xml");
                ani.id=id
                ani.setAttribute("attributeName", "transform");
                ani.setAttribute("type", "translate");
                ani.setAttribute("values", val);
                ani.setAttribute("fill", "freeze");
                ani.setAttribute("begin", "indefenite");
                ani.setAttribute("additive", "sum");
                ani.setAttribute("dur", dur);
                return ani
            }
            ani=createAnim("mark_anim","0 0;0 0;0 55;0 55;0 55;0 55;0 55;0 55;0 55;0 0","10s")
            mh.appendChild(ani)
            var ani_p = createAnim("punch_anim","0 0;0 0;0 0;0 0;0 0;0 0;0 0;0 40;0 40;0 0","6s")
            ph.appendChild(ani_p)
            var ani_pl=createAnim("plate_anim","0 0;450 0; 450 0;900 0","12s")
            plate.appendChild(ani_pl)
        }, false);
    },

    watch: {
        roll_on: function(){
            if(this.test_mode)
                this.setPresent("roll_run",this.roll_on)

        },
        punch_cont_sensor:function(){
            this.setSensorColor("punch_head_contact",this.punch_cont_sensor)
        },
        up_punch_sensor: function () {
            this.setSensorColor("punch_head_pos",this.up_punch_sensor)
        },
        mark_up_sensor: function () {
            this.setSensorColor("mark_up_pos",this.mark_up_sensor)
        },
        mark_work_sensor: function () {
            this.setSensorColor("mark_work_pos",this.mark_work_sensor)
        },
        trav_pos1: function () {
            if(this.test_mode)
            this.setSensorColor("tr_pos_1",this.trav_pos1)
        },
        trav_pos2: function () {
            if(this.test_mode)
            this.setSensorColor("tr_pos_2",this.trav_pos2)
        },
        plate_present_sensor: function () {
            this.setPresent("plate",this.plate_present_sensor)
        },
        mark_wheel_move: function () {
            this.svgDoc.getElementById("mark_wheel1").setAttribute("style",this.mark_wheel_move ?
                "fill:#00ff00;fill-opacity:1;stroke:#000000;stroke-width:0.5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
                :"fill:#7f7f7f;fill-opacity:1;stroke:#000000;stroke-width:0.5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1")
            this.svgDoc.getElementById("mark_wheel2").setAttribute("style",this.mark_wheel_move ?
                "fill:#00ff00;fill-opacity:1;stroke:#000000;stroke-width:0.5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"
                :"fill:#7f7f7f;fill-opacity:1;stroke:#000000;stroke-width:0.5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1")
        },
        mark_data_status: function(){
            if(this.mark_data_status<16) this.kmm_view=0;
        }

    }
})
