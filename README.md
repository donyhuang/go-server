# Server



## golang 框架
支持http grpc kafka consumer

封装了对于日志，配置文件，nacos远程配置，mysql(gorm) clickhouse



## 示例配置文件 configs/app.yml

## 开启一个http server
go get gitlab.nongchangshijie.com/go-base/server/pkg/server
```golang
//开启一个http server

routes := route.NewRoutes()
routes.Add(
    []route.Route{
    {
    Method: "POST",
    Path:   "spread",
    Parent: []string{backEndGroup, AdvertisePrefix},
    Handle: []gin.HandlerFunc{ad.Create},
    Desc:   "广告计划创建",
    },
    {
    Method: "GET",
    Path:   "spread/list",
    Parent: []string{backEndGroup, AdvertisePrefix},
    Handle: []gin.HandlerFunc{ad.List},
    Desc:   "广告计划查看",
    },}...
)

svr := server.NewServer("./configs/backend.yml")
svr.RegisterHttpRoutes("advertisement-spread", route.GinRoutes())
svr.Start()



//自己的工作协程需要优雅退出
tick := time.NewTicker(time.Minute)
doneChan := make(chan struct{})
appEvent := event.NewEvent(svr.GetContext(), doneChan, tick, db.GetDb(), db.GetClickHouseConn())
appEvent.Run()
proc.AddDoneFn(appEvent.Stop) //在收到信号sigterm|SIGINT之后等待自己的协程优雅关闭
```
## 开启一个http server 以及grpc server 同时等待自己的工作协程结束
```golang
svr := server.NewServer("./configs/backend.yml")
svr.RegisterHttpRoutes("advertisement-platform", route.EngineRoute(string(svr.GetGlobalConf().Server.Mode))
handle := server2.NewGrpcHandler(proto.Order_ServiceDesc, &ErpOrder{})
svr.RegisterGrpcHandle("advertisement-order", handle)
svr.Start()


```
