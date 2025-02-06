package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/DomTok/app/user/controllers"
	"github.com/west2-online/DomTok/app/user/repository/mysql"
	"github.com/west2-online/DomTok/app/user/repository/rpc/template"
	"github.com/west2-online/DomTok/app/user/usecase"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/kitex_gen/user/userservice"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

// 存在的问题
// 1. 对于外部的依赖初始化（例如数据库），如果失败应该直接 panic，我认为 init 阶段不应该返回错误
// 2. 我不推荐 clientset 的用法，这个实在是太高层抽象了。如果是担心 client 太多的话，我建议是使用第三方的 go wire 实现依赖注入, kratos 框架也是默认使用 go wire
var (
	dbAdapter       *mysql.DBAdapter            // 定义 db 的具体实现类
	templateAdapter *template.TemplateRPCClient // 定义 template 的具体实现类
	serviceAdapter  *usecase.UseCase            // usecase 的具体实现类
	serviceName     = "user"
)

func init() {
	config.Init(serviceName)
	logger.Init(serviceName, config.GetLoggerLevel())

	// mysql 的 client 获取, 对应问题 1
	// 注入 mysql 依赖
	dbClient, _ := client.InitMySQL(mysql.User{}.TableName())
	dbAdapter = mysql.NewDBAdapter(dbClient)

	// 注入 use case 依赖
	serviceAdapter = usecase.NewUserCase(dbAdapter, templateAdapter)
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		logger.Fatalf("User: new etcd registry failed, err: %v", err)
	}
	listenAddr, err := utils.GetAvailablePort()
	if err != nil {
		logger.Fatalf("User: get available port failed, err: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("User: resolve tcp addr failed, err: %v", err)
	}
	svr := userservice.NewServer(
		// 注入 controller 依赖
		controllers.NewUserHandler(serviceAdapter),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		server.WithMuxTransport(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),
	)
	if err = svr.Run(); err != nil {
		logger.Fatalf("User: run server failed, err: %v", err)
	}
}
