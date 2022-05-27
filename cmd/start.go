package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/conf"

	"github.com/staryjie/restful-api-demo/apps/host/http"

	// 注册所有的服务实例
	_ "github.com/staryjie/restful-api-demo/apps/all"
)

var (
	// pusher service config option
	confType string
	confFile string
	confETCD string
)

// 程序的启动时 组装都在这里进行
// StartCmd represents the base command when called without any subcommands
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动 demo 后端API",
	Long:  "启动 demo 后端API",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置文件，初始化全局变量
		err := conf.LoadConfigFromToml(confFile)
		// err := conf.LoadConfigFromEnv()
		if err != nil {
			panic(err)
		}

		// 加载Host Service实体类
		// service := impl.NewHostServiceImpl()

		// 注册HostService的实例到IOC中
		// apps.HostService = impl.NewHostServiceImpl()
		// 采用：_ "github.com/staryjie/restful-api-demo/apps/host/impl" 完成注册

		// apps.HostService 是一个host.Service的接口，并没有实例初始化(Config)的方法
		apps.Init()

		// 通过Host API Handler对外提供HTTP RESTful API接口
		// api := http.NewHostHTTPHandler(service)
		api := http.NewHostHTTPHandler()
		// 从IOC中获取依赖，解除了相互依赖关系
		api.Config()

		// 提供一个Gin的Router
		g := gin.Default()
		api.Registry(g)

		return g.Run(conf.C().App.HttpAddr())
	},
}

// 还没有初始化Logger实例

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo api配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
