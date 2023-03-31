package server

type Protocol string

const (
	Http  Protocol = "http"
	Grpc  Protocol = "grpc"
	Kafka Protocol = "kafka"
)

type Mode string

const (
	DEBUG   Mode = "debug"
	RELEASE Mode = "release"
)

type Service struct {
	Name      string
	Port      int
	Network   string
	Protocol  Protocol
	WorkerNum uint32
	Http      HttpService
}
type HttpService struct {
	Health string
}
