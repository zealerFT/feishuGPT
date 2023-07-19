package cmd

import (
	"context"
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"feishu/config"
	"feishu/dep"
	"feishu/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "start http server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetRole("http")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("feishugpt api server")
		bootHTTPServer(config.Options())
	},
}

func bootHTTPServer(cfg *config.AppConfig) {
	// 依赖注入
	dependency := dep.DIHttpDependency()

	log.Info().Msg("feishugpt: boot HTTP server")
	server := http.New(
		http.ExportLogOption(),
		http.WithDependency(dependency),
		http.SetRouteOption(),
	)

	// 创建服务器
	srv := &nethttp.Server{
		Addr:    cfg.HTTPServerAddr,
		Handler: server,
	}
	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("feishugpt Shutting down server...")

	// 创建一个context，并设置超时时间，保重不会长时间的等待下去
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 调用Shutdown()方法，优雅地关闭服务器，阻止新的请求进来，并可以关闭RegisterOnShutdown注册的方法
	// 注意：这里的关闭每一个协程，其实是每一个用户请求，也就是说对外接口请求过长的，这里可以等待请求完成（没到Timeout），但是单个接口请求立即返回，但这个请求里其实开启了一个协程来处理耗时任务，这种是无法shutdown的
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err)
	}

	time.Sleep(1 * time.Second) // 如果某些耗时的go routine,为了可以在net/http shutdown以后也关闭，就需要这样一个sleep来等待所以协程终止

	fmt.Println("Server exiting")
}
