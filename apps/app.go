package apps

import (
	"fmt"

	"github.com/staryjie/restful-api-demo/apps/host"
)

// IOC容器层：管理所有的服务实例

// 1. HostService的实例必须注册过来，不然HostService就是一个nil，在服务启动的时候注册
// 2. HTTP暴露模块，依赖IOC中的HostService
var (
	HostService host.Service
	svcs        = map[string]Service{}
)

func Registry(svc Service) {
	// 通过服务名判断服务是否已经注册过
	if _, ok := svcs[svc.Name()]; ok {
		panic(fmt.Sprintf("Service %s has registried!", svc.Name()))
	}

	// 服务实力注册到svcs的map中
	svcs[svc.Name()] = svc

	// 根据对象满足的接口来注册具体的服务
	if v, ok := svc.(host.Service); ok {
		HostService = v
	}
}

// 用于初始化注册到IOC中的所有服务
func Init() {
	for _, v := range svcs {
		v.Config()
	}
}

type Service interface {
	Config()
	Name() string
}
