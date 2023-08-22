const { createApp, ref } = Vue;
const app = createApp({
    data() {
        return {
            title: "",
            addWebSite: {
                api: "/api/website/add",
                param : {
                    host: "",
                    healthUri: "",
                    rate: 10,
                    alarmResTime: 3000,
                    uriDepth: 2,
                    uriUpdateRate: 24,
                }
            },
            msg: "",
            websiteList: {
                api: function (){
                    return "/api/website/list?pg="+this.page;
                },
                page: 1,
                data: [],
            },
            deleteWebsite: {
                api: function (){
                    return "/api/website/delete/"+id;
                },
                id: "",
                hostName: "",
            },
            hasMail: {
                api: "/api/mail/init",
                data: {},
            },
            mailConf: {
                api: "/api/mail/conf",
                param: {
                    host: "smtp.qq.com",
                    port: 25,
                    from: "",
                    authCode: "",
                    toList: "",
                },
            },
            mailInfo: {
                api: "/api/mail/info",
            },
            mailSend: {
                api: "/api/mail/sendTest",
            },
            point: {
                hostId: "",
                hostUri: "",
                apiAdd: function (){
                    return "/api/point/add/"+this.hostId;
                },
                apiList: function () {
                    return "/api/point/list/"+this.hostId;
                },
                apiDel: function () {
                    return "/api/point/del/"+this.hostId;
                },
                param: {
                    uri:"",
                },
                uriList: [],
                nowUri: "",
            },
            websiteInfo: {
                hostId: "",
                api: function() {
                    return "/api/website/info/"+this.hostId;
                },
                data: {}
            },
            alertList: {
                api: "/api/alert/list",
                clear: "/api/alert/clear",
                list: [],
                len: 0,
            },
            monitorErrList: {
                api: "/api/monitor/err/list",
                clear: "/api/monitor/err/clear",
                list: [],
                len: 0,
            },
            isOk: "",
            monitorLog: {
                hostId: "",
                api: function (){
                    return "/api/monitor/log/" + this.hostId;
                },
                data: {},
            },
            editWebsiteConf: {
                api: "/api/website/edit",
                param: {
                    hostId: "",
                    rate: 10,
                    alarmResTime: 3000,
                    uriDepth: 2,
                },
                host: "",
            }
        }
    },
    created:function(){
        let t = this;
        t.getList();
        t.getMail();
        t.getMailInfo();
        t.getAlertList();
        t.getMonitorErrList();
        t.$nextTick(() => {
            t.DrawChart();
            }
        );
        t.timer = window.setInterval(() => {
            t.getList();
        }, 10000);
    },
    destroyed:function () {
        let t = this;
        window.clearInterval(t.timer);
    },
    methods: {
        toastShow: function (msg) {
            let t = this;
            t.msg = msg;
            const toastLiveExample = document.getElementById('liveToast')
            const toast = new bootstrap.Toast(toastLiveExample)
            toast.show()
        },
        get: function (url, func) {
            $.ajax({
                type: "get",
                url: url,
                data: "",
                dataType: 'json',
                success: function(data){
                    func(data);
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },
        post: function (url, param, func) {
            $.ajax({
                type: "post",
                url: url,
                data: JSON.stringify(param),
                dataType: 'json',
                success: function(data){
                    func(data);
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },
        addWebSiteMonitor: function () {
            let t = this;
            t.addWebSite.param.rate = Number(t.addWebSite.param.rate);
            t.addWebSite.param.alarmResTime = Number(t.addWebSite.param.alarmResTime);
            t.addWebSite.param.uriDepth = Number(t.addWebSite.param.uriDepth);
            t.post(t.addWebSite.api, t.addWebSite.param, function (data){
                if (data.code === 0) {
                    $("#addHostModal").modal('toggle');
                    t.getList();
                }
                t.toastShow(data.msg);
            });
        },
        test:function () {
            console.log("test...")
        },
        getList: function () {
            let t = this;
            t.get(t.websiteList.api(), function (data){
                t.websiteList.data = data.data;
            });
        },
        getAlertList: function () {
            let t = this;
            t.get(t.alertList.api, function (data) {
                t.alertList.list = data.data;
                t.alertList.len = t.alertList.list.length;
            });
        },
        alertClear: function () {
            let t = this;
            t.isOk = "alertClear";
            $("#isOkModal").modal("show");
        },
        alertClearSubmit: function () {
            let t = this;
            t.get(t.alertList.clear, function (data){
                t.toastShow(data.msg);
                t.getAlertList();
                $("#isOkModal").modal('toggle');
            });
        },
        getMonitorErrList: function () {
            let t = this;
            t.get(t.monitorErrList.api, function (data){
                t.monitorErrList.list = data.data;
                t.monitorErrList.len = t.monitorErrList.list.length;
            });
        },
        monitorErrClear: function () {
            let t = this;
            t.isOk = "monitorErrClear";
            $("#isOkModal").modal("show");
        },
        monitorErrClearSubmit: function () {
            let t = this;
            t.get(t.monitorErrList.clear, function (data){
                t.toastShow(data.msg);
                t.getMonitorErrList();
                $("#isOkModal").modal('toggle');
            });
        },
        gotoList: function (pg) {
            let t = this;
            t.websiteList.page = pg;
            t.getList();
        },
        getMail: function () {
            let t = this;
            t.get(t.hasMail.api, function (data){
                t.hasMail.data = data.data;
            })
        },
        setMailConf: function () {
            let t = this;
            if (Array.isArray(t.mailConf.param.toList)) {
                t.mailConf.param.toList = t.mailConf.param.toList.join("");
            }
            t.post(t.mailConf.api, t.mailConf.param, function (data){
                if (data.code === 0) {
                    $("#mailSetModal").modal('toggle');
                    t.getMail();
                }
                t.toastShow(data.msg);
            });
        },
        getMailInfo: function () {
            let t = this;
            t.get(t.mailInfo.api, function (data){
                t.mailConf.param.host = data.data.Host;
                t.mailConf.param.port = data.data.Port;
                t.mailConf.param.from = data.data.From;
                t.mailConf.param.authCode = data.data.AuthCode;
                t.mailConf.param.toList = data.data.ToList;
            });
        },
        mailSendTest: function () {
            let t = this;
            if (Array.isArray(t.mailConf.param.toList)) {
                t.mailConf.param.toList = t.mailConf.param.toList.join("");
            }
            t.mailConf.param.port = Number(t.mailConf.param.port);
            t.post(t.mailSend.api, t.mailConf.param, function (data){
                t.toastShow(data.msg);
            });
        },
        setUriPoint: function (item) {
            let t = this;
            t.point.hostUri = item.HealthUri + "/";
            t.point.hostId = item.ID;
            t.getUriPoint();
            $("#setUriModal").modal('show');
        },
        getUriPoint: function () {
            let t = this;
            t.get(t.point.apiList(), function (data){
                t.point.uriList = []
                if (data.code === 0) {
                    t.point.uriList = data.data;
                }
            });
        },
        addUriPoint: function () {
            let t = this;
            if (t.point.nowUri === "") {
                t.toastShow("请输入URI");
                return
            }
            t.point.param.uri = t.point.hostUri + t.point.nowUri;
            t.post(t.point.apiAdd(), t.point.param, function (data){
                t.toastShow(data.msg);
                t.getUriPoint();
            });
        },
        gotoUriPoint: function (hostId, uri) {
            let t = this;
            t.point.hostId = hostId;
            t.point.param.uri = uri;
            t.post(t.point.apiAdd(), t.point.param, function (data){
                t.toastShow(data.msg);
                t.getUriPoint();
            });
        },
        delUriPoint: function (uri) {
            let t = this;
            t.point.param.uri = uri;
            t.post(t.point.apiDel(), t.point.param, function (data){
                t.toastShow(data.msg);
                t.getUriPoint();
            });
        },
        openWebsiteInfo: function (id) {
            let t = this;
            t.websiteInfo.hostId = id;
            t.get(t.websiteInfo.api(), function (data){
                if (data.code === 0) {
                    t.websiteInfo.data = data.data;
                    $("#websiteInfoModal").modal('show');
                } else {
                    t.toastShow(data.msg);
                }
            });
        },
        logShow: function (id) {
            let t = this;
            t.monitorLog.hostId = id;
            t.get(t.monitorLog.api(), function (data){
                if (data.code === 0) {
                    t.monitorLog.data = data.data;
                    $("#logModal").modal('show');
                } else {
                    t.toastShow(data.msg);
                }
            });
        },
        deleteWebsiteOpen: function (item) {
            let t = this;
            t.deleteWebsite.hostId = item.ID;
            t.deleteWebsite.hostName = item.Host;
            t.isOk = "deleteWebsite";
            $("#isOkModal").modal("show");
        },
        deleteWebsiteSubmit: function () {
            let t = this;
            t.get(t.deleteWebsite.api(), function (data){
                t.toastShow(data.msg);
                t.getMonitorErrList();
                $("#isOkModal").modal('toggle');
                t.getList();
            });
        },
        openEditWebsiteConf: function (item) {
            let t = this;
            t.editWebsiteConf.host = item.Host;
            t.editWebsiteConf.param.hostId = item.ID;
            t.editWebsiteConf.param.rate = item.Rate;
            t.editWebsiteConf.param.alarmResTime = item.AlarmResTime;
            t.editWebsiteConf.param.uriDepth = item.UriDepth;
            console.log(item)
            console.log(t.editWebsiteConf)
            $("#setAlertModal").modal("show");
        },
        editWebsiteConfSubmit: function () {
            let t = this;
            t.editWebsiteConf.param.rate = Number(t.editWebsiteConf.param.rate);
            t.editWebsiteConf.param.alarmResTime = Number(t.editWebsiteConf.param.alarmResTime);
            t.editWebsiteConf.param.uriDepth = Number(t.editWebsiteConf.param.uriDepth);
            t.post(t.editWebsiteConf.api, t.editWebsiteConf.param, function (data){
                t.toastShow(data.msg);
                $("#setAlertModal").modal('toggle');
            })
        },
        copy: function () {
            let t = this;
            let clipboard = new ClipboardJS('.copy');
            clipboard.on('success', e => {
                t.toastShow("复制成功!");
                e.clearSelection();
            });
            clipboard.on('error', e => {
                t.toastShow("复制失败！请重试或者手动复制内容!");
            });
        },
        openChart: function (){
            $("#chartModal").modal("show");
        },
        DrawChart: function () {
            // let base = +new Date(1988, 9, 3);
            //
            // let oneDay = 24 * 3600 * 1000;
            // let data = [[base, Math.random() * 300]];
            // for (let i = 1; i < 20000; i++) {
            //     let now = new Date((base += oneDay));
            //     data.push([+now, Math.round((Math.random() - 0.5) * 20 + data[i - 1][1])]);
            // }
            //
            // console.log(data)

            let data = [[619632000000, 125],[619718400000, 245],[619804800000, 268],[619891200000, 133],]

            option = {
                tooltip: {
                    trigger: 'axis',
                    position: function (pt) {
                        return [pt[0], '10%'];
                    }
                },
                title: {
                    left: 'center',
                    text: '站点: www.33633.cn'
                },
                toolbox: {
                    feature: {
                        dataZoom: {
                            yAxisIndex: 'none'
                        },
                        restore: {},
                        saveAsImage: {}
                    }
                },
                xAxis: {
                    type: 'time',
                    boundaryGap: false
                },
                yAxis: {
                    type: 'value',
                    boundaryGap: [0, '100%']
                },
                dataZoom: [
                    {
                        type: 'inside',
                        start: 60,
                        end: 100
                    },
                    {
                        start: 60,
                        end: 100
                    }
                ],
                series: [
                    {
                        name: 'Fake Data',
                        type: 'line',
                        smooth: true,
                        symbol: 'none',
                        areaStyle: {},
                        data: data
                    }
                ]
            };
            let myChart = echarts.init(document.getElementById('myChart'), '', {
                renderer: 'canvas',
                useDirtyRect: false
            });
            // 使用刚指定的配置项和数据显示图表。
            myChart.setOption(option);
        },
    },
    computed: {
    },
    mounted:function(){
    },
});
app.mount('#app');