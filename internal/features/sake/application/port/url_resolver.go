package port

import "context"

// URLResolver はオブジェクトキーから署名付きURLに変換するインターフェース
type URLResolver interface {
	// ResolveURL はオブジェクトキーから署名付きダウンロードURLを生成する
	ResolveURL(ctx context.Context, objectKey string) (string, error)
}
