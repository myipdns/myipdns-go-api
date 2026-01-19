package config

import (
	"flag"
	"os"
)

// Config 结构体存储全局配置
type Config struct {
	Port           string
	MMDBCityPath   string
	MMDBASNPath    string
	ISPDictDir     string
	MainDomain     string
	IP2ProxyDBPath string // [新增] IP2Proxy 数据库路径
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "port", "8080", "Server port")
	flag.StringVar(&cfg.MMDBCityPath, "city-db", "/usr/share/GeoIP/GeoLite2-City.mmdb", "Path to City MMDB")
	flag.StringVar(&cfg.MMDBASNPath, "asn-db", "/usr/share/GeoIP/GeoLite2-ASN.mmdb", "Path to ASN MMDB")
	flag.StringVar(&cfg.ISPDictDir, "isp-dir", "data/isp", "Directory containing ISP translation JSON files")
	flag.StringVar(&cfg.MainDomain, "domain", "api.myipdns.com", "Main domain served via Cloudflare")

	// [新增] 默认路径对应你的脚本生成位置
	flag.StringVar(&cfg.IP2ProxyDBPath, "ip2proxy-db", "/usr/share/ip2location/ip2proxy.db", "Path to IP2Proxy SQLite DB")

	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		cfg.Port = envPort
	}

	return cfg
}
