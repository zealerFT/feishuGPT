package snowflakesvc

import (
	"feishu/util"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

type Snowflake struct {
	*snowflake.Node
}

func MestNewSnowflake() *Snowflake {
	return &Snowflake{NewSnowflake()}
}

func NewSnowflake() *snowflake.Node {
	node, err := snowflake.NewNode(int64(util.GetMachineNo())) // node 后期需要根据实际情况处理，为机器号
	if err != nil {
		log.Panic().Msgf("雪花算法初始化失败～ %v", err)
	}
	return node
}
