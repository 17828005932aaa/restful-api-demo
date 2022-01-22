package http

import (
	"restful-api-demo/apps/host"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

//host 模块的HTTP API 服务实例
var API = &handler{}

type handler struct {
	host host.ServiceServer
	log  logger.Logger
}

//初始化时，依赖外部host Service的实例对象
func (h *handler) Config() error {
	h.log = zap.L().Named("HOST API")
	//因为grpc已经将service实现 直接获取grpc获取到的服务实现
	h.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)
	return nil
	// if app.Host == nil {
	// 	panic("dependence host service is nil")
	// }
	// h.host = app.Host

}

func (h *handler) Name() string {
	return host.AppName
}

// 把Handler实现的方法 注册给主路由
func (h *handler) Registry(root router.SubRouter) {
	root.Handle("POST", "/hosts", h.CreateHost)
	root.Handle("GET", "/hosts", h.QueryHost)
	root.Handle("GET", "/hosts/:id", h.DescribeHost)
	root.Handle("PUT", "/hosts/:id", h.UpdateHost)
	root.Handle("PATCH", "/hosts/:id", h.PatchHost)
	root.Handle("DELETE", "/hosts/:id", h.DeleteHost)
}

func init() {
	app.RegistryHttpApp(API)
}
