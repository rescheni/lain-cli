#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印信息的函数
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否以 root 权限运行
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_warn "不建议以 root 用户运行此安装脚本"
    fi
}

# 检测操作系统和架构
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"
    
    case $OS in
        Linux*)
            OS='linux'
            ;;
        Darwin*)
            OS='darwin'
            ;;
        *)
            print_error "不支持的操作系统: $OS"
            exit 1
            ;;
    esac
    
    case $ARCH in
        x86_64)
            ARCH='amd64'
            ;;
        aarch64|arm64)
            ARCH='arm64'
            ;;
        *)
            print_error "不支持的架构: $ARCH"
            exit 1
            ;;
    esac
    
    print_info "检测到平台: $OS/$ARCH"
}

# 创建临时目录
create_temp_dir() {
    TEMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TEMP_DIR"' EXIT
    print_info "创建临时目录: $TEMP_DIR"
}

# 下载最新版本
download_latest() {
    print_info "正在获取最新版本信息..."
    
    # 这里应该从 GitHub releases 下载最新版本
    # 目前使用固定版本进行演示
    VERSION="v0.1.0"
    BINARY_NAME="lain-cli"
    DOWNLOAD_URL="https://github.com/rescheni/lain-cli/releases/download/${VERSION}/lain-cli_${VERSION}_${OS}_${ARCH}.tar.gz"
    
    print_info "下载地址: $DOWNLOAD_URL"
    
    # 下载文件
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${TEMP_DIR}/lain-cli.tar.gz" "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${TEMP_DIR}/lain-cli.tar.gz" "$DOWNLOAD_URL"
    else
        print_error "需要 curl 或 wget 来下载文件"
        exit 1
    fi
    
    # 解压文件
    tar -xzf "${TEMP_DIR}/lain-cli.tar.gz" -C "$TEMP_DIR"
}

# 从源码构建（如果没有预编译版本）
build_from_source() {
    print_info "从源码构建..."
    
    if ! command -v go >/dev/null 2>&1; then
        print_error "未找到 Go 编译器，请先安装 Go 1.25+"
        exit 1
    fi
    
    # 克隆仓库或复制当前目录内容
    cp -r ./* "$TEMP_DIR/"
    cd "$TEMP_DIR"
    go build -o lain-cli .
}

# 安装文件
install_files() {
    INSTALL_DIR="${HOME}/.local/bin"
    CONFIG_DIR="${HOME}/.config/lain-cli"
    
    # 创建安装目录
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    
    # 复制二进制文件
    if [ -f "${TEMP_DIR}/lain-cli" ]; then
        cp "${TEMP_DIR}/lain-cli" "$INSTALL_DIR/"
        chmod +x "${INSTALL_DIR}/lain-cli"
        print_info "已安装 lain-cli 到 $INSTALL_DIR/lain-cli"
    else
        print_error "未找到编译后的二进制文件"
        exit 1
    fi
    
    # 复制配置文件示例（如果不存在）
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        if [ -f "${TEMP_DIR}/config.yaml.example" ]; then
            cp "${TEMP_DIR}/config.yaml.example" "${CONFIG_DIR}/config.yaml"
            print_info "已创建配置文件: ${CONFIG_DIR}/config.yaml"
        else
            print_warn "未找到配置文件示例"
        fi
    else
        print_info "配置文件已存在: ${CONFIG_DIR}/config.yaml"
    fi
}

# 显示安装后说明
post_install_info() {
    print_info "安装完成！"
    echo ""
    print_info "请确保将 $INSTALL_DIR 添加到 PATH 环境变量中:"
    echo "  export PATH=\$PATH:$INSTALL_DIR"
    echo ""
    print_info "配置文件位置: ${CONFIG_DIR}/config.yaml"
    print_info "你可以根据需要修改配置文件"
    echo ""
    print_info "运行 'lain-cli --help' 查看可用命令"
}

# 主函数
main() {
    print_info "开始安装 lain-cli"
    
    check_root
    detect_platform
    create_temp_dir
    
    # 尝试下载预编译版本，如果失败则从源码构建
    if ! download_latest; then
        print_warn "无法下载预编译版本，尝试从源码构建"
        build_from_source
    fi
    
    install_files
    post_install_info
}

# 如果直接运行此脚本，则执行 main 函数
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi