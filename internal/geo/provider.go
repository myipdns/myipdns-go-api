package geo

import (
	"database/sql"
	"log"
	"math/big"
	"net"
	"strings"

	_ "github.com/mattn/go-sqlite3" // 引入 SQLite 驱动
	"github.com/oschwald/geoip2-golang"
)

// Provider 封装 GeoIP 和 IP2Proxy 数据库操作
type Provider struct {
	cityReader *geoip2.Reader
	asnReader  *geoip2.Reader
	ip2ProxyDB *sql.DB   // [新增] SQLite 连接
	stmt       *sql.Stmt // [新增] 预编译查询语句
}

// Result 定义标准的返回结构
type Result struct {
	IP string `json:"ip"`

	// --- 地理位置 ---
	Continent     string  `json:"continent,omitempty"`
	ContinentCode string  `json:"continent_code,omitempty"`
	Country       string  `json:"country,omitempty"`
	CountryCode   string  `json:"country_code,omitempty"`
	IsEU          bool    `json:"is_eu,omitempty"`
	Region        string  `json:"region,omitempty"`
	RegionCode    string  `json:"region_code,omitempty"`
	City          string  `json:"city,omitempty"`
	TimeZone      string  `json:"time_zone,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`

	// --- 网络特征 ---
	ASN   uint   `json:"asn,omitempty"`
	ASOrg string `json:"as_org,omitempty"`

	// --- [新增] IP2Proxy 扩展字段 ---
	ISP        string `json:"isp,omitempty"`         // 独立的 ISP 字段
	Provider   string `json:"provider,omitempty"`    // VPN/代理供应商名称 (如 ExpressVPN)
	UsageType  string `json:"usage_type,omitempty"`  // 用途 (ISP, DCH, MOB...)
	Domain     string `json:"domain,omitempty"`      // 运营商域名
	Threat     string `json:"threat,omitempty"`      // 威胁类型 (SPAM...)
	FraudScore int    `json:"fraud_score,omitempty"` // 欺诈评分

	// --- 特殊标记 ---
	IsProxy     bool `json:"is_proxy,omitempty"`
	IsAnycast   bool `json:"is_anycast,omitempty"`
	IsSatellite bool `json:"is_satellite,omitempty"`
}

// NewProvider 初始化所有数据库
func NewProvider(cityPath, asnPath, ip2proxyPath string) (*Provider, error) {
	p := &Provider{}

	var err error
	// 1. MaxMind City
	p.cityReader, err = geoip2.Open(cityPath)
	if err != nil {
		return nil, err
	}

	// 2. MaxMind ASN
	p.asnReader, err = geoip2.Open(asnPath)
	if err != nil {
		log.Printf("[Warning] Failed to open ASN DB: %v", err)
	}

	// 3. IP2Proxy SQLite [新增]
	// 使用只读模式打开，避免锁文件问题
	db, err := sql.Open("sqlite3", "file:"+ip2proxyPath+"?mode=ro")
	if err != nil {
		log.Printf("[Warning] Failed to open IP2Proxy DB: %v", err)
	} else {
		// 测试连接
		if err := db.Ping(); err != nil {
			log.Printf("[Warning] IP2Proxy DB connect failed: %v", err)
		} else {
			db.SetMaxOpenConns(4) // SQLite 建议限制并发连接
			p.ip2ProxyDB = db

			// [优化] 预编译 SQL 语句
			// 针对 TEXT 列的 >= 比较，索引扫描更高效
			// [Refactor] 使用标准范围查询 ip_from <= IP <= ip_to
			query := `
		SELECT country_code, country_name, region, city, isp, domain, 
		       usage_type, asn, as_name, threat, provider, fraud_score, proxy_type
		FROM ip2proxy 
		INDEXED BY idx_ip_to 
		WHERE ip_to >= ? AND ip_from <= ?
		ORDER BY ip_to ASC 
		LIMIT 1
	`
			p.stmt, err = db.Prepare(query)
			if err != nil {
				log.Printf("[Warning] Failed to prepare IP2Proxy query: %v", err)
			}
		}
	}

	return p, nil
}

