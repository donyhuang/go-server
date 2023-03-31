package prometheus

import (
	"fmt"
	"github.com/donyhuang/go-server/pkg/reflectx"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	once          sync.Once
	defaultConfig = Config{
		Path: "/metrics",
		Port: 8084,
	}
)

func StartHandler(c Config) {
	once.Do(func() {
		dc := defaultConfig
		_ = reflectx.SetStructNotEmpty(&dc, c)
		http.Handle(c.Path, promhttp.Handler())
		addr := fmt.Sprintf(":%d", c.Port)
		go func() {
			if err := http.ListenAndServe(addr, nil); err != nil {
				panic(err)
			}
		}()
	})

}
