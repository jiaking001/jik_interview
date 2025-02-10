package constant

import "fmt"

const UserSignInRedisKeyPrefix = "user:signins"

// GetUserSignInRedisKey 用户签到Key
func GetUserSignInRedisKey(year string, userId string) string {
	return fmt.Sprintf("%s:%s:%s", UserSignInRedisKeyPrefix, year, userId)
}

// NextSetBit 获取下一次签到是哪一天
func NextSetBit(bitmapData []int64, start int) int {
	for i := start / 64; i < len(bitmapData); i++ {
		for j := start % 64; j < 64; j++ {
			if (bitmapData[i]>>j)&1 == 1 {
				return i*64 + (63 - j)
			}
		}
		start = (i + 1) * 64
	}
	return -1
}
