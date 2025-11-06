package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	mui "github.com/rescheni/lain-cli/internal/ui"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rescheni/lain-cli/config"
	"github.com/rescheni/lain-cli/logs"
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

func Init() error {
	err := initMCPs(config.Conf.Mcp.Json)
	if err != nil {
		logs.Err("open mcp.json", err)
		logs.Err("MCP Location " + config.Conf.Mcp.Json + " open error")
	}
	return err
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
		logs.Info(fmt.Sprintf("初始化 MCP 客户端: %s (cmd=%s args=%v)\n", name, srv.Command, srv.Args))
		cmd := exec.Command(srv.Command, srv.Args...)
		cmd.Env = os.Environ()
		for k, v := range srv.Env {
			config.Check_ENV(&v)
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}

		cmd.Stderr = os.Stderr
		client := mcp.NewClient(&mcp.Implementation{
			Name:    "lain-cli",
			Version: "v1.0.0",
		}, nil)

		session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
		if err != nil {
			logs.Err("❌ 初始化 "+name+" 失败: ", err)
			continue
		}
		Mcps[name] = session
		// time.Sleep(500 * time.Millisecond)
	}

	logs.Info("✅ 所有 MCP 初始化完成")
	return nil
}

// 列出所有 MCP 名称
func ListMCPs() []string {
	mcps := []string{}
	for name := range Mcps {
		mcps = append(mcps, name)
	}
	return mcps
}

// 调用某个 MCP 的工具列表
func ListMCPTools(ctx context.Context, name string) {
	sess, ok := Mcps[name]
	if !ok {
		logs.Err("❌ 未找到 MCP:" + name)
		return
	}

	resp, err := sess.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		logs.Err("❌ ListTools 失败:", err)
		return
	}

	if len(resp.Tools) == 0 {
		logs.Err("(没有可用工具)")
		return
	}

	fmt.Printf("🧰 %s 可用工具:\n", name)
	for i, tool := range resp.Tools {
		fmt.Printf("\t%d - %s: %s\n", i+1, tool.Name, tool.Description)
	}
}

// 调用工具
func CallTool(ctx context.Context, name, tool string, args map[string]any, tofile string) {
	sess, ok := Mcps[name]
	if !ok {
		logs.Err("❌ 未找到 MCP:" + name)
		return
	}
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
		logs.Err("调用工具失败:", err)
		return
	}

	if res.IsError {
		logs.Err("⚠️ 工具执行失败")
		return
	}

	mds := ""
	for _, c := range res.Content {
		if text, ok := c.(*mcp.TextContent); ok {
			mds += text.Text
		}
	}
	if tofile != "" {
		file, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			logs.Err("Open File error")
		}
		defer file.Close()
		file.Write([]byte(mds))
	}
	mui.PrintMarkdown(mds, false)
}

// 关闭所有 MCP 会话
func CloseAllMCPs() {
	for name, s := range Mcps {
		_ = s.Close()
		logs.Info("已关闭 MCP:" + name)
	}
}
