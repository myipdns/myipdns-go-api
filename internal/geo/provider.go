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
// 字段命名参考 MaxMind CSV 格式标准
type Result struct {
	IP string `json:"ip"`

	// --- 地理位置 (Location) ---
	Continent     string  `json:"continent,omitempty"`      // 大洲
	ContinentCode string  `json:"continent_code,omitempty"` // AS, EU, NA 等
	Country       string  `json:"country,omitempty"`        // 国家
	CountryCode   string  `json:"country_code,omitempty"`   // CN, US 等
	IsEU          bool    `json:"is_eu,omitempty"`          // 是否为欧盟
	Region        string  `json:"region,omitempty"`         // 省/州 (原 Subdivision 1)
	RegionCode    string  `json:"region_code,omitempty"`    // 省/州 ISO 代码
	City          string  `json:"city,omitempty"`           // 城市
	TimeZone      string  `json:"time_zone,omitempty"`      // 时区 (如 Asia/Shanghai)
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`

	// --- 网络特征 (Network Traits) ---
	ASN   uint   `json:"asn,omitempty"`
	ASOrg string `json:"as_org,omitempty"` // 运营商 (如 China Telecom)

	// --- 特殊标记 (Special Tags) ---
	IsProxy     bool `json:"is_proxy,omitempty"`     // 是否为匿名代理/VPN
	IsAnycast   bool `json:"is_anycast,omitempty"`   // 是否为 Anycast 广播网络
	IsSatellite bool `json:"is_satellite,omitempty"` // 是否为卫星服务商 (如 Starlink)
}

// NewProvider 初始化数据库
func NewProvider(cityPath, asnPath string) (*Provider, error) {
	p := &Provider{}

	var err error
	p.cityReader, err = geoip2.Open(cityPath)
	if err != nil {
		return nil, err
	}

	p.asnReader, err = geoip2.Open(asnPath)
	if err != nil {
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
	mmLang := "en"
	switch userLang {
	case "cn", "zh", "zh-CN", "zh-cn":
		mmLang = "zh-CN"
	case "ru":
		mmLang = "ru"
	case "jp", "ja":
		mmLang = "ja"
	case "fr":
		mmLang = "fr"
	case "de":
		mmLang = "de"
	case "es":
		mmLang = "es"
	case "pt", "pt-BR", "pt-br":
		mmLang = "pt-BR"
	default:
		// 如果用户传了其他奇怪的语言，默认保持 en
		if userLang != "" {
			mmLang = "en"
		}
	}

	// 2. 查询 City 库 (地理位置)
	if p.cityReader != nil {
		if record, err := p.cityReader.City(ip); err == nil {
			// Continent
			res.ContinentCode = record.Continent.Code
			res.Continent = getName(record.Continent.Names, mmLang)

			// Country
			res.CountryCode = record.Country.IsoCode
			res.Country = getName(record.Country.Names, mmLang)
			res.IsEU = record.Country.IsInEuropeanUnion

			// Subdivisions (Region)
			if len(record.Subdivisions) > 0 {
				res.RegionCode = record.Subdivisions[0].IsoCode
				res.Region = getName(record.Subdivisions[0].Names, mmLang)
			}

			// City
			res.City = getName(record.City.Names, mmLang)

			// Location
			res.Latitude = record.Location.Latitude
			res.Longitude = record.Location.Longitude
			res.TimeZone = record.Location.TimeZone

			// Traits
			res.IsProxy = record.Traits.IsAnonymousProxy
			res.IsSatellite = record.Traits.IsSatelliteProvider
			res.IsAnycast = record.Traits.IsAnycast
		}
	}

	// 3. 查询 ASN 库 (运营商)
	if p.asnReader != nil {
		if record, err := p.asnReader.ASN(ip); err == nil {
			res.ASN = record.AutonomousSystemNumber
			res.ASOrg = record.AutonomousSystemOrganization
		}
	}

	return res, nil
}

// getName 辅助函数：安全获取 Map 值
func getName(names map[string]string, lang string) string {
	// 1. 尝试取目标语言
	if val, ok := names[lang]; ok {
		return val
	}
	// 2. 尝试降级到英文
	if val, ok := names["en"]; ok {
		return val
	}
	// 3. 如果连英文都没有，返回任意一个或空
	for _, v := range names {
		return v
	}
	return ""
}
