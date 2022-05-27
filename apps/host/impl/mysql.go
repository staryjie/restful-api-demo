package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/apps/host"
	"github.com/staryjie/restful-api-demo/conf"
)

// 接口实现的静态检查
// var _ host.Service = (*HostServiceImpl)(nil)
// var impl = NewHostServiceImpl()  // 会导致conf.C()对象并没有准备好，导致conf.C().MySQL.GetDB()发生panic

// 将对象的注册和初始化独立开
var impl = &HostServiceImpl{}

type HostServiceImpl struct {
	l  logger.Logger
	db *sql.DB
}

// 保证调用该函数之前全局conf对象已经初始化
func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		// Host Service的子logger
		// 封装Zap让其满足Logger接口
		// 为什么要封装：
		//   1. Logger全局实例
		//   2. Logger Level的动态调整，Logrus不支持Level的动态调整
		//   3. 加入日志轮转功能的集合
		l:  zap.L().Named("Host"),
		db: conf.C().MySQL.GetDB(),
	}
}

// 需要确保全局logger和全局conf对象已经加载完成
func (i *HostServiceImpl) Config() {
	i.l = zap.L().Named("Host")
	i.db = conf.C().MySQL.GetDB()
}

func (i *HostServiceImpl) Name() string {
	return host.AppName
}

// 都需要在启动的时候手动将服务注册到IOC层
// 注册HostService的实例到IOC中
// apps.HostService = impl.NewHostServiceImpl()

// mysql的驱动加载实现方式
// _ "github.com/go-sql-driver/mysql"
// 通过利用模块的init()方法，实现导入即可执行init()方法的特性来实现自动注册到IOC

// _ import app 实现自动注册到IOC
func init() {
	// 对象注册到IOC层
	apps.Registry(impl)
	// apps.HostService = impl
}
