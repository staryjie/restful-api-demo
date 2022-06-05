package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 全局Config实例对象
// 也就是城西在内存中的配置对象
// 程序内部获取配置都通过读取该对象
// Condig对象什么时候被初始化
//    配置加载时：
//        1. LoadConfigFromToml
//        2. LoadConfigFromEnv
// 为了不被程序在运行时进行修改，设置为私有变量
var config *Config

// 全局MySQL客户端实例
var db *sql.DB

// 获取到配置，单独提供函数
// 全局Config对象获取函数
func C() *Config {
	return config
}

// Config 应用配置
type Config struct {
	App   *App   `toml:"app"`
	Log   *Log   `toml:"log"`
	MySQL *MySQL `toml:"mysql"`
}

// 初始化一个有默认值的Config对象
func NewDefaultConfig() *Config {
	return &Config{
		App:   NewDefaultApp(),
		Log:   NewDefaultLog(),
		MySQL: NewDefaultMySQL(),
	}
}

type App struct {
	Name string `toml:"name" env:"APP_NAME"`
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
	// Key  string `toml:"key" env:"APP_KEY"`
}

func NewDefaultApp() *App {
	return &App{
		Name: "demo",
		Host: "127.0.0.1",
		Port: "8050",
	}
}

func (a *App) HttpAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func (a *App) GrpcAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, fmt.Sprintf("1%s", a.Port))
}

func (a *App) RestfulAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, fmt.Sprintf("2%s", a.Port))
}

// MySQL todo
type MySQL struct {
	Host     string `toml:"host" env:"MYSQL_HOST"`
	Port     string `toml:"port" env:"MYSQL_PORT"`
	UserName string `toml:"username" env:"MYSQL_USERNAME"`
	Password string `toml:"password" env:"MYSQL_PASSWORD"`
	Database string `toml:"database" env:"MYSQL_DATABASE"`
	// 因为使用的是MySQL的连接池，需要对连接池做一些规划配置
	// 控制当前程序打开的MySQL连接数
	MaxOpenConn int `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	// 控制MySQL连接的复用，比如 5 最多允许5个复用
	MaxIdleConn int `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	// 控制MySQL的最大连接生命周期，和MySQL Server的配置有关，必须小于等于Server端的配置
	// 一个连接使用12h，到期后必须换一个连接，保证一定的可用性
	MaxLifeTime int `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	// Idle连接最多允许存活时间
	MaxIdleTime int `toml:"max_idle_time" env:"MYSQL_MAX_idle_TIME"`

	// 作为私有变量，用于控制getDBConn
	lock sync.Mutex
}

func NewDefaultMySQL() *MySQL {
	return &MySQL{
		Host:        "127.0.0.1",
		Port:        "3306",
		UserName:    "demo",
		Password:    "123456",
		MaxOpenConn: 200,
		MaxIdleConn: 100,
	}
}

// 连接池，driverConn具体的连接对象，它维护这一个Socket
// pool []*driverConn，维护pool里面的连接都是可用的，定期检查conn的健康状况
// 当某一个连接失效了，它会清空该连接结构体数据(driverConn.Reset())，重新建立一个连接(Reconn)，让该conn借壳存活
// 避免driverConn结构体申请和释放的成本
func (m *MySQL) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

// 1.使用LoadGlobal在加载时初始化全局db实例
// 2.惰性加载，每次都获取db，动态判断是否需要初始化
func (m *MySQL) GetDB() *sql.DB {
	// 如果全局实例不存在，会报错

	// 直接加锁，锁住临界区
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		// 实例不存在，加载一个新的实例
		conn, err := m.getDBConn()
		if err != nil {
			panic(err)
		}
		db = conn
	}

	// 全局实例db一定存在
	return db
}

// Log todo
type Log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}

func NewDefaultLog() *Log {
	return &Log{
		Level:  "info",
		Format: TextFormat,
		To:     ToStdout,
	}
}
