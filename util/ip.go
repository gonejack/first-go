package util

func getIp() string {
	cn, err := net.DialTimeout("udp", "8.8.8.8:80", time.Second)

	if err == nil {
		ip, _, err := net.SplitHostPort(cn.LocalAddr().String())

		if err == nil {
			return ip
		}
	}

	fmt.Printf("解析网络IP地址失败: %s", err)

	return ""
}
