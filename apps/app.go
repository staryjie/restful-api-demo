package apps

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/staryjie/restful-api-demo/apps/host"
	"google.golang.org/grpc"
)

// IOC容器层：管理所有的服务实例

// 1. HostService的实例必须注册过来，不然HostService就是一个nil，在服务启动的时候注册
// 2. HTTP暴露模块，依赖IOC中的HostService
var (
	HostService host.Service
	// 如果有很多的Service，那不可能每个Service都要写一遍
	// 通过interface{}加断言进行抽象

	// 维护当前所有的服务
	implApps = map[string]ImplService{}
	ginApps  = map[string]GinService{}
	grpcApps = map[string]GrpcService{}
)

func RegistryImpl(svc ImplService) {
	// 通过服务名判断服务是否已经注册过
	if _, ok := implApps[svc.Name()]; ok {
		panic(fmt.Sprintf("Service %s has registried!", svc.Name()))
	}

	// 服务实力注册到svcs的map中
	implApps[svc.Name()] = svc

	// 根据对象满足的接口来注册具体的服务
	if v, ok := svc.(host.Service); ok {
		HostService = v
	}
}

// 根据名称，返回一个对象，任何类型都可以，具体的使用由使用方通过断言的方式进行判断和使用
// 从implApps中去获取这个指定名称的对象
func GetImpl(name string) interface{} {
	for k, v := range implApps {
		if name == k {
			return v
		}
	}
	return nil
}

// 用于初始化注册到IOC中的所有服务
func InitImpl() {
	for _, v := range grpcApps {
		v.Config()
	}

	for _, v := range implApps {
		v.Config()
	}
}

// 已经加载的完成的Gin App列表
func LoadedGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}
	return
}

// 已经加载的完成的Grpc App列表
func LoadedGrpcApps() (names []string) {
	for k := range grpcApps {
		names = append(names, k)
	}
	return
}

type ImplService interface {
	Config()
	Name() string
}

func RegistryGin(svc GinService) {
	// 通过服务名判断服务是否已经注册过
	if _, ok := ginApps[svc.Name()]; ok {
		panic(fmt.Sprintf("Service %s has registried!", svc.Name()))
	}

	// 服务实力注册到svcs的map中
	ginApps[svc.Name()] = svc
}

func RegistryGrpc(svc GrpcService) {
	// 通过服务名判断服务是否已经注册过
	if _, ok := grpcApps[svc.Name()]; ok {
		panic(fmt.Sprintf("Service %s has registried!", svc.Name()))
	}

	// 服务实力注册到svcs的map中
	grpcApps[svc.Name()] = svc
}

func InitGin(r gin.IRouter) {
	// 先初始化所有的对象
	for _, v := range ginApps {
		v.Config()
	}
	// 完成HTTP Handler注册
	for _, v := range ginApps {
		v.Registry(r)
	}
}

// 注册由Gin编写的Http Handler
// 比如实现了HTTP A，只需要实现Registry()方法，就能把Handler注册给Root Handler
type GinService interface {
	Registry(r gin.IRouter)
	Config()
	Name() string
}

// 初始化grpc服务，将所有实现grpc接口的实体类，注册到grpcServer中
func InitGrpc(r *grpc.Server) {
	// 先初始化所有的对象
	for _, v := range grpcApps {
		v.Config()
	}
	// 完成所有实现grpc接口的实体类注册
	for _, v := range grpcApps {
		v.Registry(r)
	}
}

type GrpcService interface {
	Registry(r *grpc.Server)
	Config()
	Name() string
}
