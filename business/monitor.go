package business

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	"small-website-monitor/model"
	"sync"
	"time"
)

func HttpCodeIsHealth(code int) bool {
	rse := true
	if code == 400 || code == 404 || code > 500 {
		return false
	}
	return rse
}

/*
	// 追加设置报警的http state code, 默认400, 404, >500
	AlarmStateCode []int64

*/

/*

监控设计

前置: 服务启动后获取监测站点列表持久数据到站点对象

场景-增删改查: 站点设置需要同步更新站点对象

场景-站点监测:
前置: 启动一个1s的定时器，每秒执行如下操作
1. 站点对象 站点监测频率值减1到0触发执行，然后复位
2. 触发执行: 执行请求生命Uri + 随机监测点 + 对照组
3. 记录结果


记录结果存储到日志文件，定义一下格式  ok
TODO 1小时采集一次URI, 没有数据则立即采集
TODO 内存缓存记录记总监测结果


*/

func Monitor() {

	// 初始化站点对象
	InitMonitor()

	// 启动监测
	go func() {
		timer := time.NewTimer(time.Second * 1) //初始化定时器
		for {
			select {
			case <-timer.C:
				// log.Info("到监测执行点...")
				// 获取站点对象
				WebSiteObj.Range(func(key any, value any) bool {
					web := value.(*WebSiteItem)
					go web.Run()
					return true
				})
				timer.Reset(time.Second * 1)
			}
		}
	}()

	// 启动定时采集
	go func() {
		log.Info("启动采集...")
		Collect()
	}()

}

func InitMonitor() {
	// 初始化站点对象
	InitWebSiteObj()
	// 启动就对站点执行一次信息进行采集并检查死链等
	go func() {
		WebSiteObj.Range(func(key any, value any) bool {
			web := value.(*WebSiteItem)
			webSiteUri := model.NewWebSiteUri(web.ID)
			webSiteUri.Collect(web.HealthUri, int(web.UriDepth))
			return true
		})
	}()
}

// WebSiteObj 站点对象
var WebSiteObj sync.Map

func InitWebSiteObj() {
	Push()
}

func Store(k string, v *WebSiteItem) {
	WebSiteObj.Store(k, v)
}

// Load get ws context
func Load(k string) *WebSiteItem {
	val, ok := WebSiteObj.Load(k)
	if ok {
		return val.(*WebSiteItem)
	}
	return nil
}

// Delete delete conn
func Delete(k string) {
	WebSiteObj.Delete(k)
}

func Push() {
	log.Info("更新站点对象")
	webSiteList, err := new(model.WebSite).GetAll()
	if err != nil {
		log.Error("加载站点对象失败," + err.Error())
		return
	}
	// 清空
	WebSiteObj.Range(func(key any, value any) bool {
		Delete(key.(string))
		return true
	})
	// 写入
	for _, v := range webSiteList {
		Store(v.ID, &WebSiteItem{
			v, v.Rate,
		})
	}
}

type WebSiteItem struct {
	*model.WebSite
	RateItem int64 // 用于计算
}

