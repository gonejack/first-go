package writer

import (
	"bufio"
	"fmt"
	"github.com/gonejack/glogger"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var logger = glogger.NewLogger("abc")

// writer 监控写日志结构体
type writer struct {
	// 文件路径格式
	pathTpl string
	// 后缀名
	baseExt string
	// 文件正在写入中的后缀名
	writeExt string
	// 自定义信息
	pathInfo map[string]string
	// 文件每天更换的时间
	updateMoment string
	// 多久检查一次 文件是否到了关闭时间点
	updatePeriod int
	// 更新长度
	updateSize float64

	// 字符串通道
	strInput chan string
	// 字节通道
	bytInput chan []byte
	// 关闭通道
	sigInput chan os.Signal

	// 文件指针
	fp *os.File
	// 文件写操作符
	writer *bufio.Writer
	// 已经积累长度
	wroteLen float64

	// 写入时间
	flushTimer <-chan time.Time
	// 关闭时间
	closeTimer <-chan time.Time
}

func (w *writer) mainRoutine() {
	stopping := false

	for {
		select {
		case <-w.sigInput:
			stopping = true
		default:
			if stopping && len(w.strInput) == 0 && len(w.bytInput) == 0 {
				w.close()
				return
			} else {
				select {
				case <-w.flushTimer:
					w.flush()
				case <-w.closeTimer:
					w.close()
				case str := <-w.strInput:
					w.write(str, nil)
				case byt := <-w.bytInput:
					w.write("", byt)
				}
			}
		}
	}
}

// 写监控消息
// 参数
// s是写入的消息（字符串类型）
// b是写入的消息（字节类型）
func (w *writer) write(s string, b []byte) {
	if w.open() {
		if sLen := len(s); sLen > 0 {
			_, err := w.writer.WriteString(s)
			if err != nil {
				logger.Errorf("写内容失败: %s", err)
			}

			w.wroteLen += float64(sLen)

			if s[sLen-1] != byte('\n') {
				_, err := w.writer.WriteRune('\n')
				if err != nil {
					logger.Errorf("写内容失败: %s", err)
				}

				w.wroteLen += 1
			}
		}
		if bLen := len(b); bLen > 0 {
			_, err := w.writer.Write(b)
			if err != nil {
				logger.Errorf("写内容失败: %s", err)
			}

			w.wroteLen += float64(bLen)

			if b[bLen-1] != byte('\n') {
				_, err := w.writer.WriteRune('\n')
				if err != nil {
					logger.Errorf("写内容失败: %s", err)
				}

				w.wroteLen += 1
			}
		}

		if w.updateSize > 0 && w.wroteLen > w.updateSize {
			logger.Infof("写入长度达到限制[%.0f bytes]", w.updateSize)
			w.close()
		}
	} else {
		logger.Errorf("没有打开的文件，写入失败")
	}
}

// flush 写文件
func (w *writer) flush() {
	if w.writer != nil {
		err := w.writer.Flush()
		if err != nil {
			logger.Errorf("写内容失败: %s", err)
		}
	}
}

// open 打开文件
// 返回
// 成功返回true，失败返回false
func (w *writer) open() bool {
	if w.fp == nil {
		if p := w.getPath(); p != "" {
			fp, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err == nil {
				logger.Infof("创建文件[%s]", p)

				w.fp = fp
				w.writer = bufio.NewWriter(fp)
				w.flushTimer = time.Tick(time.Second)
				w.closeTimer = time.After(time.Duration(w.getSleepSec()) * time.Second)

				if w.updateSize > 0 {
					logger.Infof("文件计划关闭大小为: %.0f bytes", w.updateSize)
				}
			} else {
				logger.Errorf("创建文件[%s]失败：%s", p, err)
			}
		}
	}

	return w.writer != nil
}

// close 关闭文件
func (w *writer) close() {
	if w.fp != nil {
		w.flush()

		err := w.fp.Close()
		if err != nil {
			logger.Errorf("关闭文件出错: %s", err)
		}

		logger.Infof("关闭文件[%s]", w.fp.Name())

		if w.writeExt != "" {
			w.renameFile()
		}
	}

	w.fp = nil
	w.writer = nil
	w.flushTimer = nil
	w.closeTimer = nil
	w.wroteLen = 0
}

// getSleepSec 获取下次执行定时任务的睡眠时间
// 返回
// 睡眠时间（秒）
func (w *writer) getSleepSec() (sec int) {
	now := time.Now()

	start := now
	if w.updateMoment != "" {
		todayMom := fmt.Sprintf("%s %s", now.Format("2006-01-02"), w.updateMoment)

		if t, e := time.ParseInLocation("2006-01-02 15:04:05", todayMom, now.Location()); e == nil {
			start = t
		} else {
			logger.Errorf("配置的文件更新时刻[%s]解析出错：%s，缺省为现在时刻", w.updateMoment, e)
		}
	}

	next := start
	for next.Before(now) || next.Equal(now) {
		if w.updatePeriod > 0 {
			next = next.Add(time.Duration(w.updatePeriod) * time.Second)
		} else {
			next = next.Add(time.Hour * 24)
		}
	}

	sec = int(math.Ceil(next.Sub(now).Seconds()))

	logger.Infof("文件计划关闭时间为: %s", next.Format("2006-01-02 15:04:05"))

	return
}

// getPath 获取文件路径
// 返回
// 文件路径
func (w *writer) getPath() (path string) {
	now := time.Now()
	// 基本信息
	info := map[string]string{
		"{date}":      now.Format("20060102"),
		"{ts}":        strconv.FormatInt(now.Unix(), 10),
		"{base_ext}":  w.baseExt,
		"{write_ext}": w.writeExt,
	}
	// 自定义信息
	for macro, replacement := range w.pathInfo {
		info[macro] = replacement
	}

	// 模板替换
	var replaces []string
	for macro, replacement := range info {
		replaces = append(replaces, macro, replacement)
	}
	path = strings.NewReplacer(replaces...).Replace(w.pathTpl)

	// 创文件夹
	dir := filepath.Dir(path)
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		if e := os.MkdirAll(dir, 0755); e == nil {
			logger.Infof("创建文件夹[%s]", dir)
		} else {
			logger.Errorf("创建文件夹[%s]失败: %s", dir, e)

			path = ""
		}
	}

	return
}

