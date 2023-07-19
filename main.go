package main

import (
	"flag"
	"fmt"
	"os"

	"feishu/cmd"

	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()
	log.Info().Msg("feishuGPT: is starting")
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Info().Msg("feishuGPT: is end")

}
