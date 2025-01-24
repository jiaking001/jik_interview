package middleware

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

// SnowflakeNode 全局雪花算法节点
var SnowflakeNode *snowflake.Node

func init() {
	// 初始化雪花算法节点
	var err error
	SnowflakeNode, err = snowflake.NewNode(1) // 1 是节点 ID，需确保唯一
	if err != nil {
		fmt.Println("Failed to create snowflake node:", err)
		return
	}
}
