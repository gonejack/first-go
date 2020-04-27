package writer

// Config 初始化监控服务配置信息
type Config struct {
	PathTpl      string
	BaseExt      string
	WriteExt     string
	PathInfo     map[string]string
	UpdateMoment string
	UpdatePeriod int
	UpdateSize   float64
}
