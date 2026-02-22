package port

import "context"

// FileDeleter はストレージ上のファイル削除インターフェース
type FileDeleter interface {
	// DeleteObject はオブジェクトキーに対応するファイルをストレージから削除する
	DeleteObject(ctx context.Context, objectKey string) error
}
