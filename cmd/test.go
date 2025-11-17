/*
Copyright © 2025 re
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/rescheni/lain-cli/internal/exec"
	"github.com/rescheni/lain-cli/internal/tools"
	logs "github.com/rescheni/lain-cli/logs"

	"github.com/spf13/cobra"
)

var (
	nuiFlag  bool
	openFlag bool
	start    int
	end      int
	ip       string
	ports    []int
)

// ---------- testCmd ----------
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Basic network testing tools		#基本的网络测试工具",
	Long:  `提供网络测试相关功能，包括端口扫描和网速测试。`,
}

// ---------- test port ----------
var testPortCmd = &cobra.Command{
	Use:   "port <ip> [flags]",
	Short: "Test network ports  测试端口连通性",
	Example: `
  lain-cli test port 192.168.1.1                	# 默认扫描常用端口
  lain-cli test port 192.168.1.1 -s 20 -e 8080    	# 扫描端口范围 20-8080
  lain-cli test port 192.168.1.1 -p 22 80 443   	# 测试指定端口
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ip = args[0] // 第一个参数是 IP 或域名

		if openFlag {
			tools.SetScannerOpen()
		}

		logs.Info("目标地址:" + ip)

		// ----- 模式一：端口范围 -----
		if start > 0 && end > 0 {
			logs.Info(fmt.Sprintf("扫描端口范围: %d - %d\n", start, end))
			exec.Rfun = func() {
				exec.RunNmap(ip, start, end)
			}
			exec.Run()
			return nil
		}
		if len(args) > 1 {
			for _, a := range args[1:] {
				if v, err := strconv.Atoi(a); err == nil {
					ports = append(ports, v)
				}
			}
		}
		// ----- 模式二：指定端口 -----
		if len(ports) > 0 {
			logs.Err(fmt.Sprintf("测试指定端口: %v\n", ports))
			exec.Rfun = func() {
				exec.RunNmapPorts(ip, ports...)
			}
			exec.Run()
			return nil
		}

		// ----- 模式三：默认常用端口 -----
		logs.Info("未指定端口，将测试默认常用端口")
		exec.Rfun = func() {
			exec.RunDefaultNmap(ip)
		}
		exec.Run()
		return nil
	},
}

// ---------- test speed ----------
var testSpeedCmd = &cobra.Command{
	Use:   "speed",
	Short: "Test network speed  测试网速",
	Example: `
  lain test speed       # 进行普通测速
  lain test speed -u    # 使用 UI 模式测速
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		exec.RunSpeedTestUI(nuiFlag)
		return nil
	},
}

func init() {
	// 子命令注册
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testPortCmd)
	testCmd.AddCommand(testSpeedCmd)

	// test port flags
	testPortCmd.Flags().IntVarP(&start, "start", "s", 0, "Port range start")
	testPortCmd.Flags().IntVarP(&end, "end", "e", 0, "Port range end")
	testPortCmd.Flags().IntSliceVarP(&ports, "port", "p", nil, "Specify ports to test")
	testPortCmd.Flags().BoolVarP(&openFlag, "isopen", "o", false, "View all open ports")

	// test speed flags
	testSpeedCmd.Flags().BoolVarP(&nuiFlag, "nui", "n", false, "Use UI mode for speed test")
}
