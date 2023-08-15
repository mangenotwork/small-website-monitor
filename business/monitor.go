package business

import (
	"github.com/mangenotwork/common/log"
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

需要:
1. 内存缓存来存储 站点对象
2. 一个统一方法将 监测站点列表持久数据刷到 站点对象
3. 启动一个go程进行站点监测

*/

func Monitor() {
	go func() {
		timer := time.NewTimer(time.Second * 1) //初始化定时器
		for {
			select {
			case <-timer.C:
				log.Info("到监测执行点...")
				timer.Reset(time.Second * 1)
			}
		}
	}()
}
