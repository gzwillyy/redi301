// redirect/redirect.go

package redirect

import (
	"fmt"
	"net"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/sirupsen/logrus"
)

var (
	rePath = regexp2.MustCompile(`(GET|HEAD?) (.*?) HTTP`, 0)
	reHost = regexp2.MustCompile(`(?i)Host:\s*([^\r\n]+)`, 0)
)

func Process(conn net.Conn) {
	var RequestHead string
	var Path string
	var Host string
	defer conn.Close()
	for {
		var buf [128]byte
		// 接受数据
		n, err := conn.Read(buf[:])
		if err != nil {
			logrus.Errorf("[redirect] read from connect failed, err: %v\n", err)
			break
		}
		logrus.Debugf("Receive data: \n%s\n", string(buf[:n]))
		RequestHead += string(buf[:n])
		logrus.Debug(RequestHead)
		if strings.Contains(RequestHead, "HTTP") {
			// 提取 Path
			if matchArr, err := rePath.FindStringMatch(RequestHead); err != nil {
				Path = "/"
			} else {
				if matchArr != nil && len(matchArr.Groups()) > 2 {
					Path = matchArr.Groups()[2].String()
				} else {
					Path = "/"
				}
			}

			// 提取 Host
			if matchHost, err := reHost.FindStringMatch(RequestHead); err != nil || matchHost == nil {
				Host = "localhost" // 设置默认 Host
			} else {
				if len(matchHost.Groups()) > 1 {
					Host = matchHost.Groups()[1].String()
				} else {
					Host = "localhost"
				}
			}

			// 构建基于 Host 和 Path 的 Location URL
			location := fmt.Sprintf("https://%s%s", Host, Path)

			// 发送 301 重定向响应
			if _, err = conn.Write([]byte("HTTP/1.1 301\r\n" +
				"Content-Type: text/html\r\n" +
				"Cache-Control: max-age=86400\r\n" +
				"Content-Length: " + "0" + "\r\n" +
				"Connection: close\r\n" +
				"Location: " + location + "\r\n\r\n")); err != nil {
				logrus.Errorf("write to client failed, err: %v\n", err)
				break
			}
			logrus.Debug("[redirect] Response sent...")
			break
		}
	}
}

// sendBadRequest 发送400 Bad Request响应
func sendBadRequest(conn net.Conn) {
	response := "HTTP/1.1 400 Bad Request\r\n" +
		"Content-Length: 0\r\n" +
		"Connection: close\r\n\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		logrus.Errorf("[redirect] Failed to write Bad Request response: %v\n", err)
	}
}

// sendMethodNotAllowed 发送405 Method Not Allowed响应
func sendMethodNotAllowed(conn net.Conn) {
	response := "HTTP/1.1 405 Method Not Allowed\r\n" +
		"Allow: GET, HEAD\r\n" +
		"Content-Length: 0\r\n" +
		"Connection: close\r\n\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		logrus.Errorf("[redirect] Failed to write Method Not Allowed response: %v\n", err)
	}
}
