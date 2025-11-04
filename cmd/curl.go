package cmd

import (
	"lain-cli/exec"
	"lain-cli/logs"

	"github.com/spf13/cobra"
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "HTTP Test Tools		# 基本的http 测试工具",
	Long: `
HTTP Test tools 
支持 POST,GET,PUT 等协议
lain curl get url	#  向目标服务器发送GET请求
lain curl post url  #  像目标服务器发送POST请求
					#  打开窗口输入请求json ctrl+s发送
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logs.Err("Place input [GET/POST...] URL")
			return
		}
		err := exec.Curl(args[0], args[1])
		if err != nil {
			logs.Err("", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(curlCmd)
}
