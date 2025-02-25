package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"strconv"
	"strings"
)

// StringToStrings 字符串转字符串数组
func StringToStrings(s string) ([]string, error) {
	// 原始的 JSON 格式的字符串
	jsonStr := s

	// 定义一个字符串数组变量来存储解析后的结果
	var tagList []string

	// 使用 json.Unmarshal 解析 JSON 字符串
	err := json.Unmarshal([]byte(jsonStr), &tagList)
	if err != nil {
		return nil, err
	}

	return tagList, nil
}

// StringsToString 将字符串数组转化为字符串
func StringsToString(s []string) string {
	// 将字符串数组转化为字符串
	quotedTags := make([]string, len(s))
	// 向字符串数组添加双引号 使连接成的字符串每个值都带有双引号
	for i, tag := range s {
		quotedTags[i] = strconv.Quote(tag) // 使用 strconv.Quote 添加双引号
	}
	tags := strings.Join(quotedTags, ",")
	tags = "[" + tags + "]"

	return tags
}

// Int64TOString int64转化成字符串
func Int64TOString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 字符串转化成int64
func StringToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// StringToUint64 字符串转uint64
func StringToUint64(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// Uint64TOString uint64转化成字符串
func Uint64TOString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// GetIPAddress 获取ip地址
func GetIPAddress(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Request.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return strings.Split(ip, ",")[0]
}

// GetDeviceType 获取用户登录设备
func GetDeviceType(userAgent string) string {
	ua := useragent.New(userAgent)
	if ua.Mobile() {
		return "mobile"
	} else if ua.Bot() {
		return "bot"
	} else {
		return "pc"
	}
}
