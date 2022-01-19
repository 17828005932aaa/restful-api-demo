package http

import (
	"net/http"
	"restful-api-demo/apps/host"
	"strconv"

	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"github.com/julienschmidt/httprouter"
)

//创建host
func (h *handler) CreateHost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 需要读取用户传底的参数,由于POST请求,我们从body里取出数据
	req := host.NewDefaultHost()
	//解析HTTP协议,通过Json反序列化,JSON --> Request
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	//组装成request对象,调用Service方法
	//1.ctx 一定要传底,如果用户中断请求,你的后端逻辑需要中断
	//2 . req:通过Http协议传递进来
	ins, err := h.host.CreateHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

//查询主机列表,分页查询
func (h *handler) QueryHost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//url中的参数
	qs := r.URL.Query()

	//设置分页默认值
	var (
		pageSize   = 20
		pageNumber = 1
	)

	//从query string读取分页参数
	psStr := qs.Get("page_size")
	if psStr != "" {
		pageSize, _ = strconv.Atoi(psStr)
	}
	pnStr := qs.Get("page_number")
	if pnStr != "" {
		pageNumber, _ = strconv.Atoi(pnStr)
	}

	req := &host.QueryHostRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}
	set, err := h.host.QueryHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

func (h *handler) DescribeHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := &host.DescribeHostRquest{
		Id: ps.ByName("id"),
	}

	set, err := h.host.DescribeHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

func (h *handler) UpdateHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := host.NewPutUpdateHostRequest()
	//解析HTTP协议,通过Json反序列化,JSON  --> Request
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	req.Id = ps.ByName("id")

	set, err := h.host.UpdateHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 传递的是一个对象
	// success, 会把你这个对象序列化成一个JSON
	// 补充返回的数据
	response.Success(w, set)
}

func (h *handler) PatchHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := host.NewPatchUpdateHostRequest()
	//解析HTTP协议,通过Json反序列化,JSON  --> Request
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	req.Id = ps.ByName("id")

	set, err := h.host.UpdateHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 传递的是一个对象
	// success, 会把你这个对象序列化成一个JSON
	// 补充返回的数据
	response.Success(w, set)
}

func (h *handler) DeleteHost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := &host.DeleteHostRequest{
		Id: ps.ByName("id"),
	}
	set, err := h.host.DeleteHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 传递的是一个对象
	// success, 会把你这个对象序列化成一个JSON
	// 补充返回的数据
	response.Success(w, set)
}
