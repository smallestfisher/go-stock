#!/bin/bash

# --- 配置区 ---
# 您可以在这里修改默认的访问账号和密码
export APP_USER=${APP_USER:-"admin"}
export APP_PASSWORD=${APP_PASSWORD:-"admin"}
export PORT=${PORT:-"8080"}
# --------------

BINARY="./go-stock-web"

# 检查二进制文件是否存在
if [ ! -f "$BINARY" ]; then
    echo "错误: 未找到 $BINARY 文件。"
    echo "请先运行 ./scripts/build-web.sh 进行编译。"
    exit 1
fi

echo "------------------------------------------"
echo "  go-stock Web 版正在启动..."
echo "  访问地址: http://localhost:$PORT"
echo "  访问账号: $APP_USER"
echo "  访问密码: $APP_PASSWORD"
echo "------------------------------------------"

# 启动服务
# 如果您想在后台运行，可以使用: nohup $BINARY > web_server.log 2>&1 &
$BINARY
