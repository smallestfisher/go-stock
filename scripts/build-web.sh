#!/bin/bash
# go-stock web version build script

# Exit on error
set -e

# 获取脚本所在目录的绝对路径
SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$SCRIPTS_DIR/.."

cd "$PROJECT_ROOT"

echo ">>> 步骤 1: 构建前端 (Vue)..."
cd frontend
# 确保安装了依赖
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖中..."
    npm install
fi

# 打包前端静态资源
export NODE_OPTIONS="--max-old-space-size=4096"
npm run build

echo ">>> 步骤 2: 构建后端 (Go)..."
cd "$PROJECT_ROOT"
# 编译二进制文件
go build -o go-stock-web main_web.go app.go app_common.go app_web_bridge.go assets.go utils.go

echo "------------------------------------------"
echo "  编译完成！"
echo "  您现在可以运行 ./start-web.sh 启动服务了。"
echo "------------------------------------------"
