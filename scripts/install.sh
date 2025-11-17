#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印信息的函数
print_info() {
    printf "%b[INFO]%b %s\n" "$GREEN" "$NC" "$1"
}

print_warn() {
    printf "%b[WARN]%b %s\n" "$YELLOW" "$NC" "$1"
}

print_error() {
    printf "%b[ERROR]%b %s\n" "$RED" "$NC" "$1"
}

# 检查是否以 root 权限运行
check_root() {
    if [ "$EUID" = "0" ]; then
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

# 获取最新版本
get_latest_version() {
    print_info "正在获取最新版本信息..."
    
    # 使用 GitHub API 获取最新版本
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s https://api.github.com/repos/rescheni/lain-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -q -O - https://api.github.com/repos/rescheni/lain-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_warn "未找到 curl 或 wget，使用默认版本"
        VERSION="v1.0.0"
    fi
    
    if [ -z "$VERSION" ]; then
        print_warn "无法获取最新版本，使用默认版本"
        VERSION="v1.0.0"
    fi
    
    print_info "使用版本: $VERSION"
}

# 下载预编译版本
download_prebuilt() {
    BINARY_NAME="lain-cli"
    DOWNLOAD_URL="https://github.com/rescheni/lain-cli/releases/download/${VERSION}/lain-cli-${VERSION}-${OS}-${ARCH}.tar.gz"
    
    print_info "下载地址: $DOWNLOAD_URL"
    
    # 下载文件
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "${TEMP_DIR}/lain-cli.tar.gz" "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "${TEMP_DIR}/lain-cli.tar.gz" "$DOWNLOAD_URL"
    else
        print_error "需要 curl 或 wget 来下载文件"
        return 1
    fi
    
    # 解压文件
    tar -xzf "${TEMP_DIR}/lain-cli.tar.gz" -C "$TEMP_DIR"
    return 0
}

# 安装文件
install_files() {
    INSTALL_DIR="${HOME}/.lain-cli"
    
    # 创建安装目录
    mkdir -p "$INSTALL_DIR"
    
    # 检查是否为更新操作
    if [ -f "${INSTALL_DIR}/lain-cli" ]; then
        print_info "检测到已有安装，执行更新操作"
        
        # 只复制二进制文件
        if [ -f "${TEMP_DIR}/lain-cli" ]; then
            cp "${TEMP_DIR}/lain-cli" "$INSTALL_DIR/"
            chmod +x "${INSTALL_DIR}/lain-cli"
            print_info "已更新 lain-cli 二进制文件到 $INSTALL_DIR/lain-cli"
        else
            print_error "未找到编译后的二进制文件"
            exit 1
        fi
        
        # 更新 logo 文件（如果存在）
        if [ -f "${TEMP_DIR}/ascii-logo.txt" ] && [ ! -f "${INSTALL_DIR}/ascii-logo.txt" ]; then
            cp "${TEMP_DIR}/ascii-logo.txt" "$INSTALL_DIR/"
            print_info "已安装logo文件到 $INSTALL_DIR/ascii-logo.txt"
        fi
    else
        # 全新安装
        print_info "执行全新安装"
        
        # 复制二进制文件
        if [ -f "${TEMP_DIR}/lain-cli" ]; then
            cp "${TEMP_DIR}/lain-cli" "$INSTALL_DIR/"
            chmod +x "${INSTALL_DIR}/lain-cli"
            print_info "已安装 lain-cli 到 $INSTALL_DIR/lain-cli"
        else
            print_error "未找到编译后的二进制文件"
            exit 1
        fi
        
        # 复制配置文件
        if [ -f "${TEMP_DIR}/config.yaml" ]; then
            cp "${TEMP_DIR}/config.yaml" "$INSTALL_DIR/"
            print_info "已安装配置文件到 $INSTALL_DIR/config.yaml"
        else
            print_warn "未找到配置文件"
        fi
        
        # 复制logo文件
        if [ -f "${TEMP_DIR}/ascii-logo.txt" ]; then
            cp "${TEMP_DIR}/ascii-logo.txt" "$INSTALL_DIR/"
            print_info "已安装logo文件到 $INSTALL_DIR/ascii-logo.txt"
        else
            print_warn "未找到logo文件"
        fi
    fi
}

# 将路径添加到 shell 配置文件中以实现永久生效
add_to_shell_config() {
    INSTALL_DIR="${HOME}/.lain-cli"
    EXPORT_LINE="export PATH=\$PATH:${INSTALL_DIR}"
    
    # 检测用户的 shell 类型
    USER_SHELL=$(basename "$SHELL")
    
    # 确定配置文件路径
    case "$USER_SHELL" in
        bash)
            SHELL_CONFIG="${HOME}/.bashrc"
            ;;
        zsh)
            SHELL_CONFIG="${HOME}/.zshrc"
            ;;
        *)
            SHELL_CONFIG="${HOME}/.profile"
            ;;
    esac
    
    # 检查配置文件是否存在，不存在则创建
    if [ ! -f "$SHELL_CONFIG" ]; then
        touch "$SHELL_CONFIG"
    fi
    
    # 检查是否已经添加过该路径
    if ! grep -qF "$EXPORT_LINE" "$SHELL_CONFIG"; then
        echo "" >> "$SHELL_CONFIG"
        echo "# Added by lain-cli installer" >> "$SHELL_CONFIG"
        echo "$EXPORT_LINE" >> "$SHELL_CONFIG"
        print_info "已将路径添加到 $SHELL_CONFIG"
        print_info "请运行 'source $SHELL_CONFIG' 或重新打开终端使更改生效"
    else
        print_info "路径已在 $SHELL_CONFIG 中配置"
    fi
}

# 显示安装后说明
post_install_info() {
    print_info "安装完成！"
    echo ""
    print_info "文件已安装到: ${HOME}/.lain-cli"
    echo ""
    print_info "目录内容:"
    ls -la "${HOME}/.lain-cli"
    echo ""
    print_info "已将 ${HOME}/.lain-cli 添加到 PATH 环境变量中:"
    echo "  export PATH=\$PATH:${HOME}/.lain-cli"
    export PATH=$PATH:${HOME}/.lain-cli
    echo ""
    echo "运行 '${HOME}/.lain-cli/lain-cli --help' 查看可用命令"
    echo ""
    echo "注意：为了让 lain-cli 命令在所有新的终端会话中都可用，请执行以下操作之一："
    echo "1. 运行 'source ~/.bashrc' (或者对应的 shell 配置文件)"
    echo "2. 重新打开您的终端"
}

# 主函数
main() {
    print_info "开始安装 lain-cli"
    
    check_root
    detect_platform
    create_temp_dir
    get_latest_version
    
    # 下载预编译版本
    if ! download_prebuilt; then
        print_error "无法下载预编译版本"
        exit 1
    fi
    
    install_files
    post_install_info
}

# 改进的脚本执行检查方式，修复 Bad substitution 错误
# 如果直接运行此脚本，则执行 main 函数
script_name="$0"
if [ "${script_name#-}" = "$script_name" ]; then
    main "$@"
fi