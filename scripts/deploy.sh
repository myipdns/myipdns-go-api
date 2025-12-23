#!/bin/bash
echo "=== 开始更新 API ==="

cd /opt/myipdns-go-api || exit

echo "1. 拉取代码..."
git pull

echo "2. 下载依赖 (新增步骤)..."
# ★★★ 关键：自动分析代码，下载缺少的包，删除多余的包 ★★★
go mod tidy

echo "3. 重新编译..."
go build -o myip-api cmd/server/main.go

# 检查编译是否成功，如果失败就不要重启了，免得把服务搞挂
if [ $? -ne 0 ]; then
    echo "❌ 编译失败！请检查报错信息。"
    exit 1
fi

echo "4. 重启服务..."
systemctl restart myip-api

echo "=== 更新完成 ==="
systemctl status myip-api --no-pager