package config

import (
	"flag"
	"os"
)

// Config 结构体存储全局配置
type Config struct {
	Port         string
	MMDBCityPath string
	MMDBASNPath  string
	ISPDictDir   string // ★ 修改：从 ISPDictPath 改为 ISPDictDir
	MainDomain   string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "port", "8080", "Server port")
	flag.StringVar(&cfg.MMDBCityPath, "city-db", "/usr/share/GeoIP/GeoLite2-City.mmdb", "Path to City MMDB")
	flag.StringVar(&cfg.MMDBASNPath, "asn-db", "/usr/share/GeoIP/GeoLite2-ASN.mmdb", "Path to ASN MMDB")

	// ★ 修改：默认路径改为 data/isp 目录，参数名改为 -isp-dir
	flag.StringVar(&cfg.ISPDictDir, "isp-dir", "data/isp", "Directory containing ISP translation JSON files")

	flag.StringVar(&cfg.MainDomain, "domain", "api.myipdns.com", "Main domain served via Cloudflare")

	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		cfg.Port = envPort
	}

	return cfg
}
