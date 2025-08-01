package tools

import (
	"io"
	"net/http"
)

func GetPublicIP() (ip string, err error) {
	// - http://myexternalip.com/raw
	// - http://ip.dhcp.cn/?ip
	// - https://www.bt.cn/Api/getIpAddress
	resp, err := http.Get("http://myexternalip.com/raw")
	content, err := io.ReadAll(resp.Body)
	if err == nil {
		return string(content), nil
	}
	return "127.0.0.1", err
}
