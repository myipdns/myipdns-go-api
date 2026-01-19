package main

import (
	"log"
	"myipdns-go-api/internal/config"
	"myipdns-go-api/internal/geo"
	"myipdns-go-api/internal/isp"
	"myipdns-go-api/internal/middleware"
	"net"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	// 初始化翻译器
	ispTrans, err := isp.NewTranslator(cfg.ISPDictDir)
	if err != nil {
		log.Printf("[Warning] Failed to load ISP dicts: %v", err)
	}

	// [修改] 传入新的 DB 路径
	geoProvider, err := geo.NewProvider(cfg.MMDBCityPath, cfg.MMDBASNPath, cfg.IP2ProxyDBPath)
	if err != nil {
		log.Fatalf("Failed to load DBs: %v", err)
	}
	defer geoProvider.Close()

	app := fiber.New(fiber.Config{
		AppName:               "MyIPDNS API/1.3", // 版本号升级
		DisableStartupMessage: false,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "GET,HEAD"}))
	app.Use(middleware.NewSelector(middleware.SelectorConfig{MainDomain: cfg.MainDomain}))

	handler := func(c *fiber.Ctx) error {
		clientIP := c.Locals(middleware.CtxClientIP).(string)
		mode := c.Locals(middleware.CtxMode).(string)

		if mode == "simple" {
			return c.SendString(clientIP)
		}

		targetIP := clientIP
		queryIP := c.Query("ip")
		var parsedIP net.IP

		if queryIP != "" {
			parsedIP = net.ParseIP(queryIP)
			if parsedIP == nil {
				return c.Status(400).JSON(fiber.Map{"error": "invalid_ip", "ip": queryIP})
			}
			targetIP = queryIP
		} else {
			// 如果没有 query ip，解析 clientIP
			parsedIP = net.ParseIP(targetIP)
		}

		// 兜底：如果 clientIP 也是无效的（理论上中间件保证了，但防御性编程）
		if parsedIP == nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid_client_ip", "ip": targetIP})
		}

		lang := c.Query("lang", "en")
		// [新增] 检查 full 参数
		isFull := c.QueryBool("full")

		var result *geo.Result
		var err error

		// [新增] 分支逻辑
		if isFull {
			result, err = geoProvider.LookupFull(parsedIP, lang)
		} else {
			result, err = geoProvider.Lookup(parsedIP, lang)
		}

		if err != nil {
			return c.JSON(fiber.Map{"ip": targetIP, "error": "geo_lookup_failed"})
		}

		// [修改] 翻译逻辑
		// 1. 翻译 ASOrg
		if result.ASOrg != "" {
			result.ASOrg = ispTrans.Translate(result.ASOrg, lang)
		}
		// 2. [新增] 翻译 ISP (复用同一个引擎)
		// 只有当 isFull=true 且 ISP 字段有值时，这里才会生效
		if result.ISP != "" {
			result.ISP = ispTrans.Translate(result.ISP, lang)
		}

		c.Set("Cache-Control", "public, max-age=3600")
		return c.JSON(result)
	}

	app.Get("/", handler)
	app.Get("/ip", handler)

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
