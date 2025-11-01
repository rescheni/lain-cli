package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

func TestDisk(t *testing.T) {
	// --- 1. 打印所有分区 ---
	// 这会告诉我们 "/" 挂载在哪个设备上 (例如 /dev/disk1s1)
	fmt.Println("--- Partitions (disk.Partitions) ---")
	partitions, err := disk.Partitions(false) // false = 只看物理设备
	if err != nil {
		log.Fatalf("获取 Partitions 失败: %v", err)
	}

	var rootDevice string
	for _, p := range partitions {
		fmt.Printf("  Device: %s, Mountpoint: %s, Fstype: %s\n", p.Device, p.Mountpoint, p.Fstype)
		if p.Mountpoint == "/" {
			rootDevice = p.Device
		}
	}
	fmt.Printf("\n>>> 你的根 (\"/\") 挂载在: %s\n", rootDevice)

	// --- 2. 打印所有 IOCounters 的 "Keys" ---
	// 这会告诉我们 gopsutil *真正* 使用的名字 (例如 "disk1s1" 或 "disk0")
	fmt.Println("\n--- IO Counters (disk.IOCounters) ---")
	ioStats, err := disk.IOCounters() // 关键：不带参数，获取所有
	if err != nil {
		log.Fatalf("获取 IOCounters 失败: %v", err)
	}

	fmt.Println("可用的 IO Counter Keys:")
	for key, stat := range ioStats {
		fmt.Printf("  Key: \"%s\" (ReadBytes: %d, WriteBytes: %d)\n", key, stat.ReadBytes, stat.WriteBytes)
	}

	fmt.Println("\n--- 诊断 ---")
	fmt.Println("请比较上面的 'Partitions' 和 'IO Counters' 列表。")
	fmt.Println("你需要找到 'IO Counters' 列表中的哪个 Key 对应你的根设备。")
	fmt.Println("例如: 如果你的根设备是 '/dev/disk1s1'，正确的 Key 可能是 'disk1s1' 或 'disk1'。")
}

func TestDiskSpeed(t *testing.T) {
	ctx := context.Background()
	// 第一次快照
	counters1, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		panic(err)
	}
	t1 := time.Now()

	// 等待一段时间，比如 1 秒
	time.Sleep(2 * time.Second)

	// 第二次快照
	counters2, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		panic(err)
	}
	t2 := time.Now()

	elapsed := t2.Sub(t1).Seconds()

	for name, cnt2 := range counters2 {
		cnt1, ok := counters1[name]
		if !ok {
			continue
		}
		deltaReadBytes := float64(cnt2.ReadBytes - cnt1.ReadBytes)
		deltaWriteBytes := float64(cnt2.WriteBytes - cnt1.WriteBytes)
		readSpeed := deltaReadBytes / elapsed   // 字节/秒
		writeSpeed := deltaWriteBytes / elapsed // 字节/秒

		fmt.Printf("Device: %s\n", name)
		fmt.Printf("  Read speed: %.2f MB/s\n", readSpeed/(1024*1024))
		fmt.Printf("  Write speed: %.2f MB/s\n", writeSpeed/(1024*1024))
	}
}
