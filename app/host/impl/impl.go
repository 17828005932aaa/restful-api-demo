package impl

import (
	"database/sql"
	"restful-api-demo/conf"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

var Service *impl = &impl{}



type impl struct {
	//定义日志属性,可以更换成熟悉的日志库,比如logrus,标准库log, zap
    //mcube log模块是包装的zap的实现
	log logger.Logger	

	//依赖数据库
	db *sql.DB

}

func (i *impl) Init() error {
	i.log = zap.L().Named("Host")

	db,err:=conf.Cf().Mysql.GetDB()
	if err != nil {
		return err
	}
	i.db = db
	return nil
}