package monitor

import (
	"encoding/json"
	"first-go/monitor/internal/writer"
	"fmt"
	"os"
	"time"
)

// monitor 健康结构体
type monitor struct {
	writer  writer.Interface
	signal  chan os.Signal
	running bool
}

// heartbeatRoutine 定时心跳
func (s *monitor) heartbeatRoutine() {
	s.running = true

	ticker := time.Tick(time.Second)

	for {
		select {
		case <-ticker:
			s.Heartbeat()
		case <-s.signal:
			logger.Infof("退出心跳线程")

			s.running = false

			return
		}
	}
}

// Heartbeat 心跳写日志
func (s *monitor) Heartbeat() {
	msg, _ := json.Marshal(map[string]interface{}{
		"source":    os.Getenv("MONITOR_APP_ID"),
		"type":      "HEARTBEAT",
		"timestamp": time.Now().Unix(),
	})

	logger.Debugf("生成心跳消息: %s", msg)

	s.writer.WriteByt(msg)
}

// Info 写提醒日志消息
// 参数
// message是消息
// args是插入到消息的参数
func (s *monitor) Info(message string, args ...interface{}) {
	msg, _ := json.Marshal(map[string]interface{}{
		"source":    os.Getenv("MONITOR_APP_ID"),
		"type":      "INFO",
		"timestamp": time.Now().Unix(),
		"detail":    fmt.Sprintf(message, args...),
	})

	logger.Debugf("生成提醒消息: %s", msg)

	s.writer.WriteByt(msg)
}

// Warning 写警告日志消息
// 参数
// message是消息
// args是插入到消息的参数
func (s *monitor) Warning(message string, args ...interface{}) {
	msg, _ := json.Marshal(map[string]interface{}{
		"source":    os.Getenv("MONITOR_APP_ID"),
		"type":      "WARNING",
		"timestamp": time.Now().Unix(),
		"detail":    fmt.Sprintf(message, args...),
	})

	logger.Debugf("生成警告消息: %s", msg)

	s.writer.WriteByt(msg)
}

// Error 写错误日志消息
// 参数
// message是消息
// args是插入到消息的参数
func (s *monitor) Error(message string, args ...interface{}) {
	msg, _ := json.Marshal(map[string]interface{}{
		"source":    os.Getenv("MONITOR_APP_ID"),
		"type":      "ERROR",
		"timestamp": time.Now().Unix(),
		"detail":    fmt.Sprintf(message, args...),
	})

	logger.Debugf("生成错误消息: %s", msg)

	s.writer.WriteByt(msg)
}

// 开启服务监控
func (s *monitor) Start() {
	s.writer.Start()

	go s.heartbeatRoutine()
}

// 关闭服务监控
func (s *monitor) Stop() {
	s.signal <- os.Interrupt

	for s.running {
		time.Sleep(time.Millisecond)
	}

	s.writer.Stop()
}

// NewMonitor 初始化监控实例
// 返回
// 监控初始化的实例
func NewMonitor() (s *monitor) {
	s = &monitor{
		signal: make(chan os.Signal),
	}

	conf := writer.Config{
		PathTpl:  "{kafka_dir}/{topic}/{app_id}_{date}_{ts}{base_ext}{write_ext}",
		BaseExt:  ".msg",
		WriteExt: "",
		PathInfo: map[string]string{
			"{kafka_dir}": os.Getenv("MONITOR_KAFKA_FILE_DIR"),
			"{topic}":     os.Getenv("MONITOR_KAFKA_TOPIC"),
			"{app_id}":    os.Getenv("MONITOR_APP_ID"),
		},
		UpdateMoment: "00:01:00",
	}

	s.writer = writer.New("心跳消息文件", conf)

	return
}
