package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github/lupingguo/go-keyborad/config"
)

func main() {
	log.Infof("app config: %s", config.GetConfig())

	// debug
	if config.GetConfig().App.DebugEvent == true {
		DebugPrintKeyPress()
	}

	// 日志初始化
	if err := initLogSetting(); err != nil {
		log.Fatalf("init log got err: %s", err)
	}

	// 注册按键开关控制
	wait := make(chan bool)
	for task := range RegisterPressHook() {
		go task.Start()
	}
	<-wait
}

// 初始化日志
func initLogSetting() error {
	cfg := config.GetLogConfig()

	// debug下输出到stdio
	if cfg.Debug == true {
		log.SetOutput(os.Stdout)
	} else {
		logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		log.SetOutput(logFile)
	}

	// 日志格式
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
	})

	return nil
}
