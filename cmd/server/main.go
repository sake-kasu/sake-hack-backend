package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sake-kasu/sake-hack-backend/internal/config"
	"github.com/sake-kasu/sake-hack-backend/internal/database"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
	"github.com/sake-kasu/sake-hack-backend/internal/server"
	"go.uber.org/zap"
)

func main() {
	// 設定読み込み(ロガー初期化前に必要)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// ロガー初期化
	if err := logger.Init(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		log.Fatalf("ロガーの初期化に失敗しました: %v", err)
	}
	defer logger.Sync()

	logger.Get().Info("サーバーを起動します...")

	// PostgreSQL接続
	postgresPool, err := database.NewPostgresPool(database.PostgresConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Database:        cfg.Database.Database,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		logger.Get().Fatal("PostgreSQLへの接続に失敗しました", zap.Error(err))
	}
	logger.Get().Info("PostgreSQLへの接続に成功しました")

	// Valkey接続
	valkeyClient, err := database.NewValkeyClient(database.ValkeyConfig{
		Host:         cfg.Cache.Host,
		Port:         cfg.Cache.Port,
		Password:     cfg.Cache.Password,
		Database:     cfg.Cache.Database,
		PoolSize:     cfg.Cache.PoolSize,
		MinIdleConns: cfg.Cache.MinIdleConns,
		MaxRetries:   cfg.Cache.MaxRetries,
		DialTimeout:  cfg.Cache.DialTimeout,
		ReadTimeout:  cfg.Cache.ReadTimeout,
		WriteTimeout: cfg.Cache.WriteTimeout,
	})
	if err != nil {
		logger.Get().Fatal("Valkeyへの接続に失敗しました", zap.Error(err))
	}
	logger.Get().Info("Valkeyへの接続に成功しました")

	// サーバー作成
	srv := server.New(cfg, postgresPool, valkeyClient)

	// サーバー起動(Goroutine)
	go func() {
		if err := srv.Start(); err != nil {
			logger.Get().Fatal("サーバーエラーが発生しました", zap.Error(err))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Get().Info("シャットダウンシグナルを受信しました。サーバーを停止します...")

	if err := srv.Shutdown(); err != nil {
		logger.Get().Error("サーバーのシャットダウン時にエラーが発生しました", zap.Error(err))
		os.Exit(1)
	}

	logger.Get().Info("サーバーが正常に停止しました")
}
