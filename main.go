package main

import (
	"flag"
	"log"
	"os"

	"github.com/yourusername/go-mockserver/internal/config"
	"github.com/yourusername/go-mockserver/internal/server"
)

func main() {
	// コマンドラインフラグの定義
	configPath := flag.String("config", "routes.yaml", "設定ファイルのパス")
	flag.Parse()

	log.SetFlags(0) // タイムスタンプなしのシンプルなログ

	log.Println("🚀 go-mockserver 起動中...")
	log.Printf("📄 設定ファイル: %s\n", *configPath)

	// 設定ファイルを読み込む
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("❌ エラー: %v", err)
		os.Exit(1)
	}

	log.Printf("📌 登録するルート (%d件):\n", len(cfg.Routes))

	// サーバーを起動
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Printf("❌ サーバーエラー: %v", err)
		os.Exit(1)
	}
}
