package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"

	"github.com/staryjie/restful-api-demo/apps/book"
)

// go-restful HTTP handler
func (u *handler) CreateBook(r *restful.Request, w *restful.Response) {
	req := book.NewCreateBookRequest()

	// 实现功能 就是反序列换
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := u.service.CreateBook(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, set)
}

func (u *handler) QueryBook(r *restful.Request, w *restful.Response) {
	req := book.NewQueryBookRequestFromHTTP(r.Request)
	set, err := u.service.QueryBook(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (u *handler) DescribeBook(r *restful.Request, w *restful.Response) {
	req := book.NewDescribeBookRequest(r.PathParameter("id"))
	ins, err := u.service.DescribeBook(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, ins)
}

// 通过PathParameter读取路径中的{id}参数
func (u *handler) UpdateBook(r *restful.Request, w *restful.Response) {
	req := book.NewPutBookRequest(r.PathParameter("id"))

	if err := r.ReadEntity(req.Data); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := u.service.UpdateBook(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (u *handler) PatchBook(r *restful.Request, w *restful.Response) {
	req := book.NewPatchBookRequest(r.PathParameter("id"))

	if err := r.ReadEntity(req.Data); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := u.service.UpdateBook(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (u *handler) DeleteBook(r *restful.Request, w *restful.Response) {
	req := book.NewDeleteBookRequestWithID(r.PathParameter("id"))
	set, err := u.service.DeleteBook(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
