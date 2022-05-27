package apps

import "github.com/staryjie/restful-api-demo/apps/host"

// IOC容器层：管理所有的服务实例

// 1. HostService的实例必须注册过来，不然HostService就是一个nil，在服务启动的时候注册
// 2. HTTP暴露模块，依赖IOC中的HostService
var (
	HostService host.Service
)
