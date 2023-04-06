package proc

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

type DownFn func()

var (
	//用户自定义关闭任务
	doneFn []DownFn
	//对外任务需要优先关闭
	outDoneFn []DownFn
	//框架基础任务需要最后关闭 ，mysql,clickhouse ,redis,remoteConf(naCos)
	baseFn []DownFn

	mut       sync.Mutex
	outMut    sync.Mutex
	baseMut   sync.Mutex
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
func AddOutDoneFn(fn DownFn) {
	outMut.Lock()
	defer outMut.Unlock()
	outDoneFn = append(outDoneFn, fn)
}
func AddBseDoneFn(fn DownFn) {
	baseMut.Lock()
	defer baseMut.Unlock()
	baseFn = append(baseFn, fn)
}
func doAllDoneFn(doneFns []DownFn) {
	if len(doneFns) == 0 {
		return
	}
	wg.Add(len(doneFns))
	for _, fn := range doneFns {
		vfn := fn
		go func() {
			vfn()
			wg.Done()
		}()
	}
	wg.Wait()
}
func ShutDown() {
	<-signalCtx.Done()
	stop()
	doAllDoneFn(outDoneFn)
	doAllDoneFn(doneFn)
	doAllDoneFn(baseFn)
}
