package repository

import "context"

// LikeRepository いいね操作のリポジトリインターフェース
type LikeRepository interface {
	// Create いいねを作成する(重複時はDO NOTHING)
	Create(ctx context.Context, sakeID int32, token string) error
	// Delete いいねを削除する(削除件数を返す)
	Delete(ctx context.Context, sakeID int32, token string) (int64, error)
	// GetCount 酒のいいね数を取得する
	GetCount(ctx context.Context, sakeID int32) (int64, error)
	// Exists いいねの存在確認
	Exists(ctx context.Context, sakeID int32, token string) (bool, error)
}
