package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"

	"github.com/staryjie/restful-api-demo/apps"
	"github.com/staryjie/restful-api-demo/conf"
	"github.com/staryjie/restful-api-demo/protocol"

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
		// g := gin.Default()

		// 注册所有HTTP Handler到IOC中
		// apps.InitGin(g)
		// api.Registry(g)

		svc := NewManager()
		ch := make(chan os.Signal, 1)
		// channel是一种复合数据结构, 可以当初一个容器, 自定义的struct make(chan int, 1000), 8bytes * 1024  1Kb
		// 如果没close gc是不会回收的
		defer close(ch)

		// Go为了并发编程设计的(CSP), 依赖Channel作为数据通信的信道
		// 出现了一个思路模式的转变:
		//    单兵作战(只有一个Groutine) --> 团队作战(多个Groutine 采用Channel来通信)
		//    main { for range channel }  这个时候的channel仅仅想到于一个缓存, 必须选择带缓存区的channl
		//    signal.Notify 当中一个Goroutine, g1
		//    go svc.WaitStop(ch) 第二Goroutine, g2
		//    g1 -- ch1 --> g2
		//    g1 <-- ch2 -- g2
		//    g1 数据发送给ch1, g2 读取channle中的数据, chanel 只要生成好了就能用, 如果channle关闭
		//    设计channel 使用数据的发送方负责关闭, 相当于表示挂电话
		//    for range   由range帮忙处理了 chnanel 关闭后， read的中断处理
		//    for v,err := <-ch { if(err == io.EOF) break }
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		go svc.WaitStop(ch)

		// 后台启动grpc Service
		go svc.grpc.Start()

		// 后台启动 restful api
		go svc.rest.Start()

		return svc.Start()

		// return g.Run(conf.C().App.HttpAddr())
	},
}

// 用于管理所有需要启动的服务
// 1.HTTP服务的启动
// 2. Grpc服务的启动
type manager struct {
	rest *protocol.RestfulService
	http *protocol.HttpService
	grpc *protocol.GRPCService
	l    logger.Logger
}

// 有两个服务需要启动,一个http, 一个grpc
// 一个前台启动，一个后台启动
func NewManager() *manager {
	return &manager{
		rest: protocol.NewRestfulService(),
		http: protocol.NewHttpService(),
		grpc: protocol.NewGRPCService(),
		l:    zap.L().Named("CLI"),
	}
}

func (m *manager) Start() error {
	return m.http.Start()
}

// 处理来自外部的中断信号, 比如Terminaled
func (m *manager) WaitStop(ch <-chan os.Signal) {
	for v := range ch {
		switch v {
		// case syscall.SIGTERM:
		// case syscall.SIGHUP:

		// 统一做停止服务处理
		default:
			// 先关闭内部调用
			// 关闭gRpc
			if err := m.grpc.Stop(); err != nil {
				m.l.Error(err)
			}
			// 关闭restful
			m.rest.Stop()
			// 关闭外部调用
			m.l.Infof("Received signal: %s, start stop server ...", v)
			m.http.Stop()
		}
	}
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
