package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/conf"
)

type RestfulService struct {
	server *http.Server
	l      logger.Logger
	r      *restful.Container
}

// RestfulService的构造函数
func NewRestfulService() *RestfulService {
	// new restful router实例，并没有加载Handler
	r := restful.DefaultContainer

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.RestfulAddr(),
		Handler:           r,
	}

	return &RestfulService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}

}

// 启动HTTP Service
func (s *RestfulService) Start() error {
	// 加载Handler，把所有模块的Handler注册给Restful Router
	apps.InitRestful(s.r)

	// 输出已加载的APP日志信息
	apps := apps.LoadedRestfulApps()
	s.l.Infof("Loaded Restful apps :%v", apps)

	// 该操作是阻塞的，监听端口，等待请求
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Infof("Restful Server stopped success!")
			return nil
		}
		return fmt.Errorf("Start Restful Server error, %s", err)
	}

	return nil
}

// 停止HTTP Service
func (s *RestfulService) Stop() {
	s.l.Info("Start graceful shutdown ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Warnf("Shutdown Restful Server error, %s", err)
	}
}
