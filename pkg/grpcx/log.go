package grpcx

import "github.com/sirupsen/logrus"

type GrpcLoggerV2 struct {
	*logrus.Logger
}

func (l *GrpcLoggerV2) V(level int) bool {
	return l.Logger.IsLevelEnabled(logrus.Level(level))
}
