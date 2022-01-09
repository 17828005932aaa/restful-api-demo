package impl

import (
	"context"
	"database/sql"
	"fmt"
	"restful-api-demo/app/host"
	"time"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/rs/xid"
)

func (i *impl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	// 校验数据合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 生成UUID的一个库,
	// snow
	// 分布式ID, app, instance, ip, mac, ......, idc(), region
	ins.Id = xid.New().String()
	if ins.CreateAt == 0 {
		ins.CreateAt = ftime.Now().Timestamp()
	}

	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败, 事务的逻辑

	// 全局异常
	var (
		resStmt  *sql.Stmt
		descStmt *sql.Stmt
		err      error
	)

	// 初始化一个事务, 所有的操作都使用这个事务来进行提交
	// 比如 用户http 请求取消了, 但是操作数据的逻辑 并不知道
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 函数执行完成后, 专门判断事务是否正常
	defer func() {
		// 事务执行有异常
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				i.log.Debugf("tx rollback error, %s", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				i.log.Debugf("tx commit error, %s", err)
			}
		}
	}()

	// 需要判断事务执行过程当中是否有异常
	// 有异常 就回滚事务, 无异常就提交事务

	// 在这个事务里面执行 Insert SQL, 先执行Prepare, 防止SQL注入攻击
	resStmt, err = tx.Prepare(insertResourceSQL)
	if err != nil {
		return nil, fmt.Errorf("prepare resource sql error, %s", err)
	}
	defer resStmt.Close()

	// 注意: Prepare语句 会占用MySQL资源, 如果你使用后不关闭会导致Prepare溢出
	_, err = resStmt.Exec(
		ins.Id, ins.Vendor, ins.Region, ins.Zone, ins.CreateAt, ins.ExpireAt, ins.Category, ins.Type, ins.InstanceId,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.SyncAccount, ins.PublicIP,
		ins.PrivateIP, ins.PayType, ins.Resource_hash, ins.Describe_hash,
	)
	if err != nil {
		return nil, fmt.Errorf("insert resource error, %s", err)
	}

	// 同样的逻辑,  我们也需要Host的数据存入
	descStmt, err = tx.Prepare(insertDescribeSQL)
	if err != nil {
		return nil, fmt.Errorf("prepare describe sql error, %s", err)
	}
	defer descStmt.Close()

	_, err = descStmt.Exec(
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec, ins.OSType, ins.OSName,
		ins.SerialNumber, ins.ImageID, ins.InternetMaxBandwidthOut,
		ins.InternetMaxBandwidthIn, ins.KeyPairName, ins.SecurityGroups,
	)
	if err != nil {
		return nil, fmt.Errorf("insert describe error, %s", err)
	}

	return ins, nil
}

func (i *impl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.HostSet, error) {

	query := sqlbuilder.NewQuery(queryHostSQL).Order("create_at").Desc().Limit(int64(req.Offset()), uint(req.PageSize))
	//build一个查询语句
	sqlStr, args := query.BuildQuery()
	i.log.Debug("sql:%s,args:%v", sqlStr, args)

	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	//初始化返回的Set
	set := host.NewDefaultHostSet()

	//迭代查询表里的数据
	for rows.Next() {
		ins := host.NewDefaultHost()
		if err := rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt,
			&ins.Category, &ins.Type, &ins.InstanceId, &ins.Name,
			&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount,
			&ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.Resource_hash, &ins.Describe_hash,
			&ins.Id, &ins.CPU,
			&ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName,
			&ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn,
			&ins.KeyPairName, &ins.SecurityGroups,
		); err != nil {
			return nil, err
		}
		set.Add(ins)

	}

	countStr, countArgs := query.BuildCount()
	countStmt, err := i.db.Prepare(countStr)
	if err != nil {
		return nil, fmt.Errorf("prepare count stmt error,%s", err)
	}
	defer countStmt.Close()
	//返回一行
	if err := countStmt.QueryRow(countArgs...).Scan(&set.Total); err != nil {
		return nil, fmt.Errorf("query count error,%s", err)
	}
	return set, nil

}

