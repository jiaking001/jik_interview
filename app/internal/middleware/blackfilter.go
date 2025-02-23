package middleware

import (
	"app/pkg/utils"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"net/http"
	"strings"
	"sync"
)

// 初始化布隆过滤器
var (
	configClient config_client.IConfigClient
	once         sync.Once
	mutex        sync.Mutex
	filter       = bloom.NewWithEstimates(1000, 0.01)
)

func initConfigClient() {
	// Nacos 服务器配置
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	// 客户端配置
	cc := constant.ClientConfig{
		NamespaceId:         "public", // 命名空间ID
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}

	// 创建配置客户端
	var err error
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create config client: %v", err)
	}
}

func loadBlacklist() {
	dataId := "jik_blacklist"
	group := "DEFAULT_GROUP"

	// 获取黑名单配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Fatalf("Failed to get blacklist config: %v", err)
	}

	// 将黑名单添加到布隆过滤器中
	blacklist := strings.Split(content, ",")
	mutex.Lock()
	defer mutex.Unlock()
	filter.ClearAll()
	for _, item := range blacklist {
		filter.AddString(item)
	}
}

func BlackFilter() {
	// 确保configClient只被初始化一次
	once.Do(initConfigClient)

	loadBlacklist()

	// 监听配置变化
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: "jik_blacklist",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			// 更新黑名单
			newBlacklist := strings.Split(data, ",")
			mutex.Lock()
			defer mutex.Unlock()
			filter.ClearAll()
			for _, item := range newBlacklist {
				filter.AddString(item)
			}
		},
	})
	if err != nil {
		log.Fatalf("Failed to listen config change: %v", err)
	}
}

func isInBlacklist(item string, filter *bloom.BloomFilter) bool {
	return filter.TestString(item)
}

func BlacklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := utils.GetIPAddress(c)
		if isInBlacklist(ip, filter) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Blacklisted IP"})
			c.Abort()
			return
		}
		c.Next()
	}
}