// Run 计算执行频率
func (item *WebSiteItem) Run() {
	item.RateItem--
	// log.Info("item.RateItem = ", item.RateItem)
	if item.RateItem <= 0 {
		item.RateItem = item.Rate
		log.Info("执行 " + item.HealthUri)
		mLog := &MonitorLog{
			LogType:     "Info",
			Time:        utils.NowDate(),
			HostId:      item.ID,
			Host:        item.Host,
			ContrastUri: Contrast,
			Ping:        Ping,
		}
		ping, err := gt.Ping(Ping)
		if err != nil {
			log.Error("网络不通请前往检查监测平台!", err.Error())
			mLog.LogType = LogTypeError
			mLog.Msg = "网络不通请前往检查监测平台!" + err.Error()
			mLog.Write()
			return
		}
		pingMs := ping.Milliseconds()
		mLog.PingMs = pingMs
		if pingMs >= 1000 {
			log.Error("网络环境缓慢，超过1s请前往检查监测平台!")
			mLog.LogType = LogTypeError
			mLog.Msg = "网络环境缓慢，超过1s请前往检查监测平台!"
			mLog.Write()
			return
		}
		contrastCode, contrastMs := request(Contrast)
		mLog.ContrastUriCode = contrastCode
		mLog.ContrastUriMs = contrastMs
		if item.AlertRuleCode(contrastCode) {
			log.Error("对照组请求失败，请前往检查监测平台的网络!")
			mLog.LogType = LogTypeError
			mLog.Msg = "对照组请求失败，请前往检查监测平台的网络!"
			mLog.Write()
			return
		}
		if contrastMs >= item.AlarmResTime {
			log.Error("请求对照组网络超时，请前往检查监测平台的网络!")
			mLog.LogType = LogTypeError
			mLog.Msg = "请求对照组网络超时，请前往检查监测平台的网络!"
			mLog.Write()
			return
		}
		// 监测生命URI
		HealthCode, HealthMs := request(item.HealthUri)
		mLog.Uri = item.HealthUri
		mLog.UriCode = HealthCode
		mLog.UriMs = HealthMs
		mLog.UriType = URIHealth
		// 监测规则
		if item.AlertRuleCode(HealthCode) {
			log.Error("执行报警!")
			mLog.LogType = LogTypeAlert
			mLog.Msg = fmt.Sprintf("请求失败，状态码:%d;", HealthCode)
			// TODO 记录内容用于发邮件
		}
		if item.AlertRuleMs(HealthMs) {
			log.Error("执行报警!")
			mLog.LogType = LogTypeAlert
			mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%d;", HealthMs)
			// TODO 记录内容用于发邮件
		}
		mLog.Write() // 记录日志

		// 随机取一个URI监测
		uri := model.NewWebSiteUri(item.ID)
		_, _ = uri.Get()
		if len(uri.AllUri) > 0 {
			mLog.LogType = LogTypeInfo
			mLog.Msg = ""
			randomUri := utils.RandomString(uri.AllUri)
			randomCode, randomMs := request(randomUri)
			mLog.Uri = randomUri
			mLog.UriCode = randomCode
			mLog.UriMs = randomMs
			mLog.UriType = URIRandom
			// 监测规则
			if item.AlertRuleCode(HealthCode) {
				log.Error("执行报警!")
				mLog.LogType = LogTypeAlert
				mLog.Msg = fmt.Sprintf("请求失败，状态码:%d", HealthCode)
				// TODO 记录内容用于发邮件
			}
			if item.AlertRuleMs(HealthMs) {
				log.Error("执行报警!")
				mLog.LogType = LogTypeAlert
				mLog.Msg = fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%d", HealthMs)
				// TODO 记录内容用于发邮件
			}
			mLog.Write() // 记录日志
		}

		// TODO 循环监测点监测

		// TODO 统一发邮件，这里需要注意 读一个时间锁，目的是为了防止短时间内频繁发送邮件

	}
}

func AlertRuleCode(code int) bool {
	if code == 400 || code == 404 || code >= 500 {
		return true
	}
	return false
}

// AlertRuleCode 报警规则 出现 400, 404, >500 的状态码视为出现问题
func (item *WebSiteItem) AlertRuleCode(code int) bool {
	if code == 400 || code == 404 || code >= 500 {
		return true
	}
	return false
}

// AlertRuleMs 响应时间超过设置的响应时间视为出现问题
func (item *WebSiteItem) AlertRuleMs(nowMs int64) bool {
	if nowMs >= item.AlarmResTime {
		return true
	}
	return false
}

// Contrast TODO 对照组
var Contrast = "www.baidu.com"
var Ping = "101.226.4.6"

func request(url string) (int, int64) {
	ctx, err := gt.Get(url)
	if err != nil {
		log.Error(err)
		return 0, 0
	}
	return ctx.StateCode, ctx.Ms.Milliseconds()
}
