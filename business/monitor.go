package business

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	"small-website-monitor/global"
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

监控设计

监控原理: 对Uri进行请求通过响应状态码和响应时间进行判断，
在此基础上需要有对照组和当前网络情况作为前置条件加以判断前提，
并且对站点进行根Uri+随机Uri+指定监测点进行组合监测，
最终实现监测。

前置: 服务启动后获取监测站点列表持久数据到站点对象

场景-增删改查: 站点设置需要同步更新站点对象

场景-站点监测:
前置: 启动一个1s的定时器，每秒执行如下操作
1. 站点对象 站点监测频率值减1到0触发执行，然后复位
2. 触发执行: 执行请求生命Uri + 随机监测点 + 对照组
3. 记录结果

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
					web.Run() // 不要使用并发，影响当前网络环境
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

// Run 执行监测
func (item *WebSiteItem) Run() {
	item.RateItem--
	if item.RateItem <= 0 {
		// 计算频率复位
		item.RateItem = item.Rate
		//log.Info("执行 " + item.HealthUri)
		// 报警数据初始化
		alert := &model.AlertBody{
			Synopsis: "监测到" + item.Host + "站点出现问题，请快速前往检查并处理!",
			Tr:       make([]*model.AlertTd, 0),
		}
		masterConf := model.GetMasterConf()
		// 日志数据
		mLog := &MonitorLog{
			LogType:     "Info",
			Time:        utils.NowDate(),
			HostId:      item.ID,
			Host:        item.Host,
			ContrastUri: masterConf.ContrastUri,
			Ping:        masterConf.Ping,
		}
		// ping一下，检查当前网络环境
		_, pingRse := item.Ping(masterConf, mLog)
		if !pingRse {
			// 网络环境异常不执行监测
			return
		}
		// 请求对照组，对照组有问题不执行监测
		if item.Contrast(masterConf, mLog) {
			return
		}
		// 监测生命URI
		isAlert1 := item.MonitorHealthUri(masterConf, mLog, alert)
		// 随机URI监测
		isAlert2 := item.MonitorRandomUri(masterConf, mLog, alert)
		// 循环监测点监测
		isAlert3 := item.MonitorPointUri(masterConf, mLog, alert)
		// 发邮件
		if isAlert1 || isAlert2 || isAlert3 {
			now := time.Now().Unix()
			if now-global.LastSendMail < 60 {
				log.Info("邮件发送太频繁，保持1分钟的间隔")
			} else {
				global.LastSendMail = now
				model.Send(alert.Synopsis, alert.Html())
			}
		}
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

func request(url string) (int, int64) {
	ctx, err := gt.Get(url, gt.SetSleepMs(10, 100))
	if err != nil {
		log.Error(err)
		model.SetMonitorErrInfo(url + " 请求失败err:" + err.Error())
		return 0, 0
	}
	return ctx.StateCode, ctx.Ms.Milliseconds()
}

func (item *WebSiteItem) Ping(masterConf *model.MasterConf, mLog *MonitorLog) (int64, bool) {
	ping, err := gt.Ping(masterConf.Ping)
	if err != nil {
		mLog.LogType = LogTypeError
		mLog.Msg = "网络不通请前往检查监测平台!" + err.Error()
		mLog.Write()
		model.SetMonitorErrInfo(mLog.Msg)
		return 0, false
	}
	pingMs := ping.Milliseconds()
	mLog.PingMs = pingMs
	if pingMs >= 1000 {
		mLog.LogType = LogTypeError
		mLog.Msg = fmt.Sprintf("网络环境缓慢，超过1s(%d)请前往检查监测平台!", pingMs)
		mLog.Write()
		model.SetMonitorErrInfo(mLog.Msg)
		return pingMs, false
	}
	return pingMs, true
}

func (item *WebSiteItem) Contrast(masterConf *model.MasterConf, mLog *MonitorLog) bool {
	contrastCode, contrastMs := request(masterConf.ContrastUri)
	mLog.ContrastUriCode = contrastCode
	mLog.ContrastUriMs = contrastMs
	contrastErr := false
	if item.AlertRuleCode(contrastCode) {
		contrastErr = true
		mLog.Msg += fmt.Sprintf("对照组请求失败code=%d!", contrastCode)
		model.SetMonitorErrInfo(mLog.Msg)
	}
	if contrastMs >= masterConf.ContrastTime {
		contrastErr = true
		mLog.Msg += fmt.Sprintf("请求对照组网络超%dms 当前为%dms!", masterConf.ContrastTime, contrastMs)
		model.SetMonitorErrInfo(mLog.Msg)
	}
	if contrastErr {
		mLog.LogType = LogTypeError
		mLog.Write()
		return contrastErr
	}
	return contrastErr
}

func (item *WebSiteItem) MonitorHealthUri(masterConf *model.MasterConf, mLog *MonitorLog, alert *model.AlertBody) bool {
	// =================================  监测生命URI
	//log.Info("=================================  监测生命URI... ")
	times := 0
R:
	healthCode, healthMs := request(item.HealthUri)
	mLog.Uri = item.HealthUri
	mLog.UriCode = healthCode
	mLog.UriMs = healthMs
	mLog.UriType = URIHealth
	mLog.LogType = LogTypeInfo
	mLog.Msg = ""
	healthAlert := false
	// 监测规则
	if item.AlertRuleCode(healthCode) {
		healthAlert = true
		mLog.LogType = LogTypeAlert
		mLog.Msg = fmt.Sprintf("请求失败，状态码:%d;", healthCode)
	}
	if item.AlertRuleMs(healthMs) {
		// 如果是超时再来一次,确保监测是连续超时并非单个网络波动
		if times == 0 {
			times++
			goto R
		}
		healthAlert = true
		mLog.LogType = LogTypeAlert
		mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms;", healthMs)
	}

	NowMonitorSet(&NowMonitor{
		HostId:     item.ID,
		Code:       healthCode,
		Ms:         healthMs,
		ContrastMs: mLog.ContrastUriMs,
		PingMs:     mLog.PingMs,
	})

	if healthAlert {
		date := utils.NowDate()
		// 记录报警
		alertObj := model.NewWebSiteAlert(item.ID)
		err := alertObj.Add(&model.AlertData{
			Date:          date,
			Uri:           item.HealthUri,
			UriCode:       healthCode,
			UriMs:         healthMs,
			ContrastUriMs: mLog.ContrastUriMs,
			PingMs:        mLog.PingMs,
			Msg:           mLog.Msg,
		})
		if err != nil {
			log.Error("记录报警信息失败:" + err.Error())
		}
		// 记录信息用于发邮件通知
		alert.Tr = append(alert.Tr, &model.AlertTd{
			Date: utils.NowDate(),
			Host: item.Host,
			Uri:  item.HealthUri,
			Code: healthCode,
			Ms:   healthMs,
			NetworkEnv: fmt.Sprintf("ping:%dms; 对照组(%s):%dms",
				mLog.PingMs, masterConf.ContrastUri, mLog.ContrastUriMs),
			Msg: mLog.Msg,
		})
	}
	if mLog.LogType == LogTypeInfo {
		mLog.Msg = "passed"
	}
	mLog.Write() // 记录日志
	return healthAlert
}

func (item *WebSiteItem) MonitorRandomUri(masterConf *model.MasterConf, mLog *MonitorLog, alert *model.AlertBody) bool {
	// =================================  随机取一个URI监测
	//log.Info("=================================  随机取一个URI监测... ")
	uri := model.NewWebSiteUri(item.ID)
	_, _ = uri.Get()
	if len(uri.AllUri) > 0 {
		mLog.LogType = LogTypeInfo // 复位
		mLog.Msg = ""              // 复位
		randomUri := utils.RandomString(uri.AllUri)
		times := 0
	R:
		randomCode, randomMs := request(randomUri)
		mLog.Uri = randomUri
		mLog.UriCode = randomCode
		mLog.UriMs = randomMs
		mLog.UriType = URIRandom
		randomAlert := false
		if item.AlertRuleCode(randomCode) {
			randomAlert = true
			mLog.LogType = LogTypeAlert
			mLog.Msg = fmt.Sprintf("请求失败，状态码:%d", randomCode)
		}
		if item.AlertRuleMs(randomMs) {
			if times == 0 {
				times++
				goto R
			}
			randomAlert = true
			mLog.LogType = LogTypeAlert
			mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms", randomMs)
		}
		if randomAlert {
			date := utils.NowDate()
			// 记录报警
			alertObj := model.NewWebSiteAlert(item.ID)
			err := alertObj.Add(&model.AlertData{
				Date:          date,
				Uri:           randomUri,
				UriCode:       randomCode,
				UriMs:         randomMs,
				ContrastUriMs: mLog.ContrastUriMs,
				PingMs:        mLog.PingMs,
				Msg:           mLog.Msg,
			})
			if err != nil {
				log.Error("记录报警信息失败:" + err.Error())
			}
			// 记录内容用于发邮件
			alert.Tr = append(alert.Tr, &model.AlertTd{
				Date: date,
				Host: item.Host,
				Uri:  randomUri,
				Code: randomCode,
				Ms:   randomMs,
				NetworkEnv: fmt.Sprintf("ping:%dms; 对照组(%s):%dms",
					mLog.PingMs, masterConf.ContrastUri, mLog.ContrastUriMs),
				Msg: mLog.Msg,
			})
		}
		if mLog.LogType == LogTypeInfo {
			mLog.Msg = "passed"
		}
		mLog.Write() // 记录日志
		return randomAlert
	}
	return false
}

func (item *WebSiteItem) MonitorPointUri(masterConf *model.MasterConf, mLog *MonitorLog, alert *model.AlertBody) bool {
	// =================================  循环监测点监测
	point := model.NewWebSitePoint(item.ID)
	err := point.Get()
	hasAlert := false
	if err == nil && len(point.Uri) > 0 {
		for _, v := range point.Uri {
			times := 0
		R:
			mLog.LogType = LogTypeInfo // 复位
			mLog.Msg = ""              // 复位
			pointCode, pointMs := request(v)
			mLog.Uri = v
			mLog.UriCode = pointCode
			mLog.UriMs = pointMs
			mLog.UriType = URIPoint
			vAlert := false
			// 监测规则
			if item.AlertRuleCode(pointCode) {
				vAlert = true
				mLog.LogType = LogTypeAlert
				mLog.Msg = fmt.Sprintf("请求失败，状态码:%d", pointCode)
			}
			if item.AlertRuleMs(pointMs) {
				if times == 0 {
					times++
					goto R
				}
				vAlert = true
				mLog.LogType = LogTypeAlert
				mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms", pointMs)
			}
			if vAlert {
				hasAlert = true
				date := utils.NowDate()
				alertObj := model.NewWebSiteAlert(item.ID)
				err = alertObj.Add(&model.AlertData{
					Date:          date,
					Uri:           v,
					UriCode:       pointCode,
					UriMs:         pointMs,
					ContrastUriMs: mLog.ContrastUriMs,
					PingMs:        mLog.PingMs,
					Msg:           mLog.Msg,
				})
				if err != nil {
					log.Error("记录报警信息失败:" + err.Error())
				}
				// 记录内容用于发邮件
				alert.Tr = append(alert.Tr, &model.AlertTd{
					Date: date,
					Host: item.Host,
					Uri:  v,
					Code: pointCode,
					Ms:   pointMs,
					NetworkEnv: fmt.Sprintf("ping:%dms; 对照组(%s):%dms",
						mLog.PingMs, masterConf.ContrastUri, mLog.ContrastUriMs),
					Msg: mLog.Msg,
				})
			}
			if mLog.LogType == LogTypeInfo {
				mLog.Msg = "passed"
			}
			mLog.Write() // 记录日志
		}
	}
	return hasAlert
}
