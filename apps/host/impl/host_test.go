package impl_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/restful-api-demo/apps/host"
	"github.com/staryjie/restful-api-demo/apps/host/impl"
)

var (
	// 定义对象是满足该接口的实例
	service host.Service
)

func TestCreate(t *testing.T) {
	ins := host.NewHost()
	ins.Resource.Name = "test"
	service.CreateHost(context.Background(), ins)
}

func init() {
	// 初始化全局logger
	// 处于性能考虑，设计为默认不打印
	// 在初始化service之前需要先初始化全局logger实例
	zap.DevelopmentSetup()

	// 接口的具体实现
	service = impl.NewHostServiceImpl()
}
