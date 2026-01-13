package main

import (
	"log"
	"myipdns-go-api/internal/config"
	"myipdns-go-api/internal/geo"
	"myipdns-go-api/internal/isp"
	"myipdns-go-api/internal/middleware"
	"net"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化 ISP 数据层 (这里改用了 ISPDictDir)
	ispTrans, err := isp.NewTranslator(cfg.ISPDictDir)
	if err != nil {
		log.Printf("[Warning] Failed to load ISP dicts: %v", err)
	}

	geoProvider, err := geo.NewProvider(cfg.MMDBCityPath, cfg.MMDBASNPath)
	if err != nil {
		log.Fatalf("Failed to load MaxMind DB: %v", err)
	}
	defer geoProvider.Close()
	// 3. 初始化 Fiber App
	app := fiber.New(fiber.Config{
		AppName:               "MyIPDNS API/1.2",
		DisableStartupMessage: false,
		Immutable:             false,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		IdleTimeout:           30 * time.Second,
		ProxyHeader:           "X-Forwarded-For",
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "GET,HEAD"}))

	app.Use(middleware.NewSelector(middleware.SelectorConfig{
		MainDomain: cfg.MainDomain,
	}))

	// 4. 核心路由处理逻辑
	handler := func(c *fiber.Ctx) error {
		clientIP := c.Locals(middleware.CtxClientIP).(string)
		mode := c.Locals(middleware.CtxMode).(string)

		// 极简模式：直接返回 IP
		if mode == "simple" {
			return c.SendString(clientIP)
		}

		// 完整模式：支持查询指定 IP
		targetIP := clientIP
		queryIP := c.Query("ip")
		if queryIP != "" {
			if net.ParseIP(queryIP) == nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "invalid_ip_format",
					"ip":    queryIP,
				})
			}
			targetIP = queryIP
		}

		// 获取语言参数 (默认 en)
		lang := c.Query("lang", "en")

		// 1. 查 MaxMind 库
		result, err := geoProvider.Lookup(targetIP, lang)
		if err != nil {
			return c.JSON(fiber.Map{"ip": targetIP, "error": "geo_lookup_failed"})
		}

		// 2. 翻译 ISP (★ 核心修改：逻辑简化 ★)
		// 如果有 ISP 信息，直接扔给翻译器。
		// 翻译器会根据 lang 自动查找对应的 json (如 ja.json)，找不到就原样返回。
		if result.ASOrg != "" {
			result.ASOrg = ispTrans.Translate(result.ASOrg, lang)
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
