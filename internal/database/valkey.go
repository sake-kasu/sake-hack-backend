package database

import (
	"context"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

// ValkeyConfig はValkey接続設定
type ValkeyConfig struct {
	Host         string
	Port         int
	Password     string
	Database     int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewValkeyClient はValkeyクライアントを作成する
func NewValkeyClient(config ValkeyConfig) (valkey.Client, error) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	clientOption := valkey.ClientOption{
		InitAddress: []string{addr},
		SelectDB:    config.Database,
	}

	// パスワードが設定されている場合のみ追加
	if config.Password != "" {
		clientOption.Password = config.Password
	}

	client, err := valkey.NewClient(clientOption)
	if err != nil {
		return nil, fmt.Errorf("valkeyのClientの作成に失敗しました: %w", err)
	}

	// 接続確認
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return nil, fmt.Errorf("valkeyへのPingに失敗しました: %w", err)
	}

	return client, nil
}

// HealthCheckValkey はValkeyの接続状態を確認する
func HealthCheckValkey(ctx context.Context, client valkey.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := client.Do(ctx, client.B().Ping().Build()).Error()
	if err != nil {
		return fmt.Errorf("valkeyへのヘルスチェックに失敗しました: %w", err)
	}

	return nil
}
