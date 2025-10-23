package test

import (
	"fmt"
	"lain-cli/tools"
	mui "lain-cli/ui"
	"sync/atomic"
	"testing"
)

var completedTasks atomic.Uint64 // 1. 定义一个 64 位的原子计数器
var totalTasks uint64

func TestNmap(t *testing.T) {

	pool := tools.NewDefaultPools()
	pool.Run()

	begin := 20
	end := 65535
	testNums := begin - end

	totalTasks = uint64(testNums)

	tools.SetScannerOpen()
	for i := 20; i <= 10000; i++ {
		pool.Add(func() {
			defer completedTasks.Add(1)
			tools.Scan("", i)
		})
	}

	pool.Stop()
	info := make([][]string, 0)
	i := 0
	for k, v := range tools.Set {
		i++
		t := []string{fmt.Sprint(i), k.IP, fmt.Sprint(k.PORT), v}
		info = append(info, t)
	}
	mui.TuiPrintTable([]string{"NUM", "IP/Domain Name", "PORT", "PROC"}, info)

}
