package server

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/donyhuang/go-server/pkg/cache"
	pkgConf "github.com/donyhuang/go-server/pkg/conf"
	"github.com/donyhuang/go-server/pkg/db"
	"github.com/donyhuang/go-server/pkg/ginx/middleware"
	"github.com/donyhuang/go-server/pkg/ginx/route"
	"github.com/donyhuang/go-server/pkg/grpcx"
	"github.com/donyhuang/go-server/pkg/grpcx/interceptor"
	log2 "github.com/donyhuang/go-server/pkg/log"
	kafka "github.com/donyhuang/go-server/pkg/mq"
	"github.com/donyhuang/go-server/pkg/netx"
	"github.com/donyhuang/go-server/pkg/proc"
	"github.com/donyhuang/go-server/pkg/prometheus"
	"github.com/donyhuang/go-server/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"sync"
)

var (
	instance *Server
	once     sync.Once
)

type doneFunc func()
type Server struct {
	ctx               context.Context
	configFile        string
	globalConf        *AppConf
	httpRoutes        map[string]*route.Routes
	grpcHandle        map[string]*server.GrpcHandler
	kafkaConsumer     map[string]sarama.ConsumerGroupHandler
	doneFunc          []doneFunc
	kInstance         *kafka.Kafka
	httpServer        *server.HttpServer
	grpcServer        *server.GrpcServer
	kafkaConsumerName string
}

func NewServer(configFile string) *Server {
	once.Do(func() {
		instance = &Server{
			configFile: configFile,
			doneFunc:   make([]doneFunc, 0),
		}
		instance.init()
	})
	return instance
}
func (s *Server) AddDoneFunc(fn doneFunc) {
	proc.AddDoneFn(proc.DownFn(fn))
}
func (s *Server) init() {
	s.ctx = proc.GetProcSignalCtx()
	viper.SetConfigFile(s.configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config err %v\n", err)
	}
	var appConf AppConf
	err = viper.Unmarshal(&appConf)
	if err != nil {
		log.Fatalf("config unmarshal err %v\n", err)
	}
	s.globalConf = &appConf
	log2.InitLog(appConf.Log)
	if len(appConf.Storage.Mysql) != 0 {
		err = db.InitMysql(appConf.Storage.Mysql, string(appConf.Server.Mode))
		if err != nil {
			log.Fatalf("init mysql err %v\n", err)
		}
		proc.AddBseDoneFn(func() {
			db.CloseMysql()
			log.Println("mysql close")
		})
	}
	if len(appConf.Storage.Clickhouse) != 0 {
		err = db.InitClickhouse(appConf.Storage.Clickhouse)
		if err != nil {
			log.Fatalf("init clickhouse err %v\n", err)
		}
		proc.AddBseDoneFn(func() {
			db.CloseClick()
			log.Println("clickhouse close")
		})
	}
	if appConf.Kafka.BrokerList != "" {
		myKafka, err := kafka.NewKafka(&appConf.Kafka)

		if err != nil {
			log.Fatalf("init kafka err %v\n", err)
		}
		s.kInstance = myKafka
		proc.AddBseDoneFn(func() {
			myKafka.Close()
			log.Println("kafka Close ")
		})
	}
	var cancelListenFnList []pkgConf.CancelListenFunc
	if len(appConf.Conf.Nacos.Data) != 0 {
		var err error
		cancelListenFnList, err = pkgConf.InitNacos(&appConf.Conf.Nacos)
		if err != nil {
			log.Fatalf("init nacos err %v\n", err)
		}
		proc.AddBseDoneFn(func() {
			for _, v := range cancelListenFnList {
				vCopy := v
				_ = vCopy()
			}
			log.Println(" naCos stop listen")
		})
	}
	if len(appConf.Cache.Redis) > 0 {
		cache.NewPools(appConf.Cache.Redis)
		proc.AddBseDoneFn(cache.StopAllPools)
	}
	s.initServer()
}

// Start server and wait shutDown
func (s *Server) Start() {
	s.start()
	proc.ShutDown()
}
func (s *Server) GetGlobalConf() *AppConf {
	return s.globalConf
}
func (s *Server) RegisterGrpcHandle(name string, handler *server.GrpcHandler) {
	if s.grpcHandle == nil {
		s.grpcHandle = make(map[string]*server.GrpcHandler)
	}
	s.grpcHandle[name] = handler
}
func (s *Server) RegisterKafkaConsumer(name string, handler sarama.ConsumerGroupHandler) {
	if s.kafkaConsumer == nil {
		s.kafkaConsumer = make(map[string]sarama.ConsumerGroupHandler)
	}
	s.kafkaConsumer[name] = handler
}

