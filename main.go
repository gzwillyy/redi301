package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/gzwillyy/redi301/config"
	"github.com/gzwillyy/redi301/http"
	"github.com/gzwillyy/redi301/lagran"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(config.LogLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	config.Init()
}

func main() {
	color.Redln(config.Header)
	go lagran.Run()
	go http.Start()
	color.Greenf("[%v] App is running...\n", time.Now().Format("2006-06-02 15:04:05"))
	color.Greenf("[%v] You can use command '%v -h' to view more instructions... \n", time.Now().Format("2006-06-02 15:04:05"), os.Args[0])
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		lagran.UnsetIptable(config.HttpPort)
		logrus.Debug("Unset iptables...")
		done <- true
	}()
	<-done
	color.Redf("[%v] App terminated...\n", time.Now().Format("2006-06-02 15:04:05"))
}

// GOOS=linux GOARCH=amd64 go build -o tcp_server tcp_server.go
// chmod +x ./tcp_server
// nohup ./tcp_server --tcp-port 25125 --tls-port 25126 > server.log 2>&1 &
// ps aux | grep tcp_server
// kill 12345
