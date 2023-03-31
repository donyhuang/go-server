package proc

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

type DownFn func()

var (
	doneFn    []DownFn
	mut       sync.Mutex
	wg        sync.WaitGroup
	signalCtx context.Context
	stop      context.CancelFunc
)

func init() {
	signalCtx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}
func GetProcSignalCtx() context.Context {
	return signalCtx
}
func AddDoneFn(fn DownFn) {
	mut.Lock()
	defer mut.Unlock()
	doneFn = append(doneFn, fn)
}

func ShutDown() {
	<-signalCtx.Done()
	stop()
	mut.Lock()
	defer mut.Unlock()
	wg.Add(len(doneFn))
	for _, fn := range doneFn {
		vfn := fn
		go func() {
			vfn()
			wg.Done()
		}()
	}
	wg.Wait()
}
