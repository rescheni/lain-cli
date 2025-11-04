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
|  | 输出结构化、可读性强 | ✅ |
| **网络工具集** | 网络连通性、速率测试、端口扫描 | ✅ |
|  | HTTP 请求测试与结果展示 | ✅ |
| **系统监控** | CPU、内存、磁盘、网络实时刷新 | ✅ |
| **系统信息展示** | 类似 Neofetch 的静态信息输出 | ✅ |
| **Markdown 工具** | 输出转换与格式美化 | ✅ |
| **MCP 协议支持** | 初步兼容、支持配置调用 | ✅ |
| **Linux 完全支持** | 适配与优化中 | ✅ |
|  | MCP适配测试失败 | todo |


---

## 使用示例

```bash
# 模型调用（管道方式）
echo "show network info" | lain-cli ag "帮我总结一下"

# 网络测试
lain-cli net test --host example.com

# 系统监控
lain-cli monitor

# 系统信息展示
lain-cli info
```
### 安装

未来将支持一键安装脚本：
```bash

curl -fsSL https://lain.sh/install.sh | sh

```
当前可通过源码编译方式使用：
``` bash
git clone https://github.com/yourname/lain-cli.git
cd lain-cli
go build -o lain-cli
```
功能预览
系统监控界面

Neofetch 风格信息展示

项目结构
```bash
github.com/rescheni/lain-cli/
├── cmd/            # 各子命令定义（Cobra）
├── tools/          # 工具模块封装
├── tui/            # TUI 视图组件
├── config/         # 全局配置
├── main.go
└── README.md
```
### 未来计划
- 完整 Linux 适配与系统调用抽象
- 更新与一键安装脚本
- 丰富 TUI 交互（进程管理 / 网络连接追踪）
- MCP 插件使用
- 云端同步与远程配置

开源协议
本项目采用 MIT License 开源协议

致开发者
这是一个从实验起步的项目，
我希望它能成为 “命令行世界里最有温度的工具”。
如果你有兴趣参与、测试或改进，欢迎提交 Issue 或 PR。

> Author: re
> Project: github.com/rescheni/lain-cli