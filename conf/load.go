package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

//环境变量中加载
func LoadConfigFromToml(path string) error {
	//new配置对象
	cfg := NewDefaultConfig()
	//解析配置文件赋值给cfg,没有赋值的取默认值
	_,err:=toml.DecodeFile(path,cfg)
	if err != nil {
		return err
	}
	SetGlobalConfig(cfg)
	return nil
}

//配置文件中加载
func LoadConfigFromEnv() error {
	//new配置对象
	cfg := NewDefaultConfig()
	//解析配置文件赋值给cfg,没有赋值的取默认值
	if err:=env.Parse(cfg); err != nil {
		return err
	}
	SetGlobalConfig(cfg)
	return nil
}