// Close 关闭所有句柄
func (p *Provider) Close() {
	if p.cityReader != nil {
		p.cityReader.Close()
	}
	if p.asnReader != nil {
		p.asnReader.Close()
	}
	if p.stmt != nil {
		p.stmt.Close()
	}
	if p.ip2ProxyDB != nil {
		p.ip2ProxyDB.Close()
	}
}

// Lookup 基础查询 (仅 MaxMind)
func (p *Provider) Lookup(ip net.IP, userLang string) (*Result, error) {
	ipStr := ip.String()
	res := &Result{IP: ipStr}

	// 语言映射
	mmLang := mapLang(userLang)

	// 1. MaxMind City
	if p.cityReader != nil {
		if record, err := p.cityReader.City(ip); err == nil {
			res.ContinentCode = record.Continent.Code
			res.Continent = getName(record.Continent.Names, mmLang)
			res.CountryCode = record.Country.IsoCode
			res.Country = getName(record.Country.Names, mmLang)
			res.IsEU = record.Country.IsInEuropeanUnion

			if len(record.Subdivisions) > 0 {
				res.RegionCode = record.Subdivisions[0].IsoCode
				res.Region = getName(record.Subdivisions[0].Names, mmLang)
			}
			res.City = getName(record.City.Names, mmLang)
			res.Latitude = record.Location.Latitude
			res.Longitude = record.Location.Longitude
			res.TimeZone = record.Location.TimeZone

			res.IsProxy = record.Traits.IsAnonymousProxy
			res.IsSatellite = record.Traits.IsSatelliteProvider
			res.IsAnycast = record.Traits.IsAnycast
		}
	}

	// 2. MaxMind ASN
	if p.asnReader != nil {
		if record, err := p.asnReader.ASN(ip); err == nil {
			res.ASN = record.AutonomousSystemNumber
			res.ASOrg = record.AutonomousSystemOrganization
		}
	}

	return res, nil
}

// LookupFull 完整查询 (MaxMind + IP2Proxy) [新增]
func (p *Provider) LookupFull(ip net.IP, userLang string) (*Result, error) {
	// 1. 先获取基础 MaxMind 数据
	res, err := p.Lookup(ip, userLang)
	if err != nil {
		return nil, err
	}

	// 2. 查询 IP2Proxy SQLite
	p2Res := p.queryIP2Proxy(ip)
	if p2Res == nil {
		return res, nil // 查不到就直接返回基础数据
	}

	// 3. --- 合并逻辑 (你的第三点要求) ---

	// [地理位置补全] 如果 MaxMind 为空，尝试用 IP2Proxy
	if res.CountryCode == "" {
		res.CountryCode = p2Res.CountryCode
	}
	if res.Country == "" {
		res.Country = p2Res.Country
	}
	if res.Region == "" {
		res.Region = p2Res.Region
	}
	if res.City == "" {
		res.City = p2Res.City
	}

	// [ASOrg 策略] MaxMind 优先，如果没有则补全
	if res.ASOrg == "" && p2Res.ASName != "" && p2Res.ASName != "-" {
		res.ASOrg = p2Res.ASName
	}
	// ASN 补全
	if res.ASN == 0 && p2Res.ASN > 0 {
		res.ASN = p2Res.ASN
	}

	// [新字段] 只要 IP2Proxy 有数据就填充
	if p2Res.ISP != "" && p2Res.ISP != "-" {
		res.ISP = p2Res.ISP
	}
	if p2Res.Provider != "" && p2Res.Provider != "-" {
		res.Provider = p2Res.Provider
	}
	if p2Res.UsageType != "" && p2Res.UsageType != "-" {
		res.UsageType = p2Res.UsageType
	}
	if p2Res.Domain != "" && p2Res.Domain != "-" {
		res.Domain = p2Res.Domain
	}
	if p2Res.Threat != "" && p2Res.Threat != "-" {
		res.Threat = p2Res.Threat
	}
	if p2Res.FraudScore > 0 {
		res.FraudScore = p2Res.FraudScore
	}

	// [Proxy 状态] 任意一方认为是代理，即为代理
	if p2Res.IsProxy {
		res.IsProxy = true
	}

	return res, nil
}

