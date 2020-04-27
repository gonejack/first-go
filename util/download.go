package util

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var logger = NewLogger("Download")

func Download(source string, target string) (uri string) {
	uri = source

	// 本地文件路径解析
	target, err := filepath.Abs(target)
	if err != nil {
		logger.Fatalf("error with save file %s: %s", source, err)
		return
	} else {
		CreativeDir(filepath.Dir(target))
	}

	// 创建本地文件
	fp, err := os.Create(target)
	if err != nil {
		logger.Errf("error with creating %s: %s", target, err)
		return
	} else {
		defer fp.Close()
	}

	// 获取远程文件
	resp, err := http.Get(source)

	// 302 过多
	if err != nil {
		logger.Errf("error with fetching %s: %s", source, err)
		return
	} else {
		defer resp.Body.Close()
	}

	// 状态码不对
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		logger.Errf("wrong http status code %d from %s", resp.StatusCode, source)
		return
	}

	// 流复制
	_, err = io.Copy(fp, resp.Body)
	if err != nil {
		logger.Errf("error downloading %s to %s", source, target)
		return
	} else {
		logger.Debf("%s => %s", source, target)
	}

	// 一切成功，返回本地url
	uri = "file://" + target

	return
}
