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
	ins.Resource.Id = "ins-01"
	ins.Resource.Region = "cn-hangzhou"
	ins.Resource.Type = "sm1"
	ins.Describe.CPU = 1
	ins.Describe.Memory = 2048

	should := assert.New(t)
	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Printf("%#v", ins)
	}
}

func TestQueryHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewQueryHostRequest()
	req.KeyWords = "postman"
	set, err := service.QueryHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Printf("Total: %d\n", set.Total)

		for i := range set.Items {
			fmt.Println(set.Items[i].Id)
		}
	}
}

func TestDescribeHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewDescribeHostRequestWithId("ins-09")
	ins, err := service.DescribeHost(context.Background(), req)
	if should.NoError(err) {

		fmt.Println(ins.Name)
	}
}

func TestUpdateHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewPutUpdateRequest("ins-07")
	req.Name = "test-put-07"
	req.Region = "cn-shanghai"
	req.Type = "sm2"
	req.CPU = 4
	req.Memory = 4096
	req.Description = "测试用例-PUT"
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Name)
	}
}

func TestPatchHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewPatchUpdateRequest("ins-07")
	// req.Name = "test-patch-07"
	// req.Region = "cn-shanghai"
	// req.Type = "sm1"
	// req.CPU = 2
	// req.Memory = 2048
	req.Description = "测试用例-PATCH-1212121"

	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Name)
	}
}

func TestDeleteHost(t *testing.T) {
	should := assert.New(t)

	req := host.NewDeleteHostRequestWithId("ins-08")
	ins, err := service.DeleteHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Name)
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
