// config/config.go

package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

var (
	HttpAddr string
	// Target   string // 移除或保留视需求而定
	HttpPort string
)

const (
	LogLevel = logrus.FatalLevel
	Header   = `=========================================================`
)

func Init() {
	flag.Usage = func() {
		color.Redln(Header)
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&HttpAddr, "a", "0.0.0.0:80", "The listen address and port.")
	// flag.StringVar(&Target, "t", "https://www.microsoft.com", "The prefix of target redirect url.") // 移除或保留视需求而定
	flag.Parse()
	ipAndPort := strings.Split(HttpAddr, ":")
	if len(ipAndPort) != 2 {
		logrus.Fatalf("http listen address error...")
	} else {
		HttpPort = ipAndPort[1]
	}
	// 如果保留 Target，确保其为有效URL
	/*
		if _, err := url.Parse(Target); err != nil {
			logrus.Fatalf("prefix of target redirect url error...")
		}
	*/
}
