package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/oschwald/maxminddb-golang"
)

// 定义 ASN 数据库的记录结构
type ASNRecord struct {
	ASOrg string `maxminddb:"autonomous_system_organization"`
}

func main() {
	dbPath := "/usr/share/GeoIP/GeoLite2-ASN.mmdb"
	
	// 打开数据库
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	// 使用 Map 去重
	uniqueISPs := make(map[string]bool)
	fmt.Println("正在遍历数据库提取 ISP 名称，这可能需要几秒钟...")

	// 遍历所有网段
	networks := db.Networks(maxminddb.SkipAliasedNetworks)
	var record ASNRecord
	count := 0

	for networks.Next() {
		_, err := networks.Network(&record)
		if err != nil {
			continue
		}
		if record.ASOrg != "" {
			uniqueISPs[record.ASOrg] = true
		}
		count++
		if count%100000 == 0 {
			fmt.Printf("已处理 %d 个网段...\n", count)
		}
	}

	// 转换为切片并排序
	var ispList []string
	for name := range uniqueISPs {
		ispList = append(ispList, name)
	}
	sort.Strings(ispList)

	// 写入文件
	file, err := os.Create("isp_list_raw.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, isp := range ispList {
		_, _ = writer.WriteString(isp + "\n")
	}
	writer.Flush()

	fmt.Printf("\n提取完成！\n总共发现 %d 个唯一的 ISP 名称。\n已保存到 isp_list_raw.txt\n", len(ispList))
}