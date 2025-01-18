package main

import (
	"os"
	"os/signal"
	"simulator/etc"
	"simulator/generate"
	"simulator/httpService"
	"simulator/util"
	"syscall"
	"time"
)

var (
	cfgPath = "./etc.yaml"
)

func main() {
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	// 解析配置文件
	etc.GetConfig(cfgPath)

	start := time.Now()
	if etc.IsVirtual() {
		generate.VirtualEnvGenerateData()
	} else {
		go func() {
			generate.TruthEnvGenerateData()
			err := httpService.Start()
			if err != nil {
				util.Log.Fatalf("启动服务失败：%v", err)
			}
		}()
		listenSignal()
	}
	util.Log.Printf("耗时：%vs", time.Now().Sub(start).Seconds())
}

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	util.Log.Printf("收到信号：%v", sig.String())
	generate.Stop()
	os.Exit(0)
}