func (s *Server) RegisterHttpRoutes(name string, routes *route.Routes) {
	if s.httpRoutes == nil {
		s.httpRoutes = make(map[string]*route.Routes)
	}
	s.httpRoutes[name] = routes
}

func (s *Server) WithGrpcOption(opt ...grpc.ServerOption) *Server {
	s.grpcServer.WithGrpcOption(opt...)
	return s
}
func (s *Server) WithHttpMiddleWare(h ...gin.HandlerFunc) {
	s.httpServer.Use(h...)
}
func (s *Server) AddRoute(r route.Route) {
	s.httpServer.AddRoute(r)
}
func (s *Server) initServer() {
	if len(s.globalConf.Server.Service) == 0 {
		return
	}

	//metrics const label
	label := map[string]string{
		"env": string(s.globalConf.Server.Mode),
		"ip":  netx.InternalIp(),
	}
	for _, si := range s.globalConf.Server.Service {
		cSi := si
		switch cSi.Protocol {
		case server.Http:
			s.httpServer = s.initHttp(cSi, s.globalConf.Server.Mode)
			if s.globalConf.Server.Prometheus.Port != 0 {
				s.httpServer.Use(middleware.PrometheusHandler(label))
			}
		case server.Grpc:
			s.grpcServer = s.initGrpc(cSi, s.globalConf.Server.Mode)
			if s.globalConf.Server.Prometheus.Port != 0 {
				s.grpcServer.WithGrpcUnaryServerInterceptor(interceptor.UnaryPrometheusInterceptor(label))
			}
		case server.Kafka:
			s.kafkaConsumerName = cSi.Name
		}

	}

}
func (s *Server) start() {
	if s.globalConf.Server.Prometheus.Port != 0 {
		prometheus.StartHandler(s.globalConf.Server.Prometheus)
	}
	if s.httpServer != nil {
		if s.httpRoutes[s.httpServer.Name] == nil {
			panic("need register http routes " + s.httpServer.Name)
		}
		s.httpServer.WithRoutes(s.httpRoutes[s.httpServer.Name])
		s.httpServer.BuildRoutes()
		go func() {
			s.httpServer.Start()
		}()
		proc.AddOutDoneFn(s.httpServer.ShutDown)
	}
	if s.grpcServer != nil {
		if s.grpcHandle[s.grpcServer.Name] == nil {
			panic("need register http routes " + s.grpcServer.Name)
		}

		go func() {
			s.grpcServer.Start(s.grpcHandle[s.grpcServer.Name])
		}()
		proc.AddOutDoneFn(s.grpcServer.ShutDown)
	}

	s.initKafkaConsumer()

}

func (s *Server) initHttp(sv server.Service, mode server.Mode) *server.HttpServer {
	svr := server.NewHttpServer(sv, mode)
	svr.Use(middleware.LoggerMiddleware, middleware.RecoveryWithWriter(log2.GetLogRecovery()))
	return svr
}

func (s *Server) initGrpc(sv server.Service, mode server.Mode) *server.GrpcServer {
	grpclog.SetLoggerV2(&grpcx.GrpcLoggerV2{Logger: log2.Logger})
	svr := server.NewGrpcServer(sv, mode, interceptor.LoggingInterceptor, interceptor.RecoveryInterceptor)
	svr.WithGrpcOption(grpc.NumStreamWorkers(sv.WorkerNum))
	return svr
}

func (s *Server) initKafkaConsumer() {
	if s.kafkaConsumerName == "" {
		return
	}
	handle := s.kafkaConsumer[s.kafkaConsumerName]
	if handle == nil {
		panic("need register kafka consumer " + s.kafkaConsumerName)
	}
	go func() {
		for {
			err := s.kInstance.Consumer(s.ctx, handle)
			if err != nil {
				log.Printf("consumer err %v", err)
				return
			}
			select {
			case <-s.ctx.Done():
				return
			default:
			}
		}
	}()
}

func (s *Server) GetContext() context.Context {
	return s.ctx
}
