package geo

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// Provider 封装 GeoIP 数据库操作
type Provider struct {
	cityReader *geoip2.Reader
	asnReader  *geoip2.Reader
}

// Result 定义标准的返回结构
type Result struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ASN         uint   `json:"asn,omitempty"`
	ISP         string `json:"isp,omitempty"`
}

// NewProvider 初始化数据库
// 使用 mmap 模式，对 1GB VPS 非常友好
func NewProvider(cityPath, asnPath string) (*Provider, error) {
	p := &Provider{}

	var err error
	p.cityReader, err = geoip2.Open(cityPath)
	if err != nil {
		return nil, err
	}

	p.asnReader, err = geoip2.Open(asnPath)
	if err != nil {
		// ASN 库是可选的，如果打开失败，只记录日志不中断程序
		log.Printf("[Warning] Failed to open ASN DB: %v. ISP info will be missing.", err)
	}

	return p, nil
}

// Close 关闭数据库句柄
func (p *Provider) Close() {
	if p.cityReader != nil {
		p.cityReader.Close()
	}
	if p.asnReader != nil {
		p.asnReader.Close()
	}
}

// Lookup 执行查询并处理多语言逻辑
func (p *Provider) Lookup(ipStr string, userLang string) (*Result, error) {
	ip := net.ParseIP(ipStr)
	res := &Result{IP: ipStr}

	// 1. 语言映射：将用户输入的简写映射为 MaxMind 的 Key
	// 默认为英文
	mmLang := "en"
	
	switch userLang {
	case "cn", "zh", "zh-CN", "zh-cn":
		mmLang = "zh-CN"
	case "ru":
		mmLang = "ru"
	case "jp", "ja": // 兼容用户可能传 jp
		mmLang = "ja"
	case "fr":
		mmLang = "fr"
	case "de":
		mmLang = "de"
	case "es":
		mmLang = "es"
	case "pt", "pt-BR", "pt-br":
		mmLang = "pt-BR"
	}

	// 2. 查询 City 库 (地理位置)
	if p.cityReader != nil {
		if record, err := p.cityReader.City(ip); err == nil {
			// 国家
			res.CountryCode = record.Country.IsoCode
			res.Country = getName(record.Country.Names, mmLang)

			// 城市
			res.City = getName(record.City.Names, mmLang)

			// 省份/区域 (取第一个 subdivisions)
			if len(record.Subdivisions) > 0 {
				res.Region = getName(record.Subdivisions[0].Names, mmLang)
			}
			
			// 经纬度
			res.Latitude = record.Location.Latitude
			res.Longitude = record.Location.Longitude
		}
	}

	// 3. 查询 ASN 库
	if p.asnReader != nil {
		if record, err := p.asnReader.ASN(ip); err == nil {
			res.ASN = record.AutonomousSystemNumber
			res.ISP = record.AutonomousSystemOrganization
		}
	}

	return res, nil
}

// 辅助函数：安全获取 Map 值
func getName(names map[string]string, lang string) string {
	// 1. 尝试取目标语言
	if val, ok := names[lang]; ok {
		return val
	}
	// 2. 尝试降级到英文
	if val, ok := names["en"]; ok {
		return val
	}
	// 3. 如果连英文都没有，随机取一个或者返回空 (极少见)
	for _, v := range names {
		return v
	}
	return ""
}

// 辅助函数：安全获取 Map 值，如果目标语言不存在则回退到英文
func getName(names map[string]string, lang string) string {
	if val, ok := names[lang]; ok {
		return val
	}
	if val, ok := names["en"]; ok {
		return val
	}
	return ""
}