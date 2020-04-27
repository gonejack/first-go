package monitor

import (
	"github.com/gonejack/glogger"
)

var logger = glogger.NewLogger("abc")

// 初始化全局监控实例
var instance = NewMonitor()

// Info 记录info类型的消息
// 参数
// message是消息
// args是插入到消息的参数
func Info(message string, args ...interface{}) {
	instance.Info(message, args...)
}

// Warning
// 参数
// message是消息
// args是插入到消息的参数
func Warning(message string, args ...interface{}) {
	instance.Warning(message, args...)
}

// Error 记录Error类型的消息
// 参数
// message是消息
// args是插入到消息的参数
func Error(message string, args ...interface{}) {
	instance.Error(message, args...)
}

// Start 开启监控
func Start() {
	logger.Infof("监控服务开始启动")

	instance.Start()

	logger.Infof("监控服务启动完成")
}

// Stop 关闭监控
func Stop() {
	logger.Infof("监控服务开始关闭")

	instance.Stop()

	logger.Infof("监控服务关闭完成")
}
