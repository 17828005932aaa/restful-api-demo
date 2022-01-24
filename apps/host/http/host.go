package http

import (
	"net/http"
	"restful-api-demo/apps/host"
	"strconv"
	"time"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
)

//创建host
func (h *handler) CreateHost(w http.ResponseWriter, r *http.Request) {
	// 需要读取用户传底的参数,由于POST请求,我们从body里取出数据
	//
	//用于接受前端传过来的数据
	payload := &struct {
		*host.Resource
		*host.Describe
	}{
		Resource: &host.Resource{
			CreateAt: time.Now().UnixNano() / 1000000,
		},
		Describe: &host.Describe{},
	}

	//解析HTTP协议,通过Json反序列化,JSON --> Request
	if err := request.GetDataFromRequest(r, payload); err != nil {
		response.Failed(w, err)
		return
	}
	//然后把前端传过来的结构体赋值给host
	req := host.NewDefaultHost()
	req.Resource = payload.Resource
	req.Describe = payload.Describe

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
func (h *handler) QueryHost(w http.ResponseWriter, r *http.Request) {
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
		Pagesize:   int64(pageSize),
		PageNumber: int64(pageNumber),
		Keywords:   qs.Get("keywords"),
	}
	set, err := h.host.QueryHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

func (h *handler) DescribeHost(w http.ResponseWriter, r *http.Request) {
	//从封装的context中获取ByName
	ctx := context.GetContext(r)
	req := &host.DescribeHostRquest{
		Id: ctx.PS.ByName("id"),
	}

	set, err := h.host.DescribeHost(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}

func (h *handler) UpdateHost(w http.ResponseWriter, r *http.Request) {
	//从封装的context中获取ByName
	ctx := context.GetContext(r)
	req := host.NewPutUpdateHostRequest()
	//解析HTTP协议,通过Json反序列化,JSON  --> Request
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	req.Resource.Id = ctx.PS.ByName("id")

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

func (h *handler) PatchHost(w http.ResponseWriter, r *http.Request) {
	//从封装的context中获取ByName
	ctx := context.GetContext(r)
	req := host.NewPatchUpdateHostRequest()
	//解析HTTP协议,通过Json反序列化,JSON  --> Request
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	req.Resource.Id = ctx.PS.ByName("id")

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

func (h *handler) DeleteHost(w http.ResponseWriter, r *http.Request) {
	//从封装的context中获取ByName
	ctx := context.GetContext(r)
	req := &host.DeleteHostRequest{
		Id: ctx.PS.ByName("id"),
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
