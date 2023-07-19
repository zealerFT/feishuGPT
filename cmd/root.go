package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "cli for feishu",
}

func init() {

	// HTTP Service
	RootCmd.AddCommand(httpCmd)

	// GRPC service
	RootCmd.AddCommand(grpcCmd)

}
