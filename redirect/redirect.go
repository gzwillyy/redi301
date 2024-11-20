// redirect/redirect.go

package redirect

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

func Process(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// 读取请求行
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		logrus.Errorf("[redirect] Failed to read request line: %v\n", err)
		return
	}
	logrus.Debugf("Request Line: %s", requestLine)

	// 解析请求行
	parts := strings.Fields(requestLine)
	if len(parts) < 3 {
		logrus.Errorf("[redirect] Invalid request line: %s\n", requestLine)
		sendBadRequest(conn)
		return
	}
	method, path, httpVersion := parts[0], parts[1], parts[2]

	// 验证HTTP版本
	if !strings.HasPrefix(httpVersion, "HTTP/") {
		logrus.Errorf("[redirect] Invalid HTTP version: %s\n", httpVersion)
		sendBadRequest(conn)
		return
	}

	// 读取并解析头部
	var host string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			logrus.Errorf("[redirect] Failed to read headers: %v\n", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // 头部读取结束
		}
		if strings.HasPrefix(strings.ToLower(line), "host:") {
			host = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			logrus.Debugf("Host: %s", host)
		}
	}

	if host == "" {
		logrus.Errorf("[redirect] Host header not found\n")
		sendBadRequest(conn)
		return
	}

	// 根据HTTP方法进行处理
	switch strings.ToUpper(method) {
	case "GET", "HEAD":
		// 构建目标HTTPS URL
		location := fmt.Sprintf("https://%s%s", host, path)
		response := fmt.Sprintf("HTTP/1.1 301\r\n"+
			"Content-Type: text/html\r\n"+
			"Cache-Control: max-age=86400\r\n"+
			"Content-Length: "+"0"+"\r\n"+
			"Connection: close\r\n"+
			"Location: %s\r\n\r\n", location)
		// 发送重定向响应
		_, err = conn.Write([]byte(response))
		if err != nil {
			logrus.Errorf("[redirect] Failed to write redirect response: %v\n", err)
			return
		}
		logrus.Debug("[redirect] Redirect response sent")
	default:
		// 返回405 Method Not Allowed
		sendMethodNotAllowed(conn)
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