// renameFile 重命名文件
func (w *writer) renameFile() {
	if o := w.fp.Name(); strings.HasSuffix(o, w.writeExt) {
		n := strings.TrimSuffix(o, w.writeExt)
		e := os.Rename(o, n)

		if e == nil {
			logger.Infof("重命名文件[%s => %s]", o, n)
		} else {
			logger.Errorf("重命名文件[%s => %s]失败: %s", o, n, e)
		}
	}
}

// WriteStr 写string消息
// 参数
// str消息(字符串类型)
func (w *writer) WriteStr(str string) {
	w.strInput <- str
}

// WriteByt 写byte消息
// 参数
// byts消息(字节类型)
func (w *writer) WriteByt(byts []byte) {
	w.bytInput <- byts
}

func (w *writer) Start() {
	logger.Debugf("监控写服务开始启动")

	go w.mainRoutine()

	logger.Debugf("监控写服务启动完成")
}

func (w *writer) Stop() {
	logger.Debugf("监控写服务开始关闭")

	w.sigInput <- os.Interrupt

	for w.writer != nil {
		time.Sleep(time.Millisecond * 100)
	}

	logger.Debugf("监控写服务关闭完成")
}

// New 新建写文件实例
// 参数
// name是名字
// conf是配置信息
// 返回
// 写文件实例
func New(name string, conf Config) (w Interface) {
	w = &writer{
		pathTpl:  conf.PathTpl,
		baseExt:  conf.BaseExt,
		writeExt: conf.WriteExt,
		pathInfo: conf.PathInfo,

		updateMoment: conf.UpdateMoment,
		updatePeriod: conf.UpdatePeriod,
		updateSize:   conf.UpdateSize,

		strInput: make(chan string, 100),
		bytInput: make(chan []byte, 100),
		sigInput: make(chan os.Signal),
	}

	return
}
