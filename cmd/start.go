package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/staryjie/restful-api-demo/apps/host/http"
	"github.com/staryjie/restful-api-demo/apps/host/impl"
	"github.com/staryjie/restful-api-demo/conf"
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
		service := impl.NewHostServiceImpl()

		// 通过Host API Handler对外提供HTTP RESTful API接口
		api := http.NewHostHTTPHandler(service)

		// 提供一个Gin的Router
		g := gin.Default()
		api.Registry(g)

		return g.Run(conf.C().App.HttpAddr())
	},
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo api配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
