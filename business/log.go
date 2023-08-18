package business

import (
	"fmt"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"io"
	"os"
	"strings"
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

func ReadLog(hostId string) []*MonitorLog {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}
	fileName := logPath + hostId + "_" + utils.NowDateLayout("20060102") + ".log"
	log.Info("fileName = ", fileName)
	f, err := os.Open(fileName)
	if err != nil {
		log.ErrorF("open file error:%s", err.Error())
	}
	defer func() {
		_ = f.Close()
	}()
	data := make([]*MonitorLog, 0)
	buff := make([]byte, 0, 4096)
	char := make([]byte, 1)
	stat, _ := f.Stat()
	filesize := stat.Size()
	cursor := 0
	count := 0
	maxCount := 300
	for {
		cursor -= 1
		_, _ = f.Seek(int64(cursor), io.SeekEnd)
		_, err = f.Read(char)
		if err != nil {
			log.Error(err)
			break
		}
		if char[0] == '\n' {
			if len(buff) > 0 {
				revers(buff)
				d := toMonitorLogObj(string(buff))
				if d != nil {
					data = append(data, d)
				}
				count++
				if count == maxCount {
					break
				}
			}
			buff = buff[:0]
		} else {
			buff = append(buff, char[0])
		}
		if int64(cursor) == -filesize {
			break
		}
	}
	return data
}

func revers(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func toMonitorLogObj(str string) *MonitorLog {
	strList := strings.Split(str, "|")
	if len(strList) != 15 {
		return nil
	}
	return &MonitorLog{
		LogType:         strList[0],
		Time:            strList[1],
		HostId:          strList[2],
		Host:            strList[3],
		UriType:         strList[4],
		Uri:             strList[5],
		UriCode:         utils.AnyToInt(strList[6]),
		UriMs:           utils.AnyToInt64(strList[7]),
		ContrastUri:     strList[8],
		ContrastUriCode: utils.AnyToInt(strList[9]),
		ContrastUriMs:   utils.AnyToInt64(strList[10]),
		Ping:            strList[11],
		PingMs:          utils.AnyToInt64(strList[12]),
		Msg:             strList[13],
	}
}
