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

func (h *Handler) queryHost(c *gin.Context) {
	// 从HTTP请求的query string中获取参数
	req := host.NewQueryHostFromHTTP(c.Request)

	// 接口调用，有正常返回和失败返回
	set, err := h.svc.QueryHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, set)
}

func (h *Handler) describeHost(c *gin.Context) {
	// 从HTTP请求的query string中获取参数
	req := host.NewDescribeHostRequestWithId(c.Param("id"))

	// 接口调用，有正常返回和失败返回
	ins, err := h.svc.DescribeHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, ins)
}

func (h *Handler) putHost(c *gin.Context) {
	// 从HTTP请求的query string中获取参数
	req := host.NewPutUpdateRequest(c.Param("id"))

	// 解析Body
	if err := c.Bind(req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	req.Id = c.Param("id")

	// 接口调用，有正常返回和失败返回
	ins, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, ins)
}

func (h *Handler) patchHost(c *gin.Context) {
	// 从HTTP请求的query string中获取参数
	req := host.NewPatchUpdateRequest(c.Param("id"))

	// 解析Body
	if err := c.Bind(req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	req.Id = c.Param("id")

	// 接口调用，有正常返回和失败返回
	ins, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, ins)
}
