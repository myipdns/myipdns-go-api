package config

import (
	"flag"
	"os"
)

// Config 结构体存储全局配置
// 这种结构避免了多次读取环境变量，提升性能
type Config struct {
	Port         string
	MMDBCityPath string
	MMDBASNPath  string
	ISPDictPath  string
	MainDomain   string // 用于识别是否走 CDN (api.myipdns.com)
}

func Load() *Config {
	// 定义默认值，支持通过命令行覆盖
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "port", "8080", "Server port")
	flag.StringVar(&cfg.MMDBCityPath, "city-db", "/usr/share/GeoIP/GeoLite2-City.mmdb", "Path to City MMDB")
	flag.StringVar(&cfg.MMDBASNPath, "asn-db", "/usr/share/GeoIP/GeoLite2-ASN.mmdb", "Path to ASN MMDB")
	flag.StringVar(&cfg.ISPDictPath, "isp-dict", "data/isp_zh.json", "Path to ISP translation map")
	flag.StringVar(&cfg.MainDomain, "domain", "api.myipdns.com", "Main domain served via Cloudflare")

	flag.Parse()

	// 环境变量优先 (兼容 Docker/Systemd)
	if envPort := os.Getenv("PORT"); envPort != "" {
		cfg.Port = envPort
	}
	// ... 可以按需添加其他 ENV 覆盖 ...

	return cfg
}