package http

import (
	"github.com/gin-gonic/gin"
	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/apps/host"
)

// 通过写一个实体类，把内部的接口通过HTTP协议暴露出去
// 需要依赖内部接口的实现
// 该实体类，需要实现gin handler
type Handler struct {
	svc host.Service
}

// 面向接口，真正Service的实现，在服务实例化的时候传递进行
// 也就是在程序通过CLI start的时候
// func NewHostHTTPHandler() *Handler {
// 	return &Handler{}
// }

var handler = &Handler{}

func (h *Handler) Config() {
	// if apps.HostService == nil {
	// 	panic("Dependence host service required!")
	// }
	// 从IOC中获取HostService的实例对象
	h.svc = apps.GetImpl(host.AppName).(host.Service)
}

// 完成HTTP Handler的注册
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.createHost)
}

func (h *Handler) Name() string {
	return host.AppName
}

// 完成HTTP Handler的自注册
func init() {
	apps.RegistryGin(handler)
}