func (i *impl) DescribeHost(ctx context.Context, req *host.DescribeHostRquest) (*host.Host, error) {
	query := sqlbuilder.NewQuery(queryHostSQL).Where("r.id= ?", req.Id)

	//build 查询语句
	sqlStr, args := query.BuildQuery()
	i.log.Debugf("sql: %s, args: %v", sqlStr, args)

	//Prepare
	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("prepare query host sql error %s", err)
	}
	defer stmt.Close()
	ins := host.NewDefaultHost()
	err = stmt.QueryRow(args...).Scan(
		&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt,
		&ins.Category, &ins.Type, &ins.InstanceId, &ins.Name,
		&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount,
		&ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.Resource_hash, &ins.Describe_hash,
		&ins.Id, &ins.CPU,
		&ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName,
		&ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn,
		&ins.KeyPairName, &ins.SecurityGroups,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("host %s not found", req.Id)
		}
		return nil, fmt.Errorf("stmt query error, %s", err)
	}
	return ins, nil
}

func (i *impl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {

	//重新查询出来
	ins, err := i.DescribeHost(ctx, host.NewDesribeHostRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}
	//对象更新(PATCH/PUT)
	switch req.UpdateMode {
	case host.PUT:
		//全量更新
		ins.Update(req.Resource, req.Describe)
		// 校验数据合法性
		// if err := ins.Validate(); err != nil {
		// 	return nil, err
		// }
	case host.PATCH:
		//部分更新
		err := ins.Patch(req.Resource, req.Describe)
		if err != nil {
			return nil, err
		}
	}
	//Prepare
	stmt, err := i.db.Prepare(updateResourceSQL)
	if err != nil {
		i.log.Debugf("peapare %s", err)
		return nil, err
	}
	defer stmt.Close()

	// DML
	// vendor=?,region=?,zone=?,expire_at=?,name=?,description=? WHERE id = ?
	ins.UpdateAt = time.Now().UnixNano() / 1000000
	_, err = stmt.Exec(ins.Vendor, ins.Region, ins.Zone, ins.ExpireAt, ins.Name, ins.Description, ins.UpdateAt, ins.Id)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *impl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {


	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败, 事务的逻辑

	// 全局异常
	var (
		resStmt  *sql.Stmt
		descStmt *sql.Stmt
		err      error
	)

	//重新查询出来
	ins,err := i.DescribeHost(ctx,host.NewDesribeHostRequestWithID(req.Id))
	if err != nil {
		return nil,err
	}

	// 初始化一个事务, 所有的操作都使用这个事务来进行提交
	// 比如 用户http 请求取消了, 但是操作数据的逻辑 并不知道
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 函数执行完成后, 专门判断事务是否正常
	defer func() {
		// 事务执行有异常
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				i.log.Debugf("tx rollback error, %s", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				i.log.Debugf("tx commit error, %s", err)
			}
		}
	}()

	// 需要判断事务执行过程当中是否有异常
	// 有异常 就回滚事务, 无异常就提交事务

	// 在这个事务里面执行 Insert SQL, 先执行Prepare, 防止SQL注入攻击
	resStmt, err = tx.Prepare(deleteResourceSQL)
	if err != nil {
		return nil, fmt.Errorf("prepare resource sql error, %s", err)
	}
	defer resStmt.Close()

	// 注意: Prepare语句 会占用MySQL资源, 如果你使用后不关闭会导致Prepare溢出
	_, err = resStmt.Exec(req.Id)
	if err != nil {
		return nil, err
	}

	// 同样的逻辑,  我们也需要Host的数据存入
	descStmt, err = tx.Prepare(deleleHostSQL)
	if err != nil {
		return nil, err
	}
	defer descStmt.Close()

	_, err = descStmt.Exec(req.Id)
	if err != nil {
		return nil, err
	}

	return ins, nil
}
