package server

import (
	"context"
	"fmt"
	"github.com/donyhuang/go-server/pkg/ginx/route"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

type HttpServer struct {
	Name   string
	e      *gin.Engine
	svr    *http.Server
	routes *route.Routes
}

func NewHttpServer(sv Service, mode Mode) *HttpServer {
	gin.SetMode(string(mode))
	e := gin.New()
	if sv.Http.Health != "" {
		e.Handle(http.MethodGet, sv.Http.Health, func(c *gin.Context) {
		})
	}
	return &HttpServer{
		e: e,
		svr: &http.Server{
			Addr:    fmt.Sprintf(":%v", sv.Port),
			Handler: e,
		},
		Name: sv.Name,
	}
}
func (s *HttpServer) WithRoutes(r *route.Routes) {
	s.routes = r
}
func (s *HttpServer) Use(h ...gin.HandlerFunc) {
	s.e.Use(h...)
}

func (s *HttpServer) BuildRoutes() {
	s.routes.BuildRoute(s.e)
}
func (s *HttpServer) AddRoute(r ...route.Route) {
	s.routes.Add(r...)
}

func (s *HttpServer) Start() {
	if err := s.svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
func (s *HttpServer) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.svr.Shutdown(ctx); err != nil {
		log.Fatal("HttpServer forced to stop: ", err)
	}
	log.Printf(" http server %v stop\n", s.Name)
}

type GrpcServer struct {
	Mode                   Mode
	Name                   string
	Port                   int
	NetWork                string
	unaryServerInterceptor []grpc.UnaryServerInterceptor
	grpcOpt                []grpc.ServerOption
	svr                    *grpc.Server
	handler                *GrpcHandler
}

func NewGrpcServer(sv Service, mode Mode, inter ...grpc.UnaryServerInterceptor) *GrpcServer {
	svr := &GrpcServer{
		Port:                   sv.Port,
		NetWork:                sv.Network,
		Name:                   sv.Name,
		Mode:                   mode,
		unaryServerInterceptor: inter,
	}
	return svr
}
func (g *GrpcServer) WithGrpcUnaryServerInterceptor(opt ...grpc.UnaryServerInterceptor) {
	g.unaryServerInterceptor = append(g.unaryServerInterceptor, opt...)
}
func (g *GrpcServer) WithGrpcOption(opt ...grpc.ServerOption) {
	g.grpcOpt = append(g.grpcOpt, opt...)
}

func (g *GrpcServer) Start(handle *GrpcHandler) {
	lis, err := net.Listen(g.NetWork, fmt.Sprintf(":%d", g.Port))
	if err != nil {
		log.Fatalf("init grpc %v err %v", g.Name, err)
	}
	log.Printf("grpc server listening at %v\n", lis.Addr())
	g.grpcOpt = append([]grpc.ServerOption{grpc.ChainUnaryInterceptor(g.unaryServerInterceptor...)}, g.grpcOpt...)
	g.svr = grpc.NewServer(g.grpcOpt...)
	g.svr.RegisterService(&handle.serviceDesc, handle.handle)
	if g.Mode != RELEASE {
		reflection.Register(g.svr)
	}
	if err := g.svr.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (g *GrpcServer) ShutDown() {
	g.svr.GracefulStop()
	log.Printf("grpc %v stop \n", g.Name)

}
