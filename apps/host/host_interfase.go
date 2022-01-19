package host

import "context"

//定义host业务需要实现的功能模块,实现统一接口
type Service interface {
	//录入主机信息
	CreateHost(context.Context,*Host) (*Host,error)
	//查询主机列表信息
	QueryHost(context.Context,*QueryHostRequest) (*HostSet,error)
	//主机详情查询
	DescribeHost(context.Context,*DescribeHostRquest) (*Host,error)
	//主机信息修改
	UpdateHost(context.Context, *UpdateHostRequest) (*Host,error)
	//删除主机
	DeleteHost(context.Context,*DeleteHostRequest) (*Host,error)
}

//定义查询参数(数据量大时,可以做服务端分页)
type QueryHostRequest struct {
	//每页数量
	PageSize int
	//页数
	PageNumber int
}

func (req *QueryHostRequest) Offset() int {
	return (req.PageNumber - 1) * req.PageSize
}

func NewDesribeHostRequestWithID(id string) *DescribeHostRquest {
	return &DescribeHostRquest{
		Id: id,
	}
}
//主机详情查询条件请求
type DescribeHostRquest struct {
	Id string
}

const (
	PUT UpdateMode = 0
	PATCH UpdateMode = 1
)

type UpdateMode int


func NewPatchUpdateHostRequest() *UpdateHostRequest {
	return &UpdateHostRequest{
		UpdateMode: PATCH,
		Resource:   &Resource{},
		Describe:   &Describe{},
	}
}

func NewPutUpdateHostRequest() *UpdateHostRequest {
	return &UpdateHostRequest{
		UpdateMode: PUT,
		Resource:   &Resource{},
		Describe:   &Describe{},
	}
}
//更新数据请求,
type UpdateHostRequest struct {
	UpdateMode
	*Resource
	*Describe
}

//删除数据请求
type DeleteHostRequest struct {
	Id string
}