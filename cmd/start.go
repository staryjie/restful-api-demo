package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"

	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/conf"

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

		// 初始化全局日志Logger
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		// 加载Host Service实体类
		// service := impl.NewHostServiceImpl()

		// 注册HostService的实例到IOC中
		// apps.HostService = impl.NewHostServiceImpl()
		// 采用：_ "github.com/staryjie/restful-api-demo/apps/host/impl" 完成注册

		// apps.HostService 是一个host.Service的接口，并没有实例初始化(Config)的方法
		apps.InitImpl()

		// 通过Host API Handler对外提供HTTP RESTful API接口
		// api := http.NewHostHTTPHandler(service)
		// api := http.NewHostHTTPHandler()
		// 从IOC中获取依赖，解除了相互依赖关系
		// api.Config()

		// 提供一个Gin的Router
		g := gin.Default()

		// 注册所有HTTP Handler到IOC中
		apps.InitGin(g)
		// api.Registry(g)

		return g.Run(conf.C().App.HttpAddr())
	},
}

// 问题：
//    1. HTTP API， GRPC API需要启动，消息总线也需要监听，比如负责注册配置，这些模块都是独立的
//       都需要在程序启动的时候进行初始化和启动，如果都写在start中，那么start会很臃肿，难以维护
//    2. 服务的优雅关闭怎么实现？ 外部会发生一个Terminalied（中断）信号给程序，程序需要监听并处理这些信号
//       需要实现程序优雅关闭的处理逻辑：从外到内完成资源的释放和关闭处理
//         1. API层的关闭
//         2. 消息总线等关闭
//         3. 关闭数据库连接
//         4. 如果使用了注册中心，最后需要在注册中心完成服务下线操作
//         5. 程序退出

// 还没有初始化Logger实例
// log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)

	// 根据配置文件加载后的conf对象来配置全局logger对象
	lc := conf.C().Log

	// 设置日志级别
	// DebugLevel: "debug",
	// InfoLevel:  "info",
	// WarnLevel:  "warning",
	// ErrorLevel: "error",
	// FatalLevel: "fatal",
	// PanicLevel: "panic",
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		// 配置出错则使用默认日志级别info
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}

	// 使用默认配置初始化logger全局配置
	zapConfig := zap.DefaultConfig()

	// 加上用户自定义配置完成全局logger配置
	zapConfig.Level = level

	// 配置程序是否在启动的时候创建新的日志文件
	zapConfig.Files.RotateOnStartup = false

	// 配置日志输出
	switch lc.To {
	case conf.ToStdout:
		// 日志打印到标准输出
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		// 日志打印到文件
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}

	// 配置日志输出格式
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}

	// 将配置应用到全局logger对象
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo api配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
