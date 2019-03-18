


var app = new Vue({

    el: '#main_table',
    data: {
        markResults: ["Маркировка не выбрана.","Маркировка не в порядке, маркировочный цикл не был начат.",
            "Маркировка не окончена. Маркировочный цикл прерван.","Маркировка в порядке."],
        file_data: "Нет данных",
        alarm_data: [],
        mark_data: [],
        error_list: [],
        auto: false,
        state: 1,
        timer: null
    },
    created: function () {
        console.log("Created")
        v=this
        axios.get('/error_list')
            .then(function (response) {
                v.error_list=response.data;

            })
            .catch(function (error) {
                console.log(error);
            });
        axios.get('/alarms')
            .then(function (response) {
                v.state = 1
                v.alarm_data=response.data
            })
            .catch(function (error) {
                console.log(error);
            });
    },
    methods: {
        show_ctrl: function (event) {
            v=this
            axios.get('/fileLog?type=s7')
                .then(function (response) {
                    v.state = 2
                    v.file_data = response.data
                })
                .catch(function (error) {
                    console.log(error);
                });
        },
        show_db: function (event) {
            v=this
            axios.get('/fileLog?type=db')
                .then(function (response) {
                    v.state = 3
                    v.file_data = response.data
                })
                .catch(function (error) {
                    console.log(error);
                });
        },
        show_alarm: function (event) {
            v=this
            axios.get('/alarms')
                .then(function (response) {
                    console.log(response)
                    v.state = 1
                    v.alarm_data=response.data
                })
                .catch(function (error) {
                    console.log(error);
                });

        },
        show_mark: function(){
            v=this
            axios.get('/mark_res_list')
                .then(function (response) {
                    console.log(response)
                    v.state = 4
                    v.mark_data=response.data
                })
                .catch(function (error) {
                    console.log(error);
                });
        },
        getErrText: function (n) {
            for (i = 0; i < this.error_list.length; i++) {
                if (this.error_list[i].Id == n+1) return this.error_list[i].Text;

            }
            return "";
        },
        formatTime:function(time){
            tdata=time.split("T")
            date=tdata[0]
            time=tdata[1].split(".")[0]
            return date+" "+time
        },
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
        updateAuto: function(){
            if(this.auto){
                v=this
                this.timer=setInterval(function (){
                    console.log("timer "+v.state);
                    switch (v.state){
                        case 1:
                            v.show_alarm(null);
                            break;
                        case 2:
                            v.show_ctrl(null);
                            break;
                        case 3:
                            v.show_db(null);
                            break;
                    }

                },2000);
            }else{
                clearInterval(this.timer);
            }
        },
        markText: function(n){
            return this.markResults[n]
        }

    }
})





