#!/bin/bash

# 1. 运行官方更新工具
# 请确保 /etc/GeoIP.conf 已配置好您的 License Key
/usr/bin/geoipupdate

# 获取命令返回值
STATUS=$?

if [ $STATUS -eq 0 ]; then
    echo "[$(date)] GeoIP database updated successfully."
    
    # 2. 检查数据库文件是否真的变新了 (可选，geoipupdate 只有在有更新时才下载)
    # 这里简单粗暴一点，只要运行成功就尝试重启 API
    # 也可以比较文件 md5，但 MaxMind 每周才更新一次，重启一次无伤大雅
    
    echo "Reloading Go API service..."
    systemctl restart myip-api.service
else
    echo "[$(date)] GeoIP update failed."
fi