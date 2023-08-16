package business

import (
	"fmt"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"io"
	"os"
)

// 监测日志

// TODO 只保留7天的日志

type MonitorLog struct {
	LogType         string // Info  Alert Error
	Time            string
	HostId          string
	Host            string
	UriType         string // 监测的URI类型 Health:根URI,健康URI  Random:随机URI  Point:监测点URI
	Uri             string // URI
	UriCode         int    // URI响应码
	UriMs           int64  // URI响应时间
	ContrastUri     string // 对照组URI
	ContrastUriCode int    // 对照组URI响应码
	ContrastUriMs   int64  // 对照组URI响应时间
	Ping            string
	PingMs          int64
	Msg             string
}

const (
	URIHealth    = "Health"
	URIRandom    = "Random"
	URIPoint     = "Point"
	LogTypeInfo  = "Info"
	LogTypeAlert = "Alert"
	LogTypeError = "Error"
)

// 写日志
func (m *MonitorLog) Write() {
	logPath, err := conf.YamlGetString("logPath")
	log.Info("logPath = ", logPath)
	if err != nil {
		logPath = "./log/"
	}
	fileName := logPath + m.HostId + "_" + utils.NowDateLayout("20060102") + ".log"
	log.Info("fileName = ", fileName)
	var file *os.File
	if !utils.Exists(fileName) {
		_ = os.MkdirAll(logPath, 0666)
		file, _ = os.Create(fileName)
	} else {
		file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}
	defer func() {
		_ = file.Close()
	}()
	n, err := io.WriteString(file, m.DataFormat())
	if err != nil {
		fmt.Println("写入错误：", err)
		return
	}
	fmt.Println("写入成功：n=", n)
}

func (m *MonitorLog) DataFormat() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d|%d|%s|%d|%d|%s|%d|%s|\r\n",
		m.LogType, m.Time, m.HostId, m.Host, m.UriType, m.Uri, m.UriCode, m.UriMs, m.ContrastUri,
		m.ContrastUriCode, m.ContrastUriMs, m.Ping, m.PingMs, m.Msg)
}
