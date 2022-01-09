package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"restful-api-demo/app"
	"restful-api-demo/app/host/impl"
	"restful-api-demo/conf"
	"restful-api-demo/protocol"
	"syscall"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"
)

var (
	configType string
	confFile   string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Demo后端API服务",
	Long:  `Demo后端API服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//配置加载
		if err := loadGlobalConfig(configType); err != nil {
			return err
		}

		//初始化日志
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		//初始化服务层 Ioc初始化
		if err := impl.Service.Init(); err != nil {
			return err
		}
		//初始化实例注册给IOC层
		app.Host = impl.Service

		//启动服务后，需要处理的事件
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		//启动服务
		svr := NewService(*conf.Cf())

		//等待程序退出
		go svr.waitSign(ch)

		//启动服务
		return svr.Start()
	},
}

func NewService(conf conf.Config) *Service {
	return &Service{
		conf: conf,
		http: protocol.NewHTTPService(),
		log:  zap.L().Named("Service"),
	}
}

//service
//服务的整体配置
//服务可能会启动很多模块,http,grpc,crontable
type Service struct {
	conf conf.Config
	http *protocol.HTTPService
	log  logger.Logger
}

func (s *Service) Start() error {
	return s.http.Start()
}

// 当发现用户收到终止掉程序的时候，需要完成处理
func (s *Service) waitSign(sign chan os.Signal) {
	for sg := range sign {
		switch v := sg.(type) {
		default:
			// 资源清理
			s.log.Infof("receive signal '%v', start graceful shutdown", v.String())
			if err := s.http.Stop(); err != nil {
				s.log.Errorf("graceful shutdown err: %s, force exit", err)
			}
			s.log.Infof("service stop complete")
			return
		}
	}
}

//config为全局变量,只需要load即可全局可用户
func loadGlobalConfig(configType string) error {
	switch configType {
	case "file":
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			return err
		}
	case "env":
		err := conf.LoadConfigFromEnv()
		if err != nil {
			return err
		}
	case "etcd":
		return errors.New("not implemented")
	default:
		return errors.New("unknown config type")
	}
	return nil
}

// log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)
	//获取出日志配置对象
	lc := conf.Cf().Log
	//解析配置的日志级别是否正确
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		//解析失败,默认使用info
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		//解析成功,直接使用用户配置的日志
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}
	//初始化了日志的默认配置
	zapConfig := zap.DefaultConfig()
	zapConfig.Level = level
	zapConfig.Files.RotateOnStartup = false
	switch lc.To {
	case conf.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "restful-api.log"
		zapConfig.Files.Path = lc.PathDir
	}
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}
	//初始化全局logger的配置
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	//全局Logger 初始化后 就可以正常使用
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	//全局标志PersistentFlags 所有的子节点都可以使用  Flags仅本节点可用
	//第一个形参是cmd 传参传给哪个变量接受，第二个代表 该 参数的名称，第三个代表 参数指令，第四个 代表参数的默认值，第五个代表参数的帮助描述
	RootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&configType, "config_type", "t", "file", "the restful-api config type")
	startCmd.PersistentFlags().StringVarP(&confFile, "config_file", "f", "etc/restful-api.toml", "the restful-api config file path")

}
