package isp

import (
	"encoding/json"
	"os"
	"sync"
)

// Translator 管理 ISP 翻译字典
type Translator struct {
	mapping map[string]string
	mu      sync.RWMutex // 虽然只读，但为了未来可能的热更新保留锁
}

// NewTranslator 从 JSON 文件加载字典
func NewTranslator(filePath string) (*Translator, error) {
	t := &Translator{
		mapping: make(map[string]string),
	}

	// 如果文件不存在，直接返回一个空的翻译器（不报错，仅降级）
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return t, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	if err := json.Unmarshal(data, &t.mapping); err != nil {
		return nil, err
	}

	return t, nil
}

// Translate 尝试翻译 ISP 名称
// lang: 目标语言 "cn" 或 "en"
// raw: 原始英文 ISP
func (t *Translator) Translate(raw string, lang string) string {
	// 如果不是中文模式，直接返回原文
	if lang != "cn" {
		return raw
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	// 查找字典
	if val, ok := t.mapping[raw]; ok {
		return val
	}

	// 没找到则返回原文
	return raw
}