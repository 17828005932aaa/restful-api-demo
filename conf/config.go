package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/infraboard/mcube/logger/zap"
)

//conf pkg 的全局变量
//全局配置对象
var cf *Config

//全局配置对象的访问方式
func Cf() *Config {
	if cf == nil {
		panic("config required")
	}
	return cf
}

//全局配置对象的设置方式
func SetGlobalConfig(conf *Config) {
	cf = conf
}

func NewDefaultConfig() *Config {
	return &Config{
		App:   newDefaultApp(),
		Mysql: newDefaultMYSQL(),
		Log:   newDefaultLog(),
	}
}

//设置默认值
//app
func newDefaultApp() *app {
	return &app{
		Name: "restful-api",
		Host: "127.0.0.1",
		Port: "3580",
		Key:  "default app key",
	}
}

//mysql
func newDefaultMYSQL() *mysql {
	return &mysql{
		Host:        "192.168.83.145",
		Port:        "3306",
		Username:    "root",
		Password:    "Tcdn@2007",
		Database:    "restful_api",
		MaxOpenConn: 100,
		MaxIdleTime: 20,
		MaxLifeTime: 10 * 60 * 60,
		MaxIdleConn: 5 * 60 * 60,
	}
}

//log
func newDefaultLog() *log {
	return &log{
		Level:  zap.DebugLevel.String(),
		To:     ToStdout,
		Format: TextFormat,
	}
}

//对所有配置对象的整合
type Config struct {
	App   *app
	Mysql *mysql
	Log   *log
}

//配置通过对象来进行映射
//我们定义的是,配置对象的数据结构

//应用程序本身配置
type app struct {
	// restful-api
	Name string
	// 127.0.0.1,0.0.0.0
	Host string `toml:"host"`
	//8080,8050
	Port string `toml:"port"`
	// 比较敏感的数据,入库的是加密后的数据,加密的密钥就是该配置
	Key string `toml:"key"`
}

func (a *app) Addr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

//MYSQL数据库配置
type mysql struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	//最大连接数
	MaxOpenConn int `toml:"max_open_conn"`
	//最多闲置连接数
	MaxIdleConn int `toml:"max_idle_conn"`
	//最大生命周期
	MaxLifeTime int `toml:"max_life_time"`
	//最大闲置连接声明周期
	MaxIdleTime int `toml:"max_idle_time"`

	lock sync.Mutex
}

var (
	db *sql.DB
)

//利用mysql配置,构造全局MySQL单例链接
func (m *mysql) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.Username, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error,%s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//测试mysql连接是否正常
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error,%s", dsn, err.Error())
	}
	return db, nil
}

func (m *mysql) GetDB() (*sql.DB, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		conn, err := m.getDBConn()
		if err != nil {
			return nil, err
		}
		db = conn
	}
	return db, nil
}

type log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}
