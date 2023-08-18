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
                api: "/api/website/list",
                page: 1,
                data: [],
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
                apiAdd: "/api/point/add/",
                apiList: "/api/point/list/",
                apiDel: "/api/point/del/",
                param: {
                    uri:"",
                },
                uriList: [],
                nowUri: "",
            },
            websiteInfo: {
                hostId: "",
                api: "/api/website/info/",
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
                api: "/api/monitor/log/",
                data: {},
            },
        }
    },
    created:function(){
        var t = this;
        t.getList();
        t.getMail();
        t.getMailInfo();
        t.getAlertList();
        t.getMonitorErrList();
        t.timer = window.setInterval(() => {
            t.getList();
        }, 10000);
    },
    destroyed:function () {
        var t = this;
        window.clearInterval(t.timer)
    },
    methods: {

        toastShow: function (msg) {
            var t = this;
            t.msg = msg;
            const toastLiveExample = document.getElementById('liveToast')
            const toast = new bootstrap.Toast(toastLiveExample)
            toast.show()
        },

        addWebSiteMonitor: function () {
            var t = this;
            t.addWebSite.param.rate = Number(t.addWebSite.param.rate);
            t.addWebSite.param.alarmResTime = Number(t.addWebSite.param.alarmResTime);
            t.addWebSite.param.uriDepth = Number(t.addWebSite.param.uriDepth);
            $.ajax({
                type: "post",
                url: t.addWebSite.api,
                data: JSON.stringify(t.addWebSite.param),
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    if (data.code === 0) {
                        $("#addHostModal").modal('toggle');
                        t.getList();
                    }
                    t.toastShow(data.msg);
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        test:function () {
            console.log("test...")
        },

        getList: function () {
            var t = this;
            $.ajax({
                type: "get",
                url: t.websiteList.api+"?pg"+t.websiteList.page,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.websiteList.data = data.data;
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        getAlertList: function () {
            var t = this;
            $.ajax({
                type: "get",
                url: t.alertList.api,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.alertList.list = data.data;
                    t.alertList.len = t.alertList.list.length;
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        alertClear: function () {
            var t = this;
            t.isOk = "alertClear";
            $("#isOkModal").modal("show");
        },

        alertClearSubmit: function () {
            var t = this;
            console.log("alertClearSubmit")
            $.ajax({
                type: "get",
                url: t.alertList.clear,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.toastShow(data.msg);
                    t.getAlertList();
                    $("#isOkModal").modal('toggle');
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        getMonitorErrList: function () {
            var t = this;
            $.ajax({
                type: "get",
                url: t.monitorErrList.api,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.monitorErrList.list = data.data;
                    t.monitorErrList.len = t.monitorErrList.list.length;
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        monitorErrClear: function () {
            var t = this;
            t.isOk = "monitorErrClear";
            $("#isOkModal").modal("show");
        },

        monitorErrClearSubmit: function () {
            var t = this;
            console.log("monitorErrClearSubmit")
            $.ajax({
                type: "get",
                url: t.monitorErrList.clear,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.toastShow(data.msg);
                    t.getMonitorErrList();
                    $("#isOkModal").modal('toggle');
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        gotoList: function (pg) {
            var t = this;
            t.websiteList.page = pg;
            t.getList();
        },

        getMail: function () {
            var t = this;
            $.ajax({
                type: "get",
                url: t.hasMail.api,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.hasMail.data = data.data;
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        setMailConf: function () {
            var t = this;
            if (Array.isArray(t.mailConf.param.toList)) {
                t.mailConf.param.toList = t.mailConf.param.toList.join("");
            }
            $.ajax({
                type: "post",
                url: t.mailConf.api,
                data: JSON.stringify(t.mailConf.param),
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    if (data.code === 0) {
                        $("#mailSetModal").modal('toggle');
                        t.getMail();
                    }
                    t.toastShow(data.msg);
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        getMailInfo: function () {
            var t = this;
            $.ajax({
                type: "get",
                url: t.mailInfo.api,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.mailConf.param.host = data.data.Host;
                    t.mailConf.param.port = data.data.Port;
                    t.mailConf.param.from = data.data.From;
                    t.mailConf.param.authCode = data.data.AuthCode;
                    t.mailConf.param.toList = data.data.ToList;
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        mailSendTest: function () {
            var t = this;
            if (Array.isArray(t.mailConf.param.toList)) {
                t.mailConf.param.toList = t.mailConf.param.toList.join("");
            }
            t.mailConf.param.port = Number(t.mailConf.param.port);
            $.ajax({
                type: "post",
                url: t.mailSend.api,
                data: JSON.stringify(t.mailConf.param),
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.toastShow(data.msg);
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        setUriPoint: function (item) {
            var t = this;
            t.point.hostUri = item.HealthUri + "/";
            t.point.hostId = item.ID;
            console.log(t.point.hostId);
            t.getUriPoint();
            $("#setUriModal").modal('show');
        },

        getUriPoint: function () {
            var t = this;
            $.ajax({
                type: "get",
                url: t.point.apiList+t.point.hostId,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.point.uriList = []
                    if (data.code === 0) {
                        t.point.uriList = data.data;
                    }
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        addUriPoint: function () {
            var t = this;
            if (t.point.nowUri === "") {
                t.toastShow("请输入URI");
                return
            }
            t.point.param.uri = t.point.hostUri + t.point.nowUri;
            $.ajax({
                type: "post",
                url: t.point.apiAdd+t.point.hostId,
                data: JSON.stringify(t.point.param),
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.toastShow(data.msg);
                    t.getUriPoint();
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        gotoUriPoint: function (hostId, uri) {
            var t = this;
            t.point.hostId = hostId;
            t.point.param.uri = uri;
            $.ajax({
                type: "post",
                url: t.point.apiAdd+t.point.hostId,
                data: JSON.stringify(t.point.param),
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.toastShow(data.msg);
                    t.getUriPoint();
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        delUriPoint: function (uri) {
            var t = this;
            t.point.param.uri = uri;
            $.ajax({
                type: "post",
                url: t.point.apiDel+t.point.hostId,
                data: JSON.stringify(t.point.param),
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    t.toastShow(data.msg);
                    t.getUriPoint();
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        openWebsiteInfo: function (id) {
            var t = this;
            t.websiteInfo.hostId = id;
            $.ajax({
                type: "get",
                url: t.websiteInfo.api+t.websiteInfo.hostId,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    if (data.code === 0) {
                        t.websiteInfo.data = data.data;
                        $("#websiteInfoModal").modal('show');
                    } else {
                        t.toastShow(data.msg);
                    }
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        logShow: function (id) {
            var t = this;
            t.monitorLog.hostId = id;
            $.ajax({
                type: "get",
                url: t.monitorLog.api+t.monitorLog.hostId,
                data: "",
                dataType: 'json',
                success: function(data){
                    console.log(data)
                    if (data.code === 0) {
                        t.monitorLog.data = data.data;
                        $("#logModal").modal('show');
                    } else {
                        t.toastShow(data.msg);
                    }
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        },

        copy: function () {
            var t = this;
            var clipboard = new ClipboardJS('.copy');
            clipboard.on('success', e => {
                console.info('Action:', e.action);
                console.info('Text:', e.text);
                console.info('Trigger:', e.trigger);
                t.toastShow("复制成功!");
                e.clearSelection();
            });
            clipboard.on('error', e => {
                console.error('error Action:', e.action);
                console.error('error Trigger:', e.trigger);
                t.toastShow("复制失败！请重试或者手动复制内容!");
            });
        },

    },
    computed: {
    },
    mounted:function(){
    },
});

app.mount('#app');