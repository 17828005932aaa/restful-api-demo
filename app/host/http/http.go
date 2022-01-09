package http

import (
	"restful-api-demo/app"
	"restful-api-demo/app/host"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/julienschmidt/httprouter"
)

//host 模块的HTTP API 服务实例
var API = handler{}

type handler struct {
	host host.Service
	log logger.Logger
}

//初始化时，依赖外部host Service的实例对象
func (h *handler) Init() {
	h.log = zap.L().Named("HOST API")
	if app.Host == nil {
		panic("dependence host service is nil")
	}
	h.host = app.Host
	
}

// 把Handler实现的方法 注册给主路由
func (h *handler) Registry(root *httprouter.Router){
	root.POST("/hosts",h.CreateHost)
	root.GET("/hosts",h.QueryHost)
	root.GET("/hosts/:id",h.DescribeHost)
	root.PUT("/hosts/:id", h.UpdateHost)
	root.PATCH("/hosts/:id", h.PatchHost)
	root.DELETE("/hosts/:id",h.DeleteHost)
}