package impl

import (
	"context"

	"github.com/infraboard/mcube/logger"
	"github.com/staryjie/restful-api-demo/apps/host"
)

// 录入主机
func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (
	*host.Host, error) {
	// 直接打印日志
	i.l.Debug("Create Host")

	// 带格式化的日志打印
	i.l.Debugf("Create Host %s", ins.Name)

	// 携带额外的meta数据，常用于Trace系统
	i.l.With(logger.NewAny("request-id", "req01")).Debug("Create Host with meta kv")
	return ins, nil
}

// 查询主机列表
func (i *HostServiceImpl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (
	*host.HostSet, error) {
	return nil, nil
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
