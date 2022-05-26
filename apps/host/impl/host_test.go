package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/restful-api-demo/apps/host"
	"github.com/staryjie/restful-api-demo/apps/host/impl"
	"github.com/staryjie/restful-api-demo/conf"
	"github.com/stretchr/testify/assert"
)

var (
	// 定义对象是满足该接口的实例
	service host.Service
)

func TestCreate(t *testing.T) {
	ins := host.NewHost()
	ins.Resource.Name = "test"

	should := assert.New(t)
	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}
}

func init() {
	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	// err := conf.LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}
	// 初始化全局logger
	// 处于性能考虑，设计为默认不打印
	// 在初始化service之前需要先初始化全局logger实例
	zap.DevelopmentSetup()

	// 接口的具体实现
	service = impl.NewHostServiceImpl()
}
