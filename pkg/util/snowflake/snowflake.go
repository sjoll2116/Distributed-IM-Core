package snowflake

import (
	"kama_chat_server/pkg/zlog"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init 初始化雪花算法节点
// 在分布式集群中，每个机器应当分配不同的 Node ID (0-1023)
func Init(nodeID int64) {
	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		zlog.Fatal("Failed to initialize snowflake node: " + err.Error())
	}
}

// GenID 生成一个全局唯一且单调递增的 ID，返回字符串形式 (用于替代原来简单随机字符串)
func GenID() string {
	if node == nil {
		// 默认兜底为 Node 1
		Init(1)
	}
	return node.Generate().String()
}
