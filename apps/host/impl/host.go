package impl

import (
	"context"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/staryjie/restful-api-demo/apps/host"
)

// 录入主机
func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (
	*host.Host, error) {
	// 直接打印日志
	i.l.Named("Create").Debug("Create Host") // 日志分层
	i.l.Info("Create Host")

	// 带格式化的日志打印
	i.l.Debugf("Create Host %s", ins.Name)

	// 携带额外的meta数据，常用于Trace系统
	i.l.With(logger.NewAny("request-id", "req01")).Debug("Create Host with meta kv")

	// 校验数据合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 由dao层负责将对象存储到数据库
	if err := i.save(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

// 查询主机列表
func (i *HostServiceImpl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (
	*host.HostSet, error) {
	b := sqlbuilder.NewBuilder(QueryHostSQL)
	if req.KeyWords != "" {
		b.Where("r.name LIKE ? OR r.description LIKE ? OR r.private_ip LIKE ? OR r.public_ip LIKE ?",
			"%"+req.KeyWords+"%", // name
			"%"+req.KeyWords+"%", // description
			req.KeyWords+"%",     // private_ip
			req.KeyWords+"%",     // public_ip
		)
	}

	// 分页
	b.Limit(req.OffSet(), req.GetPageSize())

	querySQL, args := b.Build()
	i.l.Debugf("query sql : %s, args: %v", querySQL, args)

	// query stmt，构建一个Prepare语句
	stmt, err := i.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...) // 传入参数
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set := host.NewHostSet()
	// 遍历结果
	for rows.Next() {
		// 每扫描一行，就需要将数据读取出来
		// h.cpu, h.memory, h.gpu_spec, h.gpu_amount, h.os_type, h.os_name, h.serial_number
		ins := host.NewHost()
		if err := rows.Scan(&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt,
			&ins.ExpireAt, &ins.Type, &ins.Name, &ins.Description, &ins.Status,
			&ins.UpdateAt, &ins.SyncAt, &ins.Account, &ins.PublicIP, &ins.PrivateIP,
			&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType,
			&ins.OSName, &ins.SerialNumber,
		); err != nil {
			return nil, err
		}
		set.Add(ins)
		// i.l.Debugf("%s", ins.Name)
	}

	// Total统计
	countSQl, args := b.BuildCount()
	i.l.Debugf("count sql: %s, args: %v", countSQl, args)
	countStmt, err := i.db.PrepareContext(ctx, countSQl)
	if err != nil {
		return nil, err
	}
	defer countStmt.Close()

	if err := countStmt.QueryRowContext(ctx, args...).Scan(&set.Total); err != nil {
		return nil, err
	}

	return set, nil
}

// 查询主机详情
func (i *HostServiceImpl) DescribeHost(ctx context.Context, req *host.QueryHostRequest) (
	*host.Host, error) {
	return nil, nil
}

// 更新主机信息
func (i *HostServiceImpl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (
	*host.Host, error) {
	return nil, nil
}

// 删除主机,前端需要展示被删除的主机信息，所以需要返回当前删除的主机信息
func (i *HostServiceImpl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (
	*host.Host, error) {
	return nil, nil
}
