package http

import (
	"github.com/gin-gonic/gin"
	"github.com/staryjie/restful-api-demo/apps/host"
)

// 通过写一个实体类，把内部的接口通过HTTP协议暴露出去
// 需要依赖内部接口的实现
// 该实体类，需要实现gin handler
type Handler struct {
	svc host.Service
}

func NewHTTPHandler(svc host.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// 完成HTTP Handler的注册
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/host", h.createHost)
}
