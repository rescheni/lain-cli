package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRodmoInt() string {
	rand.Seed(time.Now().UnixNano())
	// 随机取 0~255 之间的颜色编号
	randomColor := fmt.Sprintf("%d", rand.Intn(256))
	return randomColor
}
