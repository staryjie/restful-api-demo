package http

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
	"github.com/staryjie/restful-api-demo/apps/host"
)

// 用于暴露Host Service接口
func (h *Handler) createHost(c *gin.Context) {
	ins := host.NewHost()
	// 解析用户传递的参数
	if err := c.Bind(ins); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	// 进行接口调用
	ins, err := h.svc.CreateHost(c.Request.Context(), ins)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, ins)
}
