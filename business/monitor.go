package business

import (
	"github.com/mangenotwork/common/log"
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


TODO 记录结果存储到日志文件，定义一下格式
TODO 内存缓存记录记总监测结果

*/

func Monitor() {

	// 初始化站点对象
	InitWebSiteObj()

	go func() {
		timer := time.NewTimer(time.Second * 1) //初始化定时器
		for {
			select {
			case <-timer.C:
				// log.Info("到监测执行点...")
				// 获取站点对象
				WebSiteObj.Range(func(key any, value any) bool {
					// web := value.(*WebSiteItem)
					// go web.Run()
					return true
				})
				timer.Reset(time.Second * 1)
			}
		}
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
	log.Info("item.RateItem = ", item.RateItem)
	if item.RateItem <= 0 {
		item.RateItem = item.Rate
		log.Info("执行 " + item.HealthUri)

		// 开始执行监测
		// 获取对照组，并执行
		// 记录日志
		contrastCode, contrastMs := request(Contrast)
		if item.AlertRuleCode(contrastCode) {
			log.Error("请求对照组出现报警，请前往检查对照组!")
		}
		if item.AlertRuleMs(contrastMs) {
			log.Error("请求对照组出现报警，请前往检查对照组!")
		}
		// 监测生命URI
		// 监测站点下随机URI
		// 监测监测点

	}
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

func request(url string) (int, int64) {
	ctx, err := gt.Get(Contrast)
	if err != nil {
		log.Error(err)
		return 0, 0
	}
	return ctx.StateCode, ctx.Ms.Microseconds()
}
