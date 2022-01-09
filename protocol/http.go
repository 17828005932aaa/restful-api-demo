package protocol

import (
	"context"
	"fmt"
	"net/http"
	"restful-api-demo/conf"
	"time"

	hostAPI "restful-api-demo/app/host/http"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/julienschmidt/httprouter"
)

func NewHTTPService() *HTTPService {
	r := httprouter.New()
	return &HTTPService{
		r: r,
		l: zap.L().Named("HTTP Server"),
		server: &http.Server{
			//HTTP SErver监听地址
			Addr:    conf.Cf().App.Addr(),
			Handler: r,
			//读取header 超时设置
			ReadHeaderTimeout: 60 * time.Second,
			//连接,client ---> server 超时时间
			ReadTimeout: 60 * time.Second,
			//resp 超时时间
			WriteTimeout: 60 * time.Second,
			//http tcp 复用
			IdleTimeout: 60 * time.Second,
			//header大小控制
			MaxHeaderBytes: 1 << 20, // 1M
		},
	}
}

// HTTPService http服务
type HTTPService struct {
	//router, root router,路由,method+path  ---> handler
	r *httprouter.Router
	// 日志
	l logger.Logger
	//服务实例队列,HTTP 服务器
	server *http.Server
}

// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由
	//host http api 服务模块
	hostAPI.API.Init()
	hostAPI.API.Registry(s.r)

	// 启动 HTTP服务
	s.l.Infof("HTTP服务启动成功, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

// Stop 停止server
func (s *HTTPService) Stop() error {
	s.l.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 优雅关闭HTTP服务
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}
	return nil
}
