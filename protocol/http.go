package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/conf"
)

type HttpService struct {
	server *http.Server
	l      logger.Logger
	r      gin.IRouter
}

// HttpService的构造函数
func NewHttpService() *HttpService {
	// new gin router实例，并没有加载Handler
	r := gin.Default()

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HttpAddr(),
		Handler:           r,
	}

	return &HttpService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}

}

// 启动HTTP Service
func (s *HttpService) Start() error {
	// 加载Handler，把所有模块的Handler注册给Gin Router
	apps.InitGin(s.r)

	// 输出已加载的APP日志信息
	apps := apps.LoadedGinApps()
	s.l.Infof("Loaded gin apps :%v", apps)

	// 该操作是阻塞的，监听端口，等待请求
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Infof("Http Server stopped success!")
			return nil
		}
		return fmt.Errorf("Start Http Server error, %s", err)
	}

	return nil
}

// 停止HTTP Service
func (s *HttpService) Stop() {
	s.l.Info("Start graceful shutdown ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Warnf("Shutdown Http Server error, %s", err)
	}
}
