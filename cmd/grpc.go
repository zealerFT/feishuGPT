package cmd

import (
	"feishu/config"
	"feishu/dep"
	"feishu/pkg/graceful"

	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start grpc server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SetRole("grpc")
	},
	Run: func(cmd *cobra.Command, args []string) {
		graceful.StandBy(config.Options().HTTPServerAddr, func() {
			server := dep.DIGRPCServer()
			defer func() {
				_ = server.Close()
			}()
			server.Run()
		})
	},
}
