package constant

import "fmt"

const UserSignInRedisKeyPrefix = "user:signins"

// GetUserSignInRedisKey 用户签到Key
func GetUserSignInRedisKey(year string, userId string) string {
	return fmt.Sprintf("%s:%s:%s", UserSignInRedisKeyPrefix, year, userId)
}

// NextSetBit 获取下一次签到是哪一天
func NextSetBit(bitmapData []byte, start int) int {
	for i := start / 8; i < len(bitmapData); i++ {
		for j := 7 - start%8; j >= 0; j-- {
			if (bitmapData[i]>>j)&1 == 1 {
				return i*8 + (7 - j)
			}
		}
		start = (i + 1) * 8
	}
	return -1
}
