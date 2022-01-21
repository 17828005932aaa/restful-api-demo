package protocol

import (
	"net"
	"restful-api-demo/conf"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
)

func NewGrpcService() *GrpcService {
	return &GrpcService{
		l: zap.L().Named("GRPC Server"),
		//创建一个grpc服务
		server: grpc.NewServer(),
		//获取配置中的grpc监听地址和端口信息
		GrpcAddr: conf.Cf().App.GrpcAddr(),
		
	}
}

type GrpcService struct {
	// server grpc服务对象
	server *grpc.Server
	// 日志
	l logger.Logger
	//监听地址
	GrpcAddr string
}


func (s *GrpcService) Start()  {
	//加载服务
	app.LoadGrpcApp(s.server)
	// host.RegisterServiceServer(s.server,app.Host)
	//监听GRPC端口
	ls,err:=net.Listen("tcp",s.GrpcAddr)
	if err != nil {
		s.l.Error("start grpc service error, %s",err.Error())
		return
	}
	s.l.Info("GRPC 服务监听地址:",s.GrpcAddr)
	//把监听端口加载到grpc服务中
	if err := s.server.Serve(ls);err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}
		s.l.Error("start grpc service error, %s",err.Error())
		return 
	}
}

func (s *GrpcService) Stop() {
	//优雅退出
	s.server.GracefulStop()
}
