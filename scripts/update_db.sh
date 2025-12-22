#!/bin/bash

# 1. 运行官方更新工具
# 请确保 /etc/GeoIP.conf 已配置好您的 License Key
/usr/bin/geoipupdate

# 获取命令返回值
STATUS=$?

if [ $STATUS -eq 0 ]; then
    echo "[$(date)] GeoIP database updated successfully."
    
   

    echo "Reloading Go API service..."
    systemctl restart myip-api.service
else
    echo "[$(date)] GeoIP update failed."
fi