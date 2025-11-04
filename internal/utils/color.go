package utils

import (
	"fmt"
	"math/rand"
)

func GetRodmoInt() string {
	// 随机取 0~255 之间的颜色编号
	randomColor := fmt.Sprintf("%d", rand.Intn(256))
	return randomColor
}
