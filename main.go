package main

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"insane/general/base/appconfig"
	"insane/general/insane"
	"insane/server"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Info("insane server ready")
	if err := appconfig.InitConfig("./config/app.toml"); err != nil {
		logger.Debug("insane server error ", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGPIPE)
	//open a task
	go server.TK.TaskListRun()
	//configure http server
	go insane.OnStart()
	//open listen server
	go server.InsaneLoad.Start(3)
	logger.Debug("insane server starting ")

	for {
		a := <-signalChan
		fmt.Println(a)
		if a != syscall.SIGPIPE {
			break
		}
	}
	logger.Debug("insane server stop ")
}
