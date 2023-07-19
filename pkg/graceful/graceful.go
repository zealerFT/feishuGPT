package graceful

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"feishu/pkg/healthcheck"

	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

type Do struct {
	once        sync.Once
	onceMetrics sync.Once
}

func StandBy(addr string, f func()) {
	stop := WaitSignals()

	grace := &Do{}
	done := grace.Do(addr, f)

	for {
		select {
		// 函数正常执行结束后，chan将返回空，这里是正常结束
		case <-done:
			return
		// 当遇到指定single信号，直接结束
		case <-stop:
			return
		}
	}
}

func (g *Do) Do(addr string, f func()) <-chan struct{} {
	stop := make(chan struct{})

	// start health at once，为pod提供的监控检查livenessProbe and readinessProbe
	go g.withHealthCheck(addr, stop)

	go func() {
		// 这里的defer是为了结束stop chan，这样外层的select将捕获到f()执行结束
		defer close(stop)

		func() {
			defer HandleCrash()
			f()

		}()
	}()

	return stop
}

func (g *Do) withHealthCheck(addr string, stop <-chan struct{}) {
	g.once.Do(func() {
		HTTPHealthCheck(addr, stop)
	})
}

func HTTPHealthCheck(addr string, stop <-chan struct{}) {
	server := &http.Server{Addr: addr, Handler: healthcheck.NewHandler()}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("[BgHealthCheck] health server close with err: %+v", err)
		}
	}()
	<-stop
	server.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("[BgHealthCheck] stop server graceful stop with err: %+v", err)
	}
}

func WaitSignals() chan struct{} {
	stop := make(chan struct{})

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		<-quit
		close(stop)
	}()

	return stop
}

func logPanic(r interface{}) {
	callers := getCallers(r)
	if _, ok := r.(string); ok {
		log.Printf("observed a panic: %s\n%v", r, callers)
	} else {
		log.Printf("observed a panic: %#v (%v)\n%v", r, r, callers)
	}
}

func getCallers(r interface{}) string {
	callers := ""
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = callers + fmt.Sprintf("%v:%v\n", file, line)
	}

	return callers
}

var PanicHandlers = []func(interface{}){logPanic}

func HandleCrash(additionalHandlers ...func(interface{})) {
	if r := recover(); r != nil {
		for _, fn := range PanicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		// Actually proceed to panic.
		panic(r)
	}
}
