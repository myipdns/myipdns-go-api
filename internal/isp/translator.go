package isp

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Translator 管理多语言 ISP 翻译字典
type Translator struct {
	// 数据结构升级：第一层Key是语言代码(如 "cn"), 第二层Key是英文ISP名称
	mappings map[string]map[string]string
	mu       sync.RWMutex
}

// NewTranslator 从指定目录加载所有 .json 文件
func NewTranslator(dirPath string) (*Translator, error) {
	t := &Translator{
		mappings: make(map[string]map[string]string),
	}

	// 1. 读取目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// 如果目录不存在，不报错，仅打印警告并返回空翻译器（允许降级运行）
		if os.IsNotExist(err) {
			log.Printf("[ISP] Warning: Dict directory '%s' not found. ISP translation disabled.", dirPath)
			return t, nil
		}
		return nil, err
	}

	// 2. 遍历加载每个 json 文件
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// 提取文件名作为语言代码 (例如 "cn.json" -> "cn")
		langCode := strings.TrimSuffix(entry.Name(), ".json")
		filePath := filepath.Join(dirPath, entry.Name())

		// 读取文件内容
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[ISP] Error reading %s: %v", entry.Name(), err)
			continue
		}

		// 解析 JSON
		var dict map[string]string
		if err := json.Unmarshal(data, &dict); err != nil {
			log.Printf("[ISP] Error parsing JSON %s: %v", entry.Name(), err)
			continue
		}

		t.mappings[langCode] = dict
		log.Printf("[ISP] Loaded dictionary: %s (%d entries)", langCode, len(dict))
	}

	return t, nil
}

// Translate 尝试翻译 ISP 名称
// raw: 原始英文 ISP (例如 "Google LLC")
// lang: 目标语言代码 (例如 "cn", "ja", "ru")
func (t *Translator) Translate(raw string, lang string) string {
	// 简单规范化语言代码，兼容 URL 参数
	targetLang := lang
	if lang == "zh" || lang == "zh-CN" {
		targetLang = "cn"
	}
	if lang == "jp" {
		targetLang = "ja"
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	// 1. 查找对应语言的字典
	dict, exists := t.mappings[targetLang]
	if !exists {
		// 如果没有该语言的字典（比如用户请求了 lang=fr 但你还没做 fr.json）
		// 直接返回原始英文名 (Fallback to English)
		return raw
	}

	// 2. 在字典中查找 ISP 名称
	if val, ok := dict[raw]; ok && val != "" {
		return val
	}

	// 3. 如果字典里没这个 ISP，返回原始英文名
	return raw
}
