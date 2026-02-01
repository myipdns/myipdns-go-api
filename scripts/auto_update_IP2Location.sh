#!/bin/bash

# ================= 配置区域 =================
# 1. 填入您的永久 Token
TOKEN=""

# 2. 数据库代码 (IPv6 CSV版)
CODE="PX12LITECSVIPV6"

# 3. 目标保存目录
TARGET_DIR="/usr/share/ip2location"
TARGET_DB="ip2proxy.db"

# 4. 临时工作目录
TMP_DIR="/root/ip2proxy_temp_build"
# ===========================================

# 检查 root 权限
if [ "$EUID" -ne 0 ]; then
  echo "Error: Please run as root (sudo)."
  exit 1
fi

# 设置临时的语言环境，防止 Python 报错
export LC_ALL=C.UTF-8
export LANG=C.UTF-8

# 检查 python3
if ! command -v python3 &> /dev/null; then
    echo "Error: python3 not found. Run: apt install python3 -y"
    exit 1
fi

echo ">> Start updating IP2Proxy database..."
mkdir -p $TMP_DIR
mkdir -p $TARGET_DIR

# ==========================================
# 第一步：下载
# ==========================================
echo ">> Downloading database ($CODE)..."
DOWNLOAD_URL="https://www.ip2location.com/download?token=$TOKEN&file=$CODE"
HTTP_CODE=$(curl -L -s -w "%{http_code}" -o "$TMP_DIR/database.zip" "$DOWNLOAD_URL")

if [ "$HTTP_CODE" -ne 200 ]; then
    echo "Download failed! HTTP Code: $HTTP_CODE"
    rm -rf $TMP_DIR
    exit 1
fi

FILE_SIZE=$(stat -c%s "$TMP_DIR/database.zip")
if [ "$FILE_SIZE" -lt 100000 ]; then
    echo "Error: Downloaded file is too small ($FILE_SIZE bytes)."
    rm -rf $TMP_DIR
    exit 1
fi

# ==========================================
# 第二步：解压
# ==========================================
echo ">> Unzipping..."
unzip -o -j "$TMP_DIR/database.zip" "*.CSV" -d "$TMP_DIR" > /dev/null 2>&1
CSV_FILE=$(find "$TMP_DIR" -name "*.CSV" | head -n 1)

if [ -z "$CSV_FILE" ]; then
    echo "Error: CSV file not found in zip."
    rm -rf $TMP_DIR
    exit 1
fi

# ==========================================
# 第三步：生成转换脚本 (已修复编码问题)
# ==========================================
CONVERT_SCRIPT="$TMP_DIR/convert.py"

# 注意：下面这一行 # -*- coding: utf-8 -*- 是修复的关键
cat <<EOF > "$CONVERT_SCRIPT"
# -*- coding: utf-8 -*-
import sqlite3
import csv
import sys
import os

csv_file = sys.argv[1]
db_file = sys.argv[2]

if os.path.exists(db_file):
    os.remove(db_file)

print(f"Converting CSV to SQLite... Target: {db_file}")

conn = sqlite3.connect(db_file)
c = conn.cursor()

c.execute('PRAGMA synchronous = OFF')
c.execute('PRAGMA journal_mode = MEMORY')

# --- 创建表结构 (16列) ---
# [Fix] ip_from/ip_to 改为 TEXT 以支持大数排序 (Pad to 39 chars)
c.execute('''
    CREATE TABLE ip2proxy (
        ip_from      TEXT,     -- [1] Start IP (Padded 39 chars)
        ip_to        TEXT,     -- [2] End IP (Padded 39 chars)
        proxy_type   TEXT,     -- [3] Proxy Type
        country_code TEXT,     -- [4] Country Code
        country_name TEXT,     -- [5] Country Name
        region       TEXT,     -- [6] Region
        city         TEXT,     -- [7] City
        isp          TEXT,     -- [8] ISP
        domain       TEXT,     -- [9] Domain
        usage_type   TEXT,     -- [10] Usage Type
        asn          TEXT,     -- [11] ASN
        as_name      TEXT,     -- [12] AS Name
        last_seen    TEXT,     -- [13] Last Seen
        threat       TEXT,     -- [14] Threat
        provider     TEXT,     -- [15] Provider
        fraud_score  INTEGER   -- [16] Fraud Score
    )
''')

# --- 读取 CSV 并插入 ---
# encoding='utf-8' 确保读取 CSV 时不报错
with open(csv_file, 'r', encoding='utf-8', errors='replace') as f:
    reader = csv.reader(f)
    batch = []
    batch_size = 50000
    count = 0
    
    for row in reader:
        if not row or row[0].startswith('#'):
            continue
            
        row_data = row[:16]
        if len(row_data) < 16:
             row_data += [None] * (16 - len(row_data))

        # [Fix] Padding IP to 39 digits
        try:
            # 假设 CSV 里是十进制数字符串
            row_data[0] = str(row_data[0]).zfill(39)
            row_data[1] = str(row_data[1]).zfill(39)
        except:
            continue # Skip invalid rows

        batch.append(row_data)
        count += 1
        
        if len(batch) >= batch_size:
            c.executemany('INSERT INTO ip2proxy VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)', batch)
            conn.commit()
            batch = []
            if count % 100000 == 0:
                print(f"Processed {count} rows...", end='\r')

    if batch:
        c.executemany('INSERT INTO ip2proxy VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)', batch)
        conn.commit()

print(f"\nImport finished. Total rows: {count}")

print("Creating index (idx_ip_to)...")
# INDEX on TEXT column works lexicographically (which is what we want for padded strings)
c.execute('CREATE INDEX idx_ip_to ON ip2proxy(ip_to)')

print("Creating index (idx_ip_from)...")
c.execute('CREATE INDEX idx_ip_from ON ip2proxy(ip_from)')

conn.close()
EOF

# ==========================================
# 第四步：执行转换
# ==========================================
TEMP_DB="$TMP_DIR/temp_ip2proxy.db"
# 显式使用 UTF-8 环境运行 Python
LC_ALL=C.UTF-8 python3 "$CONVERT_SCRIPT" "$CSV_FILE" "$TEMP_DB"

if [ $? -ne 0 ]; then
    echo "Conversion failed!"
    rm -rf $TMP_DIR
    exit 1
fi

# ==========================================
# 第五步：部署
# ==========================================
echo "Deploying to $TARGET_DIR/$TARGET_DB ..."
mv -f "$TEMP_DB" "$TARGET_DIR/$TARGET_DB"
chmod 644 "$TARGET_DIR/$TARGET_DB"

rm -rf $TMP_DIR
echo "Success! Database updated."
ls -lh "$TARGET_DIR/$TARGET_DB"