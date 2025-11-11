# Lain CLI
## Let's all love lain

**一个基于 [Cobra](https://github.com/spf13/cobra) 构建的现代 TUI 命令行工具。**  
融合系统信息监控、网络测试、AI 模型调用、Markdown 输出等功能，旨在成为开发者的轻量全能 CLI 助手。

---

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-informational?style=flat-square)]()
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Demo-green.svg?style=flat-square)]()
[![UI](https://img.shields.io/badge/UI-TUI-orange?style=flat-square)]()

---

## 简介

Lain-CLI 是一个以终端交互为核心的工具集，整合多种系统能力与信息输出方式。  
设计目标是 **“在命令行中完成 80% 的日常系统与网络工作”**，  

---

## 功能模块

| 模块 | 功能说明 | 状态 |
|:------|:------|:------:|
| **基础输出美化** | 自定义 CLI 输出格式与颜色 | ✅ |
| **AI 模型调用** | 支持管道输入与上下文记忆 | ✅ |
| **网络工具集** | 网络连通性、速率测试、端口扫描 | ✅ |
| **系统监控** | CPU、内存、磁盘、网络实时刷新 | ✅ |
| **系统信息展示** | 类似 Neofetch 的静态信息输出 | ✅ |
| **Markdown 工具** | 输出转换与格式美化 | ✅ |
| **MCP 协议支持** | 初步兼容、支持配置调用 | ✅ |
|  | Linux端 MCP适配测试成功 | ✅ |
|  | mcp文件ui编辑 | ✅ |
|  | 更好的mcp调用方式 | **TODO** |
| **Linux 完全支持** |  | ✅ |
| **密钥安全** |  | ✅ |



---

## 使用示例

```bash
# 模型调用（管道方式）
echo "show network info" | lain-cli ag "帮我总结一下"

# 网络测试
- 网速测试
 - lain-cli test speed #-n 不用tui测
- 端口测试
 - lain-cli test port ip/domain  # -o 返回不会立刻refuse的端口
 - lain-cli test port ip/domain -p [ports]
 - lain-cli test port ip -s startport -e endport      # 扫描端口范围 

# 基本系统监控
lain-cli top

# 系统信息展示  [ui 使用Lain 中 NAVI电脑的图标]
lain-cli info

# 在终端渲染 md 文件
lain-cli md    # -w 在新窗口展示

# 终端调用简单mcp
 lain-cli mcps  # -f [filename] 调用mcp 返回输出到文件
 lain-cli mcps repl # 交互方式的使用mcp 

# 查看版本信息
lain-cli version 

```

### config.yaml 配置

```yaml
# AI 接口信息 [支持openwebui]
ai:
  api_url: ""
  api_key: ""
  ai_model_name: ""
# 是否使用一言 ON 
yiyan: 
  status: "ON"
  api_url: "https://v1.hitokoto.cn"

# 是否使用模型上下文 [基于shell 进程]
context:
  enabled: true
  local: "context.md" # /tmp/context.txt

# mcp 文件位置
mcp:
  json: "./mcp.json"

# info 信息的logo位置
logo:
  logo_txt: "./ascii-logo.txt"

```

### json配置示例
```json
{
  "mcpServers": {
    "rss-reader-mcp": {
      "command": "npx",
      "args": [
        "-y",
        "rss-reader-mcp"
      ]
    },
    "fetch": {
      "args": [
        "mcp-server-fetch"
      ],
      "command": "uvx"
    },
    "douyin-mcp": {
      "args": [
        "douyin-mcp-server"
      ],
      "command": "uvx",
      "env": {
        "DASHSCOPE_API_KEY": "ENV_DASHSCOPE_API_KEY"
      }
    }
  }
}

```
#### 密钥配置
如果有密钥不方便出现在`mcp.json`文件中 可以设置以`ENV_`开头的环境变量将会自动映射到程序中
> export ENV_DASHSCOPE_API_KEY="sk-xxxxxxxxxx"


### 安装

一键安装脚本：
```bash
curl -fsSL https://raw.githubusercontent.com/rescheni/lain-cli/refs/heads/main/scripts/install.sh | sh
```
可以通过源码编译方式使用：
``` bash
git clone https://github.com/rescheni/lain-cli.git
cd lain-cli
go build -o lain-cli
```


### 其他
> 实际上这个项目有很多优化空间，目前本人水平有限.
> 如有建议可以联系我： reschen@126.com

### 未来计划

开源协议
本项目采用 MIT License 开源协议

致开发者
这是一个从实验学习起步的项目，
我希望它能成为 “命令行世界里最有温度的工具”。
如果你有兴趣参与、测试或改进，欢迎提交 Issue 或 PR。

> Author: re
> Project: github.com/rescheni/lain-cli