package main

import (
	"log"
	"myipdns-go-api/internal/config"
	"myipdns-go-api/internal/geo"
	"myipdns-go-api/internal/isp"
	"myipdns-go-api/internal/middleware"
	"net" // [新增] 用于校验 IP 格式
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化数据层
	ispTrans, err := isp.NewTranslator(cfg.ISPDictPath)
	if err != nil {
		log.Printf("[Warning] Failed to load ISP dict: %v", err)
	}

	geoProvider, err := geo.NewProvider(cfg.MMDBCityPath, cfg.MMDBASNPath)
	if err != nil {
		log.Fatalf("Failed to load MaxMind DB: %v", err)
	}
	defer geoProvider.Close()

	// 3. 初始化 Fiber App
	app := fiber.New(fiber.Config{
		AppName:               "MyIPDNS API/1.1", 
		DisableStartupMessage: false,
		Immutable:             false,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		IdleTimeout:           30 * time.Second,
		ProxyHeader:           "X-Forwarded-For",
	})

	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*", AllowMethods: "GET,HEAD"}))

	app.Use(middleware.NewSelector(middleware.SelectorConfig{
		MainDomain: cfg.MainDomain,
	}))

	// 4. 核心路由处理逻辑
	handler := func(c *fiber.Ctx) error {
		// 获取连接者的真实 IP (来自中间件)
		clientIP := c.Locals(middleware.CtxClientIP).(string)
		mode := c.Locals(middleware.CtxMode).(string)

		// === 分支 A: 极简模式 (直连域名) ===
		// 直连域名 (v4/v6) 通常只返回当前连接 IP，不支持查询别人
		// 这样可以防止恶意用户利用直连通道进行大规模爬取
		if mode == "simple" {
			return c.SendString(clientIP)
		}

		// === 分支 B: 完整模式 (主域名 CDN) ===
		
		// [新增功能] 支持 ?ip=x.x.x.x 查询指定 IP
		// 默认查询目标是当前连接者
		targetIP := clientIP
		
		queryIP := c.Query("ip")
		if queryIP != "" {
			// 如果指定了 ip 参数，校验格式
			if net.ParseIP(queryIP) == nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "invalid_ip_format",
					"ip":    queryIP,
				})
			}
			targetIP = queryIP
		}

		// 获取语言参数 ?lang=cn
		lang := c.Query("lang", "en")
		if lang != "cn" && lang != "en" {
			lang = "en"
		}

		// 1. 查 MaxMind 库 (查询 targetIP)
		result, err := geoProvider.Lookup(targetIP, lang)
		if err != nil {
			return c.JSON(fiber.Map{"ip": targetIP, "error": "geo_lookup_failed"})
		}

		// 2. 翻译 ISP
		if result.ISP != "" {
			result.ISP = ispTrans.Translate(result.ISP, lang)
		}

		// 3. 设置缓存头
		c.Set("Cache-Control", "public, max-age=3600")

		// 4. 返回结果
		return c.JSON(result)
	}

	// 注册路由
	app.Get("/", handler)
	app.Get("/ip", handler)

	// 5. 启动服务器
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}