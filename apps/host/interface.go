package host

import "context"

// Host app service 的接口实现
type Service interface {
	// 录入主机
	CreateHost(context.Context, *Host) (*Host, error)
	// 查询主机列表
	QueryHost(context.Context, *QueryHostRequest) (*HostSet, error)
	// 查询主机详情
	DescribeHost(context.Context, *DescribeHostRequest) (*Host, error)
	// 更新主机信息
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	// 删除主机,前端需要展示被删除的主机信息，所以需要返回当前删除的主机信息
	DeleteHost(context.Context, *DeleteHostRequest) (*Host, error)
}
