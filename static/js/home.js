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
            }
        }
    },
    created:function(){
        var t = this;
        t.getList();
        t.getMail();
        t.getMailInfo();
    },
    methods: {
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
                    t.msg = data.msg;
                    const toastLiveExample = document.getElementById('liveToast')
                    const toast = new bootstrap.Toast(toastLiveExample)
                    toast.show()
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
                    t.msg = data.msg;
                    const toastLiveExample = document.getElementById('liveToast')
                    const toast = new bootstrap.Toast(toastLiveExample)
                    toast.show()
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
                    t.msg = data.msg;
                    const toastLiveExample = document.getElementById('liveToast')
                    const toast = new bootstrap.Toast(toastLiveExample)
                    toast.show()
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
                t.msg = "请输入URI";
                const toastLiveExample = document.getElementById('liveToast')
                const toast = new bootstrap.Toast(toastLiveExample)
                toast.show()
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
                    t.msg = data.msg;
                    const toastLiveExample = document.getElementById('liveToast')
                    const toast = new bootstrap.Toast(toastLiveExample)
                    toast.show()
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
                    t.msg = data.msg;
                    const toastLiveExample = document.getElementById('liveToast')
                    const toast = new bootstrap.Toast(toastLiveExample)
                    toast.show()
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
                        t.msg = data.msg;
                        const toastLiveExample = document.getElementById('liveToast')
                        const toast = new bootstrap.Toast(toastLiveExample)
                        toast.show()
                    }
                },
                error: function(xhr,textStatus) {
                    console.log(xhr, textStatus);
                }
            });
        }

    },
    computed: {
    },
    mounted:function(){
    },
});

app.mount('#app');