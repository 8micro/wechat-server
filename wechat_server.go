package main

import "github.com/8micro/wechat-server/api"
import "github.com/8micro/wechat-server/etc"
import "github.com/8micro/wechat-server/security"
import "github.com/8micro/wechat-server/server"
import "github.com/8micro/gounits/flocker"
import "github.com/8micro/gounits/logger"
import "github.com/8micro/gounits/rand"

import (
	"flag"
	"time"
)

//WechatServer is exported
type WechatServer struct {
	RetryStartup bool
	Locker       *flocker.FileLocker
	APIServer    *api.Server
	CoreServer   *server.CoreServer
}

//NewWechatServer is exported
func NewWechatServer() (*WechatServer, error) {

	var confFile, originID string
	flag.StringVar(&confFile, "f", "./etc/config.yaml", "wechat server etc file.")
	flag.StringVar(&originID, "o", "", "wechat account originid.")
	flag.Parse()

	if err := etc.New(confFile); err != nil {
		return nil, err
	}

	logger.OPEN(etc.LoggerArgs())
	key, err := rand.UUIDFile("./8m_wechat_server.key") //服务器唯一标识文件
	if err != nil {
		return nil, err
	}

	var fLocker *flocker.FileLocker
	if pidFile := etc.PidFile(); pidFile != "" {
		fLocker = flocker.NewFileLocker(pidFile, 0)
	}

	securityData, err := security.New(etc.SecurityRoot(), originID)
	if err != nil {
		return nil, err
	}

	coreServer := server.NewCoreServer(key, securityData)
	api.RegisterStore("Key", key)
	api.RegisterStore("SystemConfig", etc.SystemConfig)
	api.RegisterStore("CoreServer", coreServer)
	listen := etc.SystemConfig.Listen
	apiServer := api.NewServer(listen.Hosts, listen.EnableCors, nil)
	return &WechatServer{
		RetryStartup: etc.RetryStartup(),
		Locker:       fLocker,
		APIServer:    apiServer,
		CoreServer:   coreServer,
	}, nil
}

//Startup is exported
func (wechatServer *WechatServer) Startup() error {

	var err error
	for {
		if err != nil {
			if wechatServer.RetryStartup == false {
				return err
			}
			time.Sleep(time.Second * 10) //retry, after sleep 10 seconds.
		}

		wechatServer.Locker.Unlock()
		if err = wechatServer.Locker.Lock(); err != nil {
			logger.ERROR("[#main#] pidfile lock error, %s", err)
			continue
		}
		break
	}

	go func() {
		logger.INFO("[#main#] API listener: %s", wechatServer.APIServer.ListenHosts())
		if err := wechatServer.APIServer.Startup(); err != nil {
			logger.ERROR("[#main#] API startup error, %s", err.Error())
		}
	}()
	logger.INFO("[#main#] wechat server started.")
	logger.INFO("[#main#] key:%s", wechatServer.CoreServer.Key)
	return nil
}

//Stop is exported
func (wechatServer *WechatServer) Stop() error {

	wechatServer.Locker.Unlock()
	logger.INFO("[#main#] wechat server stoped.")
	logger.CLOSE()
	return nil
}
