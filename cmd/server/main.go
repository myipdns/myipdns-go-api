package main

import (
	"log"
	"myipdns-go-api/internal/config"
	"myipdns-go-api/internal/geo"
	"myipdns-go-api/internal/isp"
	"myipdns-go-api/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. [新增] 初始化数据层
	// 加载 ISP 字典
	ispTrans, err := isp.NewTranslator(cfg.ISPDictPath)
	if err != nil {
		log.Printf("[Warning] Failed to load ISP dict: %v", err)
	}

	// 加载 MaxMind 数据库
	geoProvider, err := geo.NewProvider(cfg.MMDBCityPath, cfg.MMDBASNPath)
	if err != nil {
		// 如果核心库加载失败，程序必须退出，否则无法工作
		log.Fatalf("Failed to load MaxMind DB: %v", err)
	}
	defer geoProvider.Close() // 确保退出时关闭文件句柄

	// 3. 初始化 Fiber App
	app := fiber.New(fiber.Config{
		AppName:               "MyIPDNS API/1.0",
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
	
	// 智能分流
	app.Use(middleware.NewSelector(middleware.SelectorConfig{
		MainDomain: cfg.MainDomain,
	}))

	// 4. [更新] 核心路由处理逻辑
	handler := func(c *fiber.Ctx) error {
		// 从 Context 取出中间件分析好的数据
		ip := c.Locals(middleware.CtxClientIP).(string)
		mode := c.Locals(middleware.CtxMode).(string)

		// === 分支 A: 极简模式 (直连域名) ===
		// 仅返回 IP 字符串，零 CPU 消耗，不查库
		if mode == "simple" {
			return c.SendString(ip)
		}

		// === 分支 B: 完整模式 (主域名 CDN) ===
		// 获取语言参数 ?lang=cn
		lang := c.Query("lang", "en")
		if lang != "cn" && lang != "en" {
			lang = "en"
		}

		// 1. 查 MaxMind 库 (内存操作, <1ms)
		result, err := geoProvider.Lookup(ip, lang)
		if err != nil {
			// 如果查库失败（极少见），至少返回 IP
			return c.JSON(fiber.Map{"ip": ip, "error": "geo_lookup_failed"})
		}

		// 2. 翻译 ISP (内存 Map 查找, 纳秒级)
		// 如果 MaxMind 查到了 ISP 且用户请求中文
		if result.ISP != "" {
			result.ISP = ispTrans.Translate(result.ISP, lang)
		}

		// 3. 设置缓存头 (关键优化)
		// 允许浏览器和 Cloudflare 缓存 1 小时
		c.Set("Cache-Control", "public, max-age=3600")

		// 4. 返回最终 JSON
		return c.JSON(result)
	}

	// 注册路由
	app.Get("/", handler)
	app.Get("/ip", handler) // 兼容性路由

	// 5. 启动服务器
	log.Printf("Starting server on port %s...", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}