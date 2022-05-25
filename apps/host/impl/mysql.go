package impl

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/restful-api-demo/apps/host"
)

// 接口实现的静态检查
var _ host.Service = (*HostServiceImpl)(nil)

type HostServiceImpl struct {
	l logger.Logger
}

func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		// Host Service的子logger
		// 封装Zap让其满足Logger接口
		// 为什么要封装：
		//   1. Logger全局实例
		//   2. Logger Level的动态调整，Logrus不支持Level的动态调整
		//   3. 加入日志轮转功能的集合
		l: zap.L().Named("Host"),
	}
}
