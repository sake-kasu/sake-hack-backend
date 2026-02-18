package entity

import "errors"

var (
	ErrSakeNotFound    = errors.New("酒が見つかりません")
	ErrInvalidCategory = errors.New("無効なカテゴリです")
)
