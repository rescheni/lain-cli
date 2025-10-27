package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"lain-cli/config"
	mui "lain-cli/ui"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPConfig 顶层结构
type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

// MCPServer 单个服务配置
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

func Init() {
	initMCPs(config.Conf.Mcp.Json)
}

// 保存所有 MCP 连接
var Mcps = make(map[string]*mcp.ClientSession)

// 去除 JSON 注释（//）
func stripLineComments(b []byte) []byte {
	re := regexp.MustCompile(`(?m)^\s*//.*$`)
	return re.ReplaceAll(b, []byte(""))
}

// 读取配置文件
func loadMCPConfig(path string) (*MCPConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	data = stripLineComments(data)
	var cfg MCPConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return &cfg, nil
}

func initMCPs(configPath string) error {
	cfg, err := loadMCPConfig(configPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	for name, srv := range cfg.MCPServers {
		fmt.Printf("初始化 MCP 客户端: %s (cmd=%s args=%v)\n", name, srv.Command, srv.Args)

		cmd := exec.Command(srv.Command, srv.Args...)
		cmd.Env = os.Environ()
		for k, v := range srv.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}

		// 把子进程的 stderr 显示到当前终端，便于调试工具的日志输出（不会影响 MCP 的 stdout/stdin 协议）
		// 许多工具会把 human-readable 日志写到 stderr，这样不会破坏协议。
		cmd.Stderr = os.Stderr

		client := mcp.NewClient(&mcp.Implementation{
			Name:    "lain-cli",
			Version: "v1.0.0",
		}, nil)

		session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			fmt.Printf("❌ 初始化 %s 失败: %v\n", name, err)
			continue
		}

		Mcps[name] = session
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("✅ 所有 MCP 初始化完成")
	return nil
}

// 列出所有 MCP 名称
func ListMCPs() []string {
	mcps := []string{}
	fmt.Println("当前连接的 MCP:")
	for name := range Mcps {
		mcps = append(mcps, name)
		fmt.Printf(" - %s\n", name)
	}
	return mcps
}

// 调用某个 MCP 的工具列表
func ListMCPTools(ctx context.Context, name string) {
	sess, ok := Mcps[name]
	if !ok {
		fmt.Printf("❌ 未找到 MCP: %s\n", name)
		return
	}

	resp, err := sess.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		fmt.Printf("❌ ListTools 失败: %v\n", err)
		return
	}

	if len(resp.Tools) == 0 {
		fmt.Println("(没有可用工具)")
		return
	}

	fmt.Printf("🧰 %s 可用工具:\n", name)
	for _, tool := range resp.Tools {
		fmt.Printf(" - %s: %s\n", tool.Name, tool.Description)
	}
}

// 调用工具
func CallTool(ctx context.Context, name, tool string, args map[string]any) {
	sess, ok := Mcps[name]
	if !ok {
		fmt.Printf("❌ 未找到 MCP: %s\n", name)
		return
	}
	// 调试：打印请求参数
	if args == nil {
		fmt.Println("Call payload: <nil>")
	} else {
		if bb, err := json.Marshal(args); err == nil {
			fmt.Println("Call payload:", string(bb))
		} else {
			fmt.Println("Call payload: <marshal error>", err)
		}
	}

	res, err := sess.CallTool(ctx, &mcp.CallToolParams{
		Name:      tool,
		Arguments: args,
	})
	if err != nil {
		fmt.Printf("调用工具失败: %v\n", err)
		return
	}

	if res.IsError {
		fmt.Println("⚠️ 工具执行失败")
		return
	}

	mds := ""
	for _, c := range res.Content {
		if text, ok := c.(*mcp.TextContent); ok {
			mds += text.Text
		}
	}
	mui.PrintMarkdown(mds, false)
}

// 关闭所有 MCP 会话
func CloseAllMCPs() {
	for name, s := range Mcps {
		_ = s.Close()
		fmt.Printf("已关闭 MCP: %s\n", name)
	}
}