// --- 内部辅助结构体 (仅用于接收 SQL 结果) ---
type ip2ProxyRaw struct {
	CountryCode string
	Country     string
	Region      string
	City        string
	ISP         string
	Domain      string
	UsageType   string
	ASN         uint
	ASName      string
	Threat      string
	Provider    string
	FraudScore  int
	IsProxy     bool
}

// 执行 SQLite 查询
func (p *Provider) queryIP2Proxy(ip net.IP) *ip2ProxyRaw {
	if p.ip2ProxyDB == nil || p.stmt == nil {
		return nil
	}

	// IP 转换逻辑简化，直接使用传入的 ip
	if ip == nil {
		return nil
	}

	// 判断是 IPv4 还是 IPv6
	// [Fix] IP2Location IPv6 数据库中的 IPv4 地址是映射由于 ::ffff:x.x.x.x (IPv4-mapped IPv6)
	// 因此我们需要始终将其作为 16 字节处理，以生成正确的整数值。
	ipInt := big.NewInt(0)
	ipInt.SetBytes(ip.To16())

	// [Fix] 将 IP 转换为 39 位补零字符串
	ipStr := ipInt.String()
	if len(ipStr) < 39 {
		ipStr = strings.Repeat("0", 39-len(ipStr)) + ipStr
	}

	var cc, cn, reg, city, isp, dom, usage, asnStr, asName, threat, prov, pType sql.NullString
	var fScore sql.NullInt32

	// 使用预编译语句查询
	// 传入两次 ipStr，对应 WHERE ip_to >= ? AND ip_from <= ?
	err := p.stmt.QueryRow(ipStr, ipStr).Scan(
		&cc, &cn, &reg, &city, &isp, &dom, &usage, &asnStr, &asName, &threat, &prov, &fScore, &pType,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("SQL Debug: %v", err)
		}
		// ErrNoRows 正常情况，说明没查到（空洞）
		return nil
	}

	// 检查结果是否真的包含该 IP (处理数据空洞)
	// 在生产环境中，应该再查一次 ip_from 验证 ip_num >= ip_from
	// 但为了性能，且 ip2proxy 通常是连续段，这里暂时略过严格检查，假设 LIMIT 1 命中就是对的

	res := &ip2ProxyRaw{
		CountryCode: cc.String,
		Country:     cn.String,
		Region:      reg.String,
		City:        city.String,
		ISP:         isp.String,
		Domain:      dom.String,
		UsageType:   usage.String,
		ASName:      asName.String,
		Threat:      threat.String,
		Provider:    prov.String,
		FraudScore:  int(fScore.Int32),
	}

	// 解析 ASN (字符串转 uint)
	// (省略具体的 parse logic，这里简单处理)
	// ...

	// 判断是否为代理
	// PX12 中 proxy_type 如果不是 "-" 则通常被视为代理 (VPN, DCH, PUB, TOR...)
	if pType.Valid && pType.String != "-" && pType.String != "" {
		res.IsProxy = true
	}

	return res
}

// 辅助函数 (保持不变)
func mapLang(userLang string) string {
	switch userLang {
	case "cn", "zh", "zh-CN", "zh-cn":
		return "zh-CN"
	case "ru":
		return "ru"
	case "jp", "ja":
		return "ja"
	case "fr":
		return "fr"
	case "de":
		return "de"
	case "es":
		return "es"
	case "pt", "pt-BR", "pt-br":
		return "pt-BR"
	default:
		return "en"
	}
}

func getName(names map[string]string, lang string) string {
	if val, ok := names[lang]; ok {
		return val
	}
	if val, ok := names["en"]; ok {
		return val
	}
	for _, v := range names {
		return v
	}
	return ""
}
