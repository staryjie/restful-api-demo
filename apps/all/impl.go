package all

// 完成所有模块的注册
import (
	_ "github.com/staryjie/restful-api-demo/apps/book/api"
	_ "github.com/staryjie/restful-api-demo/apps/book/impl"
	_ "github.com/staryjie/restful-api-demo/apps/host/http"
	_ "github.com/staryjie/restful-api-demo/apps/host/impl"
)
