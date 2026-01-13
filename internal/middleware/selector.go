package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ContextKey 定义 Context 中存储的键名，避免硬编码字符串
const (
	CtxClientIP = "client_ip"
	CtxMode     = "app_mode" // "full" or "simple"
)

// SelectorConfig 定义中间件配置
type SelectorConfig struct {
	MainDomain string // 主域名 (CDN)
}

// NewSelector 初始化分流中间件
func NewSelector(cfg SelectorConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. 获取 Hostname (e.g., api.myipdns.com, apiv4.myipdns.com)
		host := c.Hostname()

		// 2. 判断请求模式 & 获取真实 IP
		// 默认模式：Simple (直连/纯文本)，安全性最高，不查库
		mode := "simple"
		realIP := ""

		// 检查是否是主域名 (CDN 模式)
		// 只有主域名才允许 "Full" 模式 (查库 + JSON)
		if strings.EqualFold(host, cfg.MainDomain) {
			mode = "full"

			// --- CDN 模式 IP 获取逻辑 ---
			// 优先信任 Cloudflare 传递的 Header
			// 注意：这里假设 Nginx 已经过滤了外部伪造的 Header，
			// 或者 Go 程序只监听本地端口，只有 Nginx 能连接。
			realIP = c.Get("CF-Connecting-IP")
			if realIP == "" {
				// 降级：尝试 X-Forwarded-For (Cloudflare -> Nginx -> Go)
				// Fiber 的 c.IP() 会自动处理 XFF，但我们手动控制更稳
				realIP = c.Get("X-Forwarded-For")
			}
		} else {
			// --- 直连模式 IP 获取逻辑 (apiv4 / apiv6) ---
			// 忽略 CF-Connecting-IP (防止用户伪造)
			// 直接信任 Nginx 传过来的 X-Forwarded-For 或者 RemoteIP
			// 注意：c.IP() 在配置了 ProxyHeader 后会读取 XFF
			realIP = c.IP()
		}

		// 兜底：如果上面的逻辑都没拿到，强制取 RemoteIP
		if realIP == "" {
			realIP = c.Context().RemoteAddr().String()
		}

		// 处理 IP 格式:
		// 有时 XFF 会包含多个 IP "client, proxy1, proxy2"，我们取第一个
		// 处理 IP 格式:
		// 有时 XFF 会包含多个 IP "client, proxy1, proxy2"，我们取第一个
		if idx := strings.IndexByte(realIP, ','); idx >= 0 {
			realIP = strings.TrimSpace(realIP[:idx])
		}

		// 3. 将结果存入 Fiber 的 Locals (请求上下文)
		// 后续的 Handler 直接取这个值，不需要再计算
		c.Locals(CtxClientIP, realIP)
		c.Locals(CtxMode, mode)

		return c.Next()
	}
}
