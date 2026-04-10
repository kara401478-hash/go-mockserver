package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config はJSONファイル全体の構造
type Config struct {
	Port   int     `json:"port"`
	Routes []Route `json:"routes"`
}

// Route は1つのエンドポイント設定
type Route struct {
	Path     string   `json:"path"`
	Method   string   `json:"method"`
	Response Response `json:"response"`
}

// Response はレスポンスの内容
type Response struct {
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	Delay   int               `json:"delay_ms"`
}

// Load はJSONファイルを読み込んでConfigを返す
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルを読み込めませんでした: %w", err)
	}

	cfg := &Config{Port: 8080}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("JSONのパースに失敗しました: %w", err)
	}

	return validate(cfg)
}

// validate は設定の内容が正しいかチェックする
func validate(cfg *Config) (*Config, error) {
	for i, route := range cfg.Routes {
		if route.Path == "" {
			return nil, fmt.Errorf("route[%d]: path が空です", i)
		}
		if route.Method == "" {
			return nil, fmt.Errorf("route[%d]: method が空です", i)
		}
		if cfg.Routes[i].Response.Status == 0 {
			cfg.Routes[i].Response.Status = 200
		}
	}
	return cfg, nil
}